package pop3

import (
	"bufio"
	"net"
)

type Conn struct {
	r    *bufio.Reader
	w    *bufio.Writer
	conn net.Conn
}
