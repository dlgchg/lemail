package pop3

const (
	USER = "USER"
	PASS = "PASS"
	QUIT = "QUIT"
	STAT = "STAT" //处理请求服务器发回关于邮箱的统计资料，如邮件总数和总字节数
	UIDL = "UIDL" //处理返回邮件的唯一标识符，POP3会话的每个标识符都将是唯一的
	LIST = "LIST" //处理返回邮件数量和每个邮件的大小
	RETR = "RETR" //处理返回由参数标识的邮件的全部文本
	DELE = "DELE" //处理服务器将由参数标识的邮件标记为删除，由quit命令执行
	RSET = "RSET" //处理服务器将重置所有标记为删除的邮件，用于撤消DELE命令
	TOP  = "TOP"  //处理服务器将返回由参数标识的邮件前n行内容，n必须是正整数
	NOOP = "NOOP" //处理服务器返回一个肯定的响应
)

const (
	OK  = "+OK"
	ERR = "-ERR"
)

const (
	GbkQu    = "=?gbk?Q?undefined"
	GBKQu    = "=?GBK?Q?undefined"
	GbkBu    = "=?gbk?B?undefined"
	GBKBu    = "=?GBK?B?undefined"
	GB2312Qu = "=?GB2312?Q?undefined"
	Gb2312Qu = "=?gb2312?Q?undefined"
	Gb2312Bu = "=?gb2312?B?undefined"
	GB2312Bu = "=?GB2312?B?undefined"
	GB2312B  = "=?GB2312?B?"
	Gb2312B  = "=?gb2312?B?"
	Gb2312Q  = "=?gb2312?Q?"
	GB2312Q  = "=?GB2312?Q?"
	UTF8B    = "=?UTF-8?B?"
	Utf8B    = "=?utf-8?B?"
	UTF8D    = "=?UTF-8?D?"
	Utf8D    = "=?utf-8?D?"
	GBKB     = "=?GBK?B?"
	GbkB     = "=?gbk?B?"
	GBKQ     = "=?GBK?Q?"
	GbkQ     = "=?gbk?Q?"
)

const DateFormat = "Mon, 02 Jan 2006 15:04:05 MST -0700"
