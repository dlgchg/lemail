package db

import "time"

type Account struct {
	Id          int64     `xorm:"pk autoincr"`
	UUID        string    `xorm:"unique 'uuid'"`
	Email       string    `xorm:"unique"`
	PassWord    string    `xorm:'password'"`
	Server      string    `xorm:"server"`
	SSL         int       `xorm:"ssl"`
	CreateTime  time.Time `xorm:"created"`
	UpdatedTime time.Time `xorm:"updated"`
	Aliases     string
}
