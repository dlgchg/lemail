# lemail
lemail仿 Linux shell 命令的邮箱命令行客户端.

<!-- TOC -->
- [命令列表及说明](#命令列表及说明)
    - [新增邮箱](#新增邮箱)
    - [邮箱列表](#邮箱列表)
    - [删除邮箱](#删除邮箱)
    - [修改邮箱信息](#修改邮箱信息)
    - [使用邮箱](#使用邮箱)
    - [查看正在使用的邮箱](#查看正在使用的邮箱)
    - [查看10封邮件](#查看10封邮件)
    - [发送邮件](#发送邮件)
        - [发送一封简单的邮件](#发送一封简单的邮件)
        - [发送带附件的邮件](#发送带附件的邮件)
        - [发送完整的邮件](#发送完整的邮件)
    - [查看已发送邮件](#查看已发送邮件)
    - [删除邮件](#删除邮件)
    - [后续开发](#后续开发)
- [交流反馈](#交流反馈)
<!-- /TOC -->
# 命令列表及说明

## 新增邮箱
新增一个邮箱
```
lemail add
```

例子
```
lemail add -type=1 -email=123@163.com -pass=123456
// type 邮箱类型，目前只有QQ(0),163(1),Gmail(2)
```


## 邮箱列表
显示已添加的邮箱信息
```
lemail show
```

## 删除邮箱
使用分配的uuid来删除邮箱信息
```
lemail del
```

例子
```
lemail del -uuid=9b76ea1c-d37c-44e5-a330-cf6ecb882807
```

## 修改邮箱信息
```
lemail up
```

例子
```
lemail up -uuid=9b76ea1c-d37c-44e5-a330-cf6ecb882807 -email=1234@163.com -pass=123456789
```

## 使用邮箱
使用邮箱来进行发送邮件操作
```
lemail use
```

例子
```
lemail use -uuid=9b76ea1c-d37c-44e5-a330-cf6ecb882807
```

## 查看10封邮件
```
lemail -l
```


## 查看正在使用的邮箱
```
lemail using
```

## 发送邮件
### 发送一封简单的邮件
只有接收人，标题和正文，接收人支持多人，以逗号分割
```
lemail send-simple
```

例子
```
//单发
lemail send-simple -to=xxxxx@xx.com -title=xxxx -body=xxxxxxxxg
//多发
lemail send-simple -to=xxxxx@xx.com,xxx@xx.com -title=xxxx -body=xxxxxxxxg
```

### 发送带附件的邮件
```
lemail send-attach
```

例子
```
lemail send-attach -to=xxxxx@xx.com -title=xxxx -body=xxxxxxxxg -attach=/xx/xx/xx.png
```

### 发送完整的邮件
包含以上内容，追加抄送人和暗送人，支持多人，以逗号分割
```
lemail send
```

例子
```
//简单发送
lemail send -to=xxxxx@xx.com -title=xxxx -body=xxxxxxxx
//发送多人
lemail send -to=xxxxx@xx.com,xxxxx@xx.com -title=xxxx -body=xxxxxxxx
//添加附件
lemail send -to=xxxxx@xx.com -title=xxxx -body=xxxxxxxx -attach=x/x/xxx.jpg
//添加抄送和暗送
lemail send -to=xxxxx@xx.com -cc=xxxxx@xxx.com -bcc=xxxxxx@xx.com -title=xxxx -body=xxxxxxxx -attach=x/x/xxx.jpg
```


## 查看已发送邮件
默认显示前20条
```
lemail send-list
```

例子
```
lemail send-list -all=true //-all=true 全部
```

## 删除邮件
```
lemail send-remove
```

例子
```
//根据id删除
lemail send-remove -id=1 -all=false
//全部删除
lemail send-remove -id= -all=true
```

## 后续开发
1. 添加配置文件，邮箱以配置文件添加
2. 加载收到的邮件

# 交流反馈
提交Issue: [Issues](https://github.com/UOYO/lemail/issues)

邮箱: curmido@gmail.com