package smtpservermock

type CmdRSET struct{}

func (c *CmdRSET) GetPrefix() string {
	return "RSET"
}

func (c *CmdRSET) Execute(t *Transmission, arg string) error {
	(*t).initCurrentMessage()
	return (*t).WriteResponse("250 OK")
}
