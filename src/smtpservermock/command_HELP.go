package smtpservermock

type cmdHELP struct{}

func (c *cmdHELP) getPrefix() string {
	return "HELP"
}

func (c *cmdHELP) execute(t *transmission, arg string) error {
	resp := "250 Supported commands:"
	for _, c := range (*t).commands {
		resp += "  " + c.getPrefix()
	}
	return (*t).writeResponse(resp)
}
