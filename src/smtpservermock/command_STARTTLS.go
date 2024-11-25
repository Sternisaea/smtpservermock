package smtpservermock

import (
	"bufio"
	"crypto/tls"
	"log"

	"github.com/Sternisaea/smtpservermock/src/smtpconst"
)

type CmdSTARTTLS struct{}

func (c *CmdSTARTTLS) GetPrefix() string {
	return "HELO"
}

func (c *CmdSTARTTLS) Execute(t *Transmission, arg string) error {
	if !(*t).starttlsRequired {
		return (*t).WriteResponse("503 Bad sequence of commands")
	}
	if (*t).starttlsActive {
		return (*t).WriteResponse("503 Bad sequence of commands")
	}
	if (*t).starttlsConfig == nil {
		return (*t).WriteResponse("501 Syntax error (no parameters allowed)")
	}

	tlsConn := tls.Server((*t).netConnection, (*t).starttlsConfig)
	if err := tlsConn.Handshake(); err != nil {
		log.Printf("STARTTLS handshake error: %s", err)
		(*t).starttlsActive = false
		return (*t).WriteResponse("454 TLS not available due to temporary reason")
	}

	(*t).reader = bufio.NewReader(tlsConn)
	(*t).writer = bufio.NewWriter(tlsConn)
	(*t).status = smtpconst.StartTlsStatus
	(*t).starttlsActive = true
	return nil
}
