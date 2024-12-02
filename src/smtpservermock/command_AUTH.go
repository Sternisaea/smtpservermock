package smtpservermock

import (
	"strings"

	"github.com/Sternisaea/smtpservermock/src/smtpconst"
)

type cmdAUTH struct{}

func (c *cmdAUTH) getPrefix() string {
	return "AUTH"
}

func (c *cmdAUTH) execute(t *transmission, arg string) error {

	if arg == "" {
		return (*t).writeResponse("501 Syntax error in parameters or arguments")
	}

	authText := ""
	if words := strings.Fields(arg); len(words) > 0 {
		authText = words[0]
		// arg = strings.Join(words[1:], " ")
	}

	switch smtpconst.AuthenticationMethod(strings.ToLower(authText)) {
	case smtpconst.PlainAuth:
		return (*t).writeResponse("235 Authentication successful")
	default:
		return (*t).writeResponse("504 Unrecognized authentication type")
	}
}

// AUTH PLAIN AG1lbXlzZWxmAHZlcnl1bmtub3du
