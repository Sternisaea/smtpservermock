package smtpservermock

type CmdDATA struct{}

func (c *CmdDATA) GetPrefix() string {
	return "DATA"
}

func (c *CmdDATA) Execute(t *Transmission, arg string) error {
	return nil
}
