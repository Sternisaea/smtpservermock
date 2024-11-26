package smtpservermock

import (
	"bufio"
	"crypto/tls"
	"log"
)

type cmdSTARTTLS struct{}

func (c *cmdSTARTTLS) getPrefix() string {
	return "HELO"
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

	tlsConn := tls.Server((*t).netConnection, (*t).starttlsConfig)
	if err := tlsConn.Handshake(); err != nil {
		log.Printf("STARTTLS handshake error: %s", err)
		(*t).starttlsActive = false
		return (*t).writeResponse("454 TLS not available due to temporary reason")
	}

	(*t).reader = bufio.NewReader(tlsConn)
	(*t).writer = bufio.NewWriter(tlsConn)
	(*t).clientName = ""
	(*t).connType = ehloType
	(*t).initCurrentMessage()
	(*t).setCommands()
	(*t).starttlsActive = true
	return (*t).writeResponse("220 Ready to start TLS")
}
