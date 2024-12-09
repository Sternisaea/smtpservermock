package smtpservermock

type cmdRSET struct{}

func (c *cmdRSET) getPrefix() string {
	return "RSET"
}

func (c *cmdRSET) execute(t *transmission, arg string) error {
	err := (*t).writeResponse("250 OK")
	(*t).initCurrentMessage()
	return err
}
