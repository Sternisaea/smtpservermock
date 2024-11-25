package smtpservermock

import (
	"strings"
)

type CmdRCPTTO struct{}

func (c *CmdRCPTTO) GetPrefix() string {
	return "MAIL FROM"
}

func (c *CmdRCPTTO) Execute(t *Transmission, arg string) error {
	if (*t).msgStatus != MailFromMessage && (*t).msgStatus != ReceiptToMessage {
		return (*t).WriteResponse("503 Bad sequence of commands")
	}
	if !emailAngleBracketsRegex.MatchString(arg) {
		return (*t).WriteResponse("501 Syntax error in parameters or arguments")
	}
	email := strings.TrimPrefix(strings.TrimSuffix(arg, ">"), "<")

	// Optional: check email

	(*t).currentMessage.To = append((*t).currentMessage.To, email)
	(*t).msgStatus = ReceiptToMessage
	return (*t).WriteResponse("250 OK")
}

// RCPT TO:<recipient1@example.com>
