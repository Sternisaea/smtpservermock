package smtpservermock

type cmdQuit struct{}

func (c *cmdQuit) getPrefix() string {
	return "QUIT"
}

func (c *cmdQuit) execute(t *transmission, arg string) error {
	(*t).connType = quitType
	resp := "221 " + (*t).serverName + " Goodby " + (*t).clientName
	return (*t).writeResponse(resp)
}
