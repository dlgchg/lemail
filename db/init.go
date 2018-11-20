package db

import (
	"github.com/go-xorm/xorm"
	_ "github.com/mattn/go-sqlite3"
)

var (
	Engine *xorm.Engine
)

func init() {
	Engine, _ = xorm.NewEngine("sqlite3", "email_te.db")
	Engine.Sync2(new(Account),new(SendEmail))
}
