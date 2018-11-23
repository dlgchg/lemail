package pop3

import "sync"

type MessageHeader struct {
	ContentType string
	MessageID   string
	From        string
	To          string
	Cc          string
	Bcc         string
	Date        string
	Subject     string
	MimeVersion string
}

type MessageChannel struct {
	c      chan MessageHeader
	closed bool
	mutex  sync.Mutex
}
