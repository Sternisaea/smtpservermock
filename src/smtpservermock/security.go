package smtpservermock

import (
	"errors"
	"net"

	"github.com/Sternisaea/smtpservermock/src/smtpconst"
)

type SmtpSecurity interface {
	SetupListener() (net.Listener, error)
	ShutdownListener(listener net.Listener) error
}

var (
	ErrUnknownSecurity = errors.New("unknown type ")
)

func GetSmtpSecurity(sec smtpconst.Security, addr, certFile, keyFile string) (SmtpSecurity, error) {
	switch sec {
	case smtpconst.NoSecurity:
		return NewSecurityNone(addr)
	case smtpconst.StartTlsSec:
		return NewSecurritySTARTTLS(addr, certFile, keyFile)
	case smtpconst.SslTlsSec:
		return NewSecurityTLS(addr, certFile, keyFile)
	default:
		return nil, ErrUnknownSecurity
	}
}
