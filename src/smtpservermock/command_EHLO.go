package smtpservermock

type cmdEHLO struct{}

func (c *cmdEHLO) getPrefix() string {
	return "EHLO"
}

func (c *cmdEHLO) execute(t *transmission, arg string) error {
	(*t).clientName = arg
	(*t).connType = ehloType
	(*t).initCurrentMessage()
	(*t).setCommands()
	resp := "250 " + (*t).serverName + " Hello " + (*t).clientName
	return (*t).writeResponse(resp)
}
