package smtpservermock

type CmdVRFY struct{}

func (c *CmdVRFY) GetPrefix() string {
	return "VRFY"
}

func (c *CmdVRFY) Execute(t *Transmission, arg string) error {
	resp := "252  Cannot VRFY user, but will accept message and attempt delivery"
	return (*t).WriteResponse(resp)
}

// User Name <local-part@domain>
// local-part@domain

// "553 User ambiguous"
// "550 Requested action not taken: mailbox unavailable"
