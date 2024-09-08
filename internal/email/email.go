package email

import "bytes"

type Email interface {
	SendEmail(to string, subject string, body string, attachement *bytes.Buffer, attachmentName string) error
}
