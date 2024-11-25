package smtpservermock

type CmdHELP struct{}

func (c *CmdHELP) GetPrefix() string {
	return "HELP"
}

func (c *CmdHELP) Execute(t *Transmission, arg string) error {
	// return (*t).WriteResponse("250 OK")

	// (*t).starttlsActive == true >> Do not display STARTTLS

	return nil
}
