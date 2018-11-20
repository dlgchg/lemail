package db

import (
	"time"
	"strings"
	"fmt"
)

type SendEmail struct {
	Id         int64     `xorm:"pk autoincr"`
	FromEmail  string    `xorm:"fromEmail"`
	ToEmail    string    `xorm:"toEmail"`
	CcEmail    string    `xorm:"ccEmail"`
	BccEmail   string    `xorm:"bccEmail"`
	CreateTime time.Time `xorm:"created"`
	Title      string
	Content    string
	Attach     string
	State      string
}

func NewSendEmail(account *Account, ToEmail, CcEmail, BccEmail []string, Title, Content, Attach string, err error) *SendEmail {
	var errStr string
	if err != nil {
		errStr = err.Error()
	} else {
		errStr = ""
	}
	return &SendEmail{
		FromEmail: account.Email,
		ToEmail:   strings.Replace(strings.Trim(fmt.Sprint(ToEmail), "[]"), " ", ",", -1),
		CcEmail:   strings.Replace(strings.Trim(fmt.Sprint(CcEmail), "[]"), " ", ",", -1),
		BccEmail:  strings.Replace(strings.Trim(fmt.Sprint(BccEmail), "[]"), " ", ",", -1),
		Title:     Title,
		Content:   Content,
		Attach:    Attach,
		State:     errStr,
	}
}


