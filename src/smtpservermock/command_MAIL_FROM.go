package smtpservermock

type CmdMAILFROM struct{}

func (c *CmdMAILFROM) GetPrefix() string {
	return "MAIL FROM"
}

func (c *CmdMAILFROM) Execute(t *Transmission, arg string) error {
	return nil
}
