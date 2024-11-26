package smtpservermock

import "log"

type CmdNOOP struct{}

func (c *CmdNOOP) GetPrefix() string {
	return "NOOP"
}

func (c *CmdNOOP) Execute(t *Transmission, arg string) error {
	if arg != "" {
		log.Printf("%#v", arg)
		return (*t).WriteResponse("501 Syntax error in parameters or arguments")
	}
	return (*t).WriteResponse("250 OK")
}
