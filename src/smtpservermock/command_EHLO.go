package smtpservermock

type CmdEHLO struct{}

func (c *CmdEHLO) GetPrefix() string {
	return "EHLO"
}

func (c *CmdEHLO) Execute(t *Transmission, arg string) error {
	(*t).clientName = arg
	(*t).connType = EhloType
	(*t).initCurrentMessage()
	(*t).setCommands()
	resp := "250 " + (*t).serverName + " Hello " + (*t).clientName
	return (*t).WriteResponse(resp)
}
