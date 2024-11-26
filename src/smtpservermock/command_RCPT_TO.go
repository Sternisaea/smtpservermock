package smtpservermock

import (
	"strings"
)

type cmdRCPTTO struct{}

func (c *cmdRCPTTO) getPrefix() string {
	return "RCPT TO"
}

func (c *cmdRCPTTO) execute(t *transmission, arg string) error {
	if (*t).msgStatus != mailFromMessage && (*t).msgStatus != receiptToMessage {
		return (*t).writeResponse("503 Bad sequence of commands")
	}

	texts := textAngleBracketsRegex.FindAllString(arg, -1)
	if texts == nil || len(texts) > 1 || arg[0] != ':' {
		return (*t).writeResponse("501 Syntax error. Format: RCPT TO: <mailbox>")
	}
	email := strings.TrimPrefix(strings.TrimSuffix(texts[0], ">"), "<")
	// Optional: check email
	(*t).currentMessage.to = append((*t).currentMessage.to, email)
	(*t).msgStatus = receiptToMessage
	return (*t).writeResponse("250 OK")
}

// RCPT TO:<recipient1@example.com>
