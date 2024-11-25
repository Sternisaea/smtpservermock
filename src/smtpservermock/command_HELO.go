package smtpservermock

import "github.com/Sternisaea/smtpservermock/src/smtpconst"

type CmdHELO struct{}

func (c *CmdHELO) GetPrefix() string {
	return "HELO"
}

func (c *CmdHELO) Execute(t *Transmission, arg string) error {
	(*t).clientName = arg
	(*t).status = smtpconst.HeloStatus
	resp := "250 " + (*t).serverName + " Hello " + (*t).clientName
	return (*t).WriteResponse(resp)
}
