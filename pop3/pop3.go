package pop3

import (
	"bufio"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"log"
	"my-golang-project/lemail/db"
	"my-golang-project/lemail/util"
	"net"
	"net/mail"
	"sort"
	"strconv"
	"strings"
	"sync"
)

func GetConn(account *db.Account) (c *Conn) {
	address := fmt.Sprintf("%s:%d", account.POPServer, util.POP3SSLPORT)
	conn, err := Dial(address, true)
	if err != nil {
		fmt.Println(err)
	}

	if err := conn.Auth(account.Email, account.PassWord); err != nil {
		return
	}
	return

}

func Dial(address string, isTsl bool) (c *Conn, err error) {
	conn, err := net.Dial("tcp", address)
	if isTsl {
		conn, err = tls.Dial("tcp", address, nil)
	}
	if err != nil {
		return
	}
	return NewConn(conn)
}

func NewConn(conn net.Conn) (c *Conn, err error) {
	c = &Conn{
		r:    bufio.NewReader(conn),
		w:    bufio.NewWriter(conn),
		conn: conn,
	}
	line, err := c.ReadLine()
	if err != nil {
		return
	}

	if !IsOK(line) {
		return nil, errors.New("pop3: Server did not respond with +OK")
	}
	return
}

func (c *Conn) ReadLine() (line string, err error) {
	b, _, err := c.r.ReadLine()
	if err == io.EOF {
		return
	}
	if err != nil {
		return
	}
	line = string(b)
	return
}

func (c *Conn) ReadLines() (lines []string, err error) {
	for {
		line, err := c.ReadLine()
		if err != nil {
			return nil, err
		}

		if line == "." {
			break
		}
		lines = append(lines, line)
	}
	return
}

func (c *Conn) Cmd(format string,
	args ...interface{}) (line string, err error) {
	if err = c.Send(format, args...); err != nil {
		return
	}

	line, err = c.ReadLine()
	if err != nil {
		return
	}
	if !IsOK(line) {
		return "", errors.New("pop3: Server did not respond with +OK")
	}
	return
}

func (c *Conn) User(u string) (err error) {
	if _, err = c.Cmd("%s %s\r\n", USER, u); err != nil {
		return
	}
	return
}

func (c *Conn) Pass(p string) (err error) {
	if _, err = c.Cmd("%s %s\r\n", PASS, p); err != nil {
		return
	}
	return
}

func (c *Conn) Send(format string, args ...interface{}) (err error) {
	if _, err = c.w.WriteString(fmt.Sprintf(format, args...)); err != nil {
		return
	}
	return c.w.Flush()
}

func (c *Conn) Quit() (err error) {
	if err = c.Send("%s\r\n", QUIT); err != nil {
		return
	}
	return c.conn.Close()
}

func (c *Conn) Auth(u, p string) (err error) {
	if err = c.User(u); err != nil {
		return
	}

	if err = c.Pass(p); err != nil {
		return
	}

	return c.Noop()
}

func (c *Conn) Noop() (err error) {
	if _, err = c.Cmd("%s\r\n", NOOP); err != nil {
		return
	}
	return
}

// 这里开始
func (c *Conn) STAT() (stat string, err error) {
	if stat, err = c.Cmd("%s\r\n", STAT); err != nil {
		return
	}
	return // +OK 994 27748563
}

func (c *Conn) LISTAll() (list []int, err error) {
	if _, err = c.Cmd("%s\r\n", LIST); err != nil {
		return
	}
	lines, err := c.ReadLines()
	if err != nil {
		return
	}
	for _, v := range lines {
		id, err := strconv.Atoi(strings.Fields(v)[0])
		if err != nil {
			return nil, err
		}
		list = append(list, id)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(list)))
	return
}

func (c *Conn) RETR(id int, mess chan MessageHeader, wait *sync.WaitGroup) {

	fmt.Println(id)
	if _, err := c.Cmd("%s %d\r\n", RETR, id); err != nil {
		log.Println(1)
		log.Println(err)
		return
	}
	m, err := mail.ReadMessage(c.r)
	if err != nil {
		log.Println(2)
		log.Println(err)
		return
	}

	line, err := c.ReadLine()
	if err != nil {
		log.Println(3)
		log.Println(err)
		return
	}
	if line != "." {
		if err = c.r.UnreadByte(); err != nil {
			log.Println(4)
			log.Println(err)
			return
		}
	}
	var subject string
	if len(m.Header["Subject"]) == 0 {
		subject = ""
	} else {
		subject = Subject(m.Header["Subject"][0])
	}
	messageHeader := MessageHeader{
		ContentType: strings.Join(m.Header["Content-Type"], ""),
		MessageID:   strings.Join(m.Header["Message-Id"], ""),
		From:        Other(strings.Join(m.Header["From"], "")),
		To:          Other(strings.Join(m.Header["To"], "")),
		Cc:          Other(strings.Join(m.Header["Cc"], "")),
		Bcc:         Other(strings.Join(m.Header["Bcc"], "")),
		Date:        strings.Join(m.Header["Date"], ""),
		Subject:     subject,
		MimeVersion: strings.Join(m.Header["Mime-Version"], ""),
	}
	mess <- messageHeader

	wait.Done()
}

func (c *Conn) RETRS(id int) (messageHeader *MessageHeader) {

	if _, err := c.Cmd("%s %d\r\n", RETR, id); err != nil {
		log.Println(err)
		return
	}
	m, err := mail.ReadMessage(c.r)
	if err != nil {
		log.Println(err)
		return
	}

	line, err := c.ReadLine()
	if err != nil {
		log.Println(err)
		return
	}
	if line != "." {
		if err = c.r.UnreadByte(); err != nil {
			log.Println(err)
			return
		}
	}
	var subject string
	if len(m.Header["Subject"]) == 0 {
		subject = ""
	} else {
		subject = Subject(m.Header["Subject"][0])
	}
	return &MessageHeader{
		ContentType: strings.Join(m.Header["Content-Type"], ""),
		MessageID:   strings.Join(m.Header["Message-Id"], ""),
		From:        From(Other(strings.Join(m.Header["From"], ""))),
		To:          Other(strings.Join(m.Header["To"], "")),
		Cc:          Other(strings.Join(m.Header["Cc"], "")),
		Bcc:         Other(strings.Join(m.Header["Bcc"], "")),
		Date:        strings.Join(m.Header["Date"], ""),
		Subject:     subject,
		MimeVersion: strings.Join(m.Header["Mime-Version"], ""),
	}
}
