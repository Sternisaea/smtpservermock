package smtpservermock

import (
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net"

	"github.com/Sternisaea/smtpservermock/src/smtpconst"
)

type SmtpConnection interface {
	SetupListener() (net.Listener, error)
	ShutdownListener(listener net.Listener) error
}

func getSmtpConnection(sec smtpconst.Security, servername, addr string, tlsconfig *tls.Config) (SmtpConnection, error) {
	switch sec {
	case smtpconst.NoSecurity, smtpconst.StartTlsSec:
		return NewRegularConnection(servername, addr)
	case smtpconst.SslTlsSec:
		return NewTLSConnection(servername, addr, tlsconfig)
	default:
		return nil, fmt.Errorf("unknown security type %s", sec)
	}
}

type RegularConnection struct {
	servername string
	addr       string
}

func NewRegularConnection(servername, addr string) (*RegularConnection, error) {
	return &RegularConnection{servername: servername, addr: addr}, nil
}

func (c *RegularConnection) SetupListener() (net.Listener, error) {
	listener, err := net.Listen("tcp", (*c).addr)
	if err != nil {
		return nil, fmt.Errorf("error starting server: %w", err)
	}
	log.Printf("%s listening on %s", (*c).servername, (*c).addr)
	return listener, nil

}
func (c *RegularConnection) ShutdownListener(listener net.Listener) error {
	if listener == nil {
		return errors.New("server not running")
	}
	if err := listener.Close(); err != nil {
		return err
	}
	log.Printf("%s has been shut down", (*c).servername)
	return nil
}

type TLSConnection struct {
	servername string
	addr       string
	tlsConfig  *tls.Config
}

func NewTLSConnection(servername, addr string, tlsconfig *tls.Config) (*TLSConnection, error) {
	return &TLSConnection{servername: servername, addr: addr, tlsConfig: tlsconfig}, nil
}

func (c *TLSConnection) SetupListener() (net.Listener, error) {
	listener, err := tls.Listen("tcp", (*c).addr, (*c).tlsConfig)
	if err != nil {
		return nil, fmt.Errorf("error starting server: %w", err)
	}
	log.Printf("%s listening on %s", (*c).servername, (*c).addr)
	return listener, nil
}

func (c *TLSConnection) ShutdownListener(listener net.Listener) error {
	if listener == nil {
		return errors.New("server not running")
	}
	if err := listener.Close(); err != nil {
		return err
	}
	log.Printf("%s has been shut down", (*c).servername)
	return nil
}
