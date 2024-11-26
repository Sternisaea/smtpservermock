package smtpservermock

type cmdVRFY struct{}

func (c *cmdVRFY) getPrefix() string {
	return "VRFY"
}

func (c *cmdVRFY) execute(t *transmission, arg string) error {
	resp := "252  Cannot VRFY user, but will accept message and attempt delivery"
	return (*t).writeResponse(resp)
}

// User Name <local-part@domain>
// local-part@domain

// "553 User ambiguous"
// "550 Requested action not taken: mailbox unavailable"
