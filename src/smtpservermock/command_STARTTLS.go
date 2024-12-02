package smtpservermock

import (
	"bufio"
	"crypto/tls"
)

type cmdSTARTTLS struct{}

func (c *cmdSTARTTLS) getPrefix() string {
	return "STARTTLS"
}

func (c *cmdSTARTTLS) execute(t *transmission, arg string) error {
	if (*t).starttlsActive {
		return (*t).writeResponse("503 Bad sequence of commands")
	}
	if (*t).connType != ehloType && (*t).connType != noType {
		return (*t).writeResponse("500 Syntax error, command unrecognized")
	}
	if (*t).starttlsConfig == nil {
		return (*t).writeResponse("501 Syntax error (no parameters allowed)")
	}

	if err := (*t).writeResponse("220 Ready to start TLS"); err != nil {
		return err
	}
	tlsConn := tls.Server((*t).netConnection, (*t).starttlsConfig)
	(*t).reader = bufio.NewReader(tlsConn)
	(*t).writer = bufio.NewWriter(tlsConn)
	(*t).clientName = ""
	(*t).connType = ehloType
	(*t).initCurrentMessage()
	(*t).setCommands()
	(*t).starttlsActive = true
	return nil
}
