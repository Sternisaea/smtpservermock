package smtpservermock

import "github.com/Sternisaea/smtpservermock/src/smtpconst"

type cmdEHLO struct{}

func (c *cmdEHLO) getPrefix() string {
	return "EHLO"
}

func (c *cmdEHLO) execute(t *transmission, arg string) error {
	(*t).clientName = arg
	(*t).connType = ehloType
	(*t).initCurrentMessage()
	(*t).setCommands()

	rm := make([]string, 0, 6)
	rm = append(rm, (*t).serverName+" says hello "+(*t).clientName)
	if (*t).security == smtpconst.StartTlsSec && !(*t).starttlsActive {
		rm = append(rm, "STARTTLS")
	}

	var resp string
	for i, r := range rm {
		if i == len(rm)-1 {
			resp = "250 " + r
		} else {
			resp = "250-" + r
		}
		if err := (*t).writeResponse(resp); err != nil {
			return err
		}
	}
	return nil
}

// EHLO
// 250-Nice to meet you.
// 250-8BITMIME
// 250-SIZE
// 250-SMTPUTF8
// 250-STARTTLS
// 250-AUTH=CRAM-MD5 PLAIN LOGIN
// 250 AUTH CRAM-MD5 PLAIN LOGIN
