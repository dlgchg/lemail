package main

import (
	"fmt"
	"github.com/satori/go.uuid"
	"github.com/urfave/cli"
	"io/ioutil"
	"my-golang-project/EmailTE-Go/db"
	"my-golang-project/EmailTE-Go/email"
	"my-golang-project/EmailTE-Go/util"
	"os"
	"strconv"
	"strings"
	"time"
)

var toEmail, ccEmail, bccEmail []string

func main() {

	var load bool

	app := cli.NewApp()
	app.Name = "EmailTE-Go"
	app.Author = "LiWei/EmailTE-Go"
	app.Copyright = fmt.Sprintf("(c) 2018-%s iikira.", time.Now().Format("2006"))
	app.Description = "EmailTE-Go是一个命令行的邮件客户端"
	app.Usage = "EmailTE-Go是一个命令行的邮件客户端"
	app.Version = "0.0.1"
	//----------------
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:        "load, l",
			Usage:       `查看是否有新邮件，待开发.`,
			Destination: &load,
		},
	}
	app.Action = func(c *cli.Context) {
		if c.NArg() != 0 {
			fmt.Printf("没有发现命令: %s\n运行命令 %s help 获取帮助\n", c.Args().Get(0), app.Name)
			return
		}
		fmt.Println("Load......")
	}
	//----------------
	app.Commands = []cli.Command{
		{
			Name:     "add",
			Aliases:  []string{"a"},
			Usage:    "新增一个邮箱",
			Category: "EmailTE-Go",
			Description: `
-type:
	0. qq
	1. 163
	2. gmail
-sp:
	true:  smtp
	false: pop3
示例: 
	EmailTE-Go add -type=1 -sp=true -email=xxxx@qq.com -pass=xxxxxx`,
			Action: func(c *cli.Context) error {
				emailType := c.Int("type")
				sp := c.Bool("sp")
				email := c.String("email")
				pass := c.String("pass")

				account := new(db.Account)
				account.Aliases = util.GetEmailName(emailType)
				account.Server = fmt.Sprintf(util.GetServer(sp), util.GetEmailName(emailType))
				account.SSL = util.GetSSL(sp)
				account.Email = email
				account.PassWord = pass
				account.UUID = uuid.Must(uuid.NewV4()).String()
				_, err := db.Engine.Insert(account)
				if err != nil {
					return err
				}
				fmt.Println("新增邮箱成功!")
				return nil
			},
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:  "type,t",
					Usage: "邮箱类型",
				},
				cli.BoolFlag{
					Name:  "sp",
					Usage: "SMTP 还是 POP3",
				},
				cli.StringFlag{
					Name:  "email,e",
					Usage: "邮箱地址",
				},
				cli.StringFlag{
					Name:  "pass,p",
					Usage: "邮箱密码",
				},
			},
		},
		{
			Name:        "show",
			Description: "显示已添加的邮箱信息",
			Category:    "EmailTE-Go",
			Usage:       "显示已添加的邮箱信息",
			Aliases:     []string{"s"},
			Action: func(c *cli.Context) error {
				var accounts []db.Account
				err := db.Engine.Find(&accounts)
				if err != nil {
					fmt.Println("请新增一个邮箱!")
				}

				fmt.Println("你的邮箱列表:")
				fmt.Printf("%-36s     %-20s     %-20s     %-10s     %-30v\n", "uuid", "email", "password", "aliases", "createTime")
				for _, account := range accounts {
					fmt.Printf("%s     %-20s     %-20s     %-10s     %-30v\n", account.UUID, account.Email, account.PassWord, account.Aliases, account.CreateTime)
				}
				return nil
			},
		},
		{
			Name: "del",
			Description: `
删除一个邮箱.
示例:
	EmailTE-Go del -uuid=9b76ea1c-d37c-44e5-a330-cf6ecb882807
`,
			Category: "EmailTE-Go",
			Usage:    "删除一个邮箱",
			Aliases:  []string{"d"},
			Action: func(c *cli.Context) error {
				var account db.Account
				uuId := c.String("uuid")

				if uuId == "" {
					fmt.Println("uuid不能为空!")
					return nil
				}

				account.UUID = uuId
				_, err := db.Engine.Delete(&account)
				if err != nil {
					fmt.Println("uuid不存在!")
					return err
				}
				fmt.Println("删除邮箱成功!")
				return nil
			},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "uuid, u",
					Usage: "UUID",
				},
			},
		},
		{
			Name:     "up",
			Aliases:  []string{"u"},
			Usage:    "修改邮箱信息.",
			Category: "EmailTE-Go",
			Description: `
示例: 
	EmailTE-Go up -uuid=xxxxxxx -email=xxxx@qq.com -pass=xxxxxx`,
			Action: func(c *cli.Context) error {
				var email, pass, uuId string
				uuId = c.String("uuid")
				email = c.String("email")
				pass = c.String("pass")

				account := new(db.Account)
				account.UUID = uuId
				account.Email = email
				account.PassWord = pass
				_, err := db.Engine.Update(account)
				if err != nil {
					return err
				}
				fmt.Println("修改邮箱信息成功!")
				return nil
			},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "uuid,u",
					Usage: "邮箱的uuid",
				},
				cli.StringFlag{
					Name:  "email,e",
					Usage: "邮箱的地址",
				},
				cli.StringFlag{
					Name:  "pass,p",
					Usage: "邮箱的密码",
				},
			},
		},
		{
			Name: "use",
			Description: `
使用邮箱来进行邮件发送.
示例:
	EmailTE-Go use -uuid=9b76ea1c-d37c-44e5-a330-cf6ecb882807
`,
			Category: "EmailTE-Go",
			Usage:    "使用邮箱来进行邮件发送",
			Aliases:  []string{"use"},
			Action: func(c *cli.Context) error {
				uuId := c.String("uuid")

				if uuId == "" {
					fmt.Println("uuid不能为空!")
					return nil
				}

				useAccount := new(db.Account)
				useAccount.UUID = uuId
				_, err := db.Engine.Get(useAccount)
				if err != nil {
					fmt.Println("uuid不存在!")
					return err
				}

				file, err := os.OpenFile("use.txt", os.O_CREATE|os.O_WRONLY, 0777)
				defer file.Close()
				if err != nil {
					fmt.Println("再试一次!")
					return err
				}
				_, err = file.WriteString(useAccount.Email)
				if err != nil {
					fmt.Println("再试一次!")
					return err
				}
				fmt.Println("正在使用的邮箱:", useAccount.Email)
				return nil
			},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "uuid, u",
					Usage: "UUID",
				},
			},
		},
		{
			Name: "using",
			Description: `
查询正在使用的邮箱.
示例:
	EmailTE-Go using`,
			Category: "EmailTE-Go",
			Usage:    "查询正在使用的邮箱.",
			Aliases:  []string{"using"},
			Action: func(c *cli.Context) error {
				file, err := os.Open("use.txt")
				defer file.Close()
				if err != nil {
					fmt.Println("没有正在使用的邮箱!")
					return err
				}
				bytes, err := ioutil.ReadAll(file)
				if err != nil {
					fmt.Println("没有正在使用的邮箱!")
					return err
				}
				fmt.Println("正在使用的邮箱:", string(bytes))
				return nil
			},
		},
		{
			Name:     "send",
			Aliases:  []string{"s"},
			Category: "EmailTE-Go",
			Usage:    "发送一封带抄送，暗送，附件，多人的邮件",
			Description: `
示例：
简单发送:
	EmailTE-Go send -to=xxxxx@xx.com -title=xxxx -body=xxxxxxxx
发送多人:
	EmailTE-Go send -to=xxxxx@xx.com,xxxxx@xx.com -title=xxxx -body=xxxxxxxx
添加附件:
	EmailTE-Go send -to=xxxxx@xx.com -title=xxxx -body=xxxxxxxx -attach=x/x/xxx.jpg
添加抄送和暗送:
	EmailTE-Go send -to=xxxxx@xx.com -cc=xxxxx@xxx.com -bcc=xxxxxx@xx.com -title=xxxx -body=xxxxxxxx -attach=x/x/xxx.jpg`,
			Action: func(c *cli.Context) error {
				to := c.String("to")
				cc := c.String("cc")
				bcc := c.String("bcc")
				title := c.String("title")
				body := c.String("body")
				attach := c.String("attach")

				if len(to) > 0 {
					toEmail = strings.Split(to, ",")
				}
				if len(cc) > 0 {
					ccEmail = strings.Split(cc, ",")
				}
				if len(bcc) > 0 {
					bccEmail = strings.Split(bcc, ",")
				}

				err := util.CheckIsEmpty(c, to, cc, bcc, title, body)
				if err != nil {
					return err
				}

				attach, err = util.CopyFileToProjectAttach(attach)
				if err != nil {
					return err
				}

				_, err = email.SendEmail(toEmail, ccEmail, bccEmail, title, body, attach)
				if err == nil {
					fmt.Println("发送成功!")
				} else {
					fmt.Println("发送失败,失败原因", err)
				}
				return nil
			},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "title",
					Usage: "邮件标题",
				},
				cli.StringFlag{
					Name:  "body",
					Usage: "邮件内容",
				},
				cli.StringFlag{
					Name:  "to",
					Usage: "收件人",
				},
				cli.StringFlag{
					Name:  "cc",
					Usage: "抄送人",
				},
				cli.StringFlag{
					Name:  "bcc",
					Usage: "暗送人",
				},
				cli.StringFlag{
					Name:  "attach",
					Usage: "附件",
				},
			},
		},
		{
			Name:     "send-simple",
			Aliases:  []string{"ss"},
			Category: "EmailTE-Go",
			Usage:    "发送一封简单的邮件,可发送多人",
			Description: `
示例:
	One:
		EmailTE-Go send-simple -to=xxxxx@xx.com -title=xxxx -body=xxxxxxxxg
	More:
		EmailTE-Go send-simple -to=xxxxx@xx.com,xxx@xx.com -title=xxxx -body=xxxxxxxxg`,
			Action: func(c *cli.Context) error {

				to := c.String("to")
				title := c.String("title")
				body := c.String("body")

				if len(to) > 0 {
					toEmail = strings.Split(to, ",")
				}

				err := util.CheckIsEmpty(c, to, title, body)
				if err != nil {
					return err
				}
				account, err := email.SendEmail(toEmail, ccEmail, bccEmail, title, body, "")
				if err == nil {
					fmt.Println("发送成功!")
				} else {
					fmt.Println("发送失败,失败原因", err)
				}
				sendEmail := db.NewSendEmail(account, toEmail, ccEmail, bccEmail, title, body, "", err)
				db.Engine.Insert(sendEmail)
				return nil
			},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "title",
					Usage: "邮件标题",
				},
				cli.StringFlag{
					Name:  "body",
					Usage: "邮件内容",
				},
				cli.StringFlag{
					Name:  "to",
					Usage: "收件人",
				},
			},
		},
		{
			Name:     "send-attach",
			Aliases:  []string{"sa"},
			Category: "EmailTE-Go",
			Usage:    "发送一封带附件的邮件,可发送多人",
			Description: `
示例:
	One:
		EmailTE-Go send-attach -to=xxxxx@xx.com -title=xxxx -body=xxxxxxxxg -attach=/xx/xx/xx.png
	More:
		EmailTE-Go send-attach -to=xxxxx@xx.com,xxx@xx.com -title=xxxx -body=xxxxxxxxg -attach=/xx/xx/xx.png`,
			Action: func(c *cli.Context) error {

				to := c.String("to")
				title := c.String("title")
				body := c.String("body")
				attach := c.String("attach")
				if len(to) > 0 {
					toEmail = strings.Split(to, ",")
				}

				err := util.CheckIsEmpty(c, to, title, body)
				if err != nil {
					return err
				}

				attach, err = util.CopyFileToProjectAttach(attach)
				if err != nil {
					return err
				}

				account, err := email.SendEmail(toEmail, ccEmail, bccEmail, title, body, attach)
				if err == nil {
					fmt.Println("发送成功!")
				} else {
					fmt.Println("发送失败,失败原因", err)
				}
				sendEmail := db.NewSendEmail(account, toEmail, ccEmail, bccEmail, title, body, attach, err)
				db.Engine.Insert(sendEmail)
				return nil
			},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "title",
					Usage: "邮件标题",
				},
				cli.StringFlag{
					Name:  "body",
					Usage: "邮件内容",
				},
				cli.StringFlag{
					Name:  "to",
					Usage: "收件人",
				},
				cli.StringFlag{
					Name:  "attach",
					Usage: "附件",
				},
			},
		},
		{
			Name:     "send-list",
			Aliases:  []string{"sl"},
			Category: "EmailTE-Go",
			Usage:    "查看已发送的邮件,默认查看前20条",
			Description: `
示例:
	One:
		EmailTE-Go send-list -all=true`,
			Action: func(c *cli.Context) error {
				var sendEmails []db.SendEmail
				var err error
				all := c.Bool("all")

				if all {
					err = db.Engine.Desc("create_time").Find(&sendEmails)
				} else {
					err = db.Engine.Where("limit = 20").Desc("create_time").Find(&sendEmails)
				}

				if err != nil || len(sendEmails) == 0 {
					fmt.Println("没有邮件!")
					return nil
				}

				fmt.Println("你已发送的邮件列表:")
				fmt.Printf("%-10s %-20s %-20s %-20s %-30s \n",
					"id", "fromEmail", "create_time", "title", "content")
				for _, sendEmail := range sendEmails {
					fmt.Printf("%-10d %-20s %-20s %-20s %-30s \n",
						sendEmail.Id, sendEmail.FromEmail, sendEmail.CreateTime.Format("2006-01-02 15:04:05"),
						util.Substr(sendEmail.Title, 0, 20), util.Substr(sendEmail.Content, 0, 30))
				}
				return nil
			},
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "all",
					Usage: "查看全部",
				},
			},
		},
		{
			Name:     "send-remove",
			Aliases:  []string{"sr"},
			Category: "EmailTE-Go",
			Usage:    "删除已发出的邮件",
			Description: `
示例:
	One:
		EmailTE-Go send-remove -id=1 -all=false
	All:
		EmailTE-Go send-remove -id= -all=true`,
			Action: func(c *cli.Context) error {
				var sendEmail db.SendEmail
				all := c.Bool("all")
				if all {
					db.Engine.Exec("delete from send_email")
				} else {
					id := c.String("id")

					if id == "" {
						fmt.Println("id不能为空!")
						return nil
					}

					Id, err := strconv.Atoi(id)
					if err != nil {
						fmt.Println("id不能为空!")
						return err
					}
					sendEmail.Id = int64(Id)
					_, err = db.Engine.Delete(&sendEmail)
					if err != nil {
						fmt.Println("id不存在!")
						return err
					}

				}
				fmt.Println("删除邮箱成功!")
				return nil
			},
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "all",
					Usage: "删除全部",
				},
				cli.StringFlag{
					Name:  "id",
					Usage: "邮件id",
				},
			},
		},
	}

	app.Run(os.Args)
}
