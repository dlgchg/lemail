package email

import (
	"github.com/go-gomail/gomail"
	"my-golang-project/lemail/db"
)

func SendEmail(toEmail []string, ccEmail []string, bccEmail []string,
	title, body, attach string) (*db.Account, error) {
	account := db.NowUsingEmailInfo()
	m := gomail.NewMessage()
	m.SetHeader("From", account.Email) // your
	m.SetHeader("To", toEmail...) //he
	if len(ccEmail) > 0 {
		m.SetHeader("Cc", ccEmail...)
	}
	if len(bccEmail) > 0 {
		m.SetHeader("Bcc", bccEmail...)
	}
	m.SetHeader("Subject", title)
	m.SetBody("text/html", body)
	if len(attach) > 0 {
		m.Attach(attach)
	}
	dialer := gomail.NewDialer(account.Server, account.SSL, account.Email, account.PassWord)
	if err := dialer.DialAndSend(m); err != nil {
		return account, err
	}
	return account, nil
}
