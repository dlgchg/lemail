package db

import (
	"os"
	"log"
	"io/ioutil"
)

func NowUsingEmailInfo() (useAccount *Account) {
	file, err := os.Open("use.txt")
	defer func() {
		log.Println("defer file close")
		file.Close()
	}()
	if err != nil {
		log.Println("没有可用邮箱!")
		return
	}
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println("没有可用邮箱!")
		return
	}
	email := string(bytes)

	useAccount = new(Account)
	useAccount.Email = email
	_, err = Engine.Get(useAccount)
	if err != nil {
		log.Println("uuid不存在!")
		return
	}
	return
}
