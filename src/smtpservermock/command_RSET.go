package smtpservermock

type cmdRSET struct{}

func (c *cmdRSET) getPrefix() string {
	return "RSET"
}

func (c *cmdRSET) execute(t *transmission, arg string) error {
	(*t).initCurrentMessage()
	return (*t).writeResponse("250 OK")
}
