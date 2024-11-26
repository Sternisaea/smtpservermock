package smtpservermock

import "log"

type cmdNOOP struct{}

func (c *cmdNOOP) getPrefix() string {
	return "NOOP"
}

func (c *cmdNOOP) execute(t *transmission, arg string) error {
	if arg != "" {
		log.Printf("%#v", arg)
		return (*t).writeResponse("501 Syntax error in parameters or arguments")
	}
	return (*t).writeResponse("250 OK")
}
