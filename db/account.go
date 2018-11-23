package db

import "time"

type Account struct {
	Id          int64     `xorm:"pk autoincr"`
	UUID        string    `xorm:"unique 'uuid'"`
	Email       string    `xorm:"unique"`
	PassWord    string    `xorm:'password'"`
	SMTPServer  string    `xorm:"smtp_server"`
	POPServer   string    `xorm:"pop_server"`
	POPSSL      int       `xorm:"pop_ssl"`
	SMTPSSL     int       `xorm:"smtp_ssl"`
	CreateTime  time.Time `xorm:"created"`
	UpdatedTime time.Time `xorm:"updated"`
	Aliases     string
}
