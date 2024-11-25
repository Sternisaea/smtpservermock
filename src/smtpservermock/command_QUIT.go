package smtpservermock

type CmdQuit struct{}

func (c *CmdQuit) GetPrefix() string {
	return "QUIT"
}

func (c *CmdQuit) Execute(t *Transmission, arg string) error {
	(*t).connType = QuitType
	resp := "221 " + (*t).serverName + " Goodby " + (*t).clientName
	return (*t).WriteResponse(resp)
}
