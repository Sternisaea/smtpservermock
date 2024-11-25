package smtpservermock

type CmdRCPTTO struct{}

func (c *CmdRCPTTO) GetPrefix() string {
	return "MAIL FROM"
}

func (c *CmdRCPTTO) Execute(t *Transmission, arg string) error {
	return nil
}
