package smtpservermock

type cmdHELO struct{}

func (c *cmdHELO) getPrefix() string {
	return "HELO"
}

func (c *cmdHELO) execute(t *transmission, arg string) error {
	if (*t).security == StartTlsSec && !(*t).starttlsActive {
		return (*t).writeResponse("530 Must issue a STARTTLS command first")
	}
	(*t).clientName = arg
	(*t).connType = heloType
	(*t).initCurrentMessage()
	(*t).setCommands()
	resp := "250 " + (*t).serverName + " Hello " + (*t).clientName
	return (*t).writeResponse(resp)
}
