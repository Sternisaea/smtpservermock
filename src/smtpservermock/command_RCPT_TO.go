package smtpservermock

import (
	"strings"
)

type CmdRCPTTO struct{}

func (c *CmdRCPTTO) GetPrefix() string {
	return "RCPT TO"
}

func (c *CmdRCPTTO) Execute(t *Transmission, arg string) error {
	if (*t).msgStatus != MailFromMessage && (*t).msgStatus != ReceiptToMessage {
		return (*t).WriteResponse("503 Bad sequence of commands")
	}

	texts := textAngleBracketsRegex.FindAllString(arg, -1)
	if texts == nil || len(texts) > 1 || arg[0] != ':' {
		return (*t).WriteResponse("501 Syntax error. Format: RCPT TO: <mailbox>")
	}
	email := strings.TrimPrefix(strings.TrimSuffix(texts[0], ">"), "<")
	// Optional: check email
	(*t).currentMessage.To = append((*t).currentMessage.To, email)
	(*t).msgStatus = ReceiptToMessage
	return (*t).WriteResponse("250 OK")
}

// RCPT TO:<recipient1@example.com>
