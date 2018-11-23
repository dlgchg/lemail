package util

const (
	QW = iota
	WY
	GMAIL
)

func GetEmailName(a int) (name string) {
	switch a {
	case QW:
		name = "qq"
	case WY:
		name = "163"
	case GMAIL:
		name = "gmail"
	default:
		name = "qq"
	}
	return
}

func GetSSL(c bool) (ssl int) {
	if c { //smtp
		ssl = SMTPSSLPORT
	} else { //pop3
		ssl = POP3SSLPORT
	}
	return
}

const (
	SMTPSERVER  = "smtp.%s.com"
	SMTPSSLPORT = 465
	POP3SERVER  = "pop.%s.com"
	POP3SSLPORT = 995
)
