package smtpservermock

import "github.com/Sternisaea/smtpservermock/src/smtpconst"

type CmdEHLO struct{}

func (c *CmdEHLO) GetPrefix() string {
	return "EHLO"
}

func (c *CmdEHLO) Execute(t *Transmission, arg string) error {
	(*t).clientName = arg
	(*t).status = smtpconst.EhloStatus
	resp := "250 " + (*t).serverName + " Hello " + (*t).clientName
	return (*t).WriteResponse(resp)
}
