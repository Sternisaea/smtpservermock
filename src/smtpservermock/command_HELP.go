package smtpservermock

type CmdHELP struct{}

func (c *CmdHELP) GetPrefix() string {
	return "HELP"
}

func (c *CmdHELP) Execute(t *Transmission, arg string) error {
	resp := "250 Supported commands:"
	for _, c := range (*t).commands {
		resp += "  " + c.GetPrefix()
	}
	return (*t).WriteResponse(resp)
}
