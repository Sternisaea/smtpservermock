package smtpservermock

import (
	"strings"
)

type CmdMAILFROM struct{}

func (c *CmdMAILFROM) GetPrefix() string {
	return "MAIL FROM"
}

func (c *CmdMAILFROM) Execute(t *Transmission, arg string) error {
	if (*t).connType != HeloType && (*t).connType != EhloType {
		return (*t).WriteResponse("503 Bad sequence of commands")
	}
	if (*t).starttlsRequired && !(*t).starttlsActive {
		return (*t).WriteResponse("530 Must issue a STARTTLS command first")
	}
	if (*t).msgStatus != EmptyMessage {
		return (*t).WriteResponse("503 Bad sequence of commands")
	}
	if !emailAngleBracketsRegex.MatchString(arg) {
		return (*t).WriteResponse("501 Syntax error in parameters or arguments")
	}
	email := strings.TrimPrefix(strings.TrimSuffix(arg, ">"), "<")

	// Optional: check email

	(*t).currentMessage.From = email
	(*t).msgStatus = MailFromMessage
	return (*t).WriteResponse("250 OK")
}

// MAIL FROM:<sender@example.com>

// "550 Requested action not taken: mailbox unavailable"
// "451 Requested action aborted: local error in processing"
