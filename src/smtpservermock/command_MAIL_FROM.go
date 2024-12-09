package smtpservermock

import (
	"strings"
)

type cmdMAILFROM struct{}

func (c *cmdMAILFROM) getPrefix() string {
	return "MAIL FROM"
}

func (c *cmdMAILFROM) execute(t *transmission, arg string) error {
	if (*t).connType != heloType && (*t).connType != ehloType {
		return (*t).writeResponse("503 Bad sequence of commands")
	}
	if (*t).security == StartTlsSec && !(*t).starttlsActive {
		return (*t).writeResponse("530 Must issue a STARTTLS command first")
	}
	if (*t).msgStatus != emptyMessage {
		return (*t).writeResponse("503 Bad sequence of commands")
	}

	texts := textAngleBracketsRegex.FindAllString(arg, -1)
	if texts == nil || arg[0] != ':' {
		return (*t).writeResponse("501 Syntax error. Format: MAIL FROM: <mailbox>")
	}
	email := strings.TrimPrefix(strings.TrimSuffix(texts[0], ">"), "<")
	// Optional: check email
	(*t).currentMessage.From = email
	(*t).msgStatus = mailFromMessage
	return (*t).writeResponse("250 OK")
}

// MAIL FROM:<sender@example.com>

// "550 Requested action not taken: mailbox unavailable"
// "451 Requested action aborted: local error in processing"
