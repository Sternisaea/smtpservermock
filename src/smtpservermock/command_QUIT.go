package smtpservermock

import "github.com/Sternisaea/smtpservermock/src/smtpconst"

type CmdQuit struct{}

func (c *CmdQuit) GetPrefix() string {
	return "QUIT"
}

func (c *CmdQuit) Execute(t *Transmission, arg string) error {
	(*t).status = smtpconst.QuitStatus
	resp := "221 " + (*t).serverName + " Goodby " + (*t).clientName
	return (*t).WriteResponse(resp)
}
