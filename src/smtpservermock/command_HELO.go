package smtpservermock

type CmdHELO struct{}

func (c *CmdHELO) GetPrefix() string {
	return "HELO"
}

func (c *CmdHELO) Execute(t *Transmission, arg string) error {
	if (*t).starttlsRequired && !(*t).starttlsActive {
		return (*t).WriteResponse("530 Must issue a STARTTLS command first")
	}
	(*t).clientName = arg
	(*t).connType = HeloType
	(*t).initCurrentMessage()
	resp := "250 " + (*t).serverName + " Hello " + (*t).clientName
	return (*t).WriteResponse(resp)
}
