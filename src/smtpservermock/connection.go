package smtpservermock

import (
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net"

	"github.com/Sternisaea/smtpservermock/src/smtpconst"
)

type smtpConnection interface {
	setupListener() (net.Listener, error)
	shutdownListener(listener net.Listener) error
}

func getSmtpConnection(sec smtpconst.Security, servername, addr string, tlsconfig *tls.Config) (smtpConnection, error) {
	switch sec {
	case smtpconst.NoSecurity, smtpconst.StartTlsSec:
		return newRegularConnection(servername, addr)
	case smtpconst.SslTlsSec:
		return newTLSConnection(servername, addr, tlsconfig)
	default:
		return nil, fmt.Errorf("unknown security type %s", sec)
	}
}

type regularConnection struct {
	servername string
	addr       string
}

func newRegularConnection(servername, addr string) (*regularConnection, error) {
	return &regularConnection{servername: servername, addr: addr}, nil
}

func (c *regularConnection) setupListener() (net.Listener, error) {
	listener, err := net.Listen("tcp", (*c).addr)
	if err != nil {
		return nil, fmt.Errorf("error starting server: %w", err)
	}
	log.Printf("%s listening on %s", (*c).servername, (*c).addr)
	return listener, nil

}
func (c *regularConnection) shutdownListener(listener net.Listener) error {
	if listener == nil {
		return errors.New("server not running")
	}
	if err := listener.Close(); err != nil {
		return err
	}
	log.Printf("%s has been shut down", (*c).servername)
	return nil
}

type tlsConnection struct {
	servername string
	addr       string
	tlsConfig  *tls.Config
}

func newTLSConnection(servername, addr string, tlsconfig *tls.Config) (*tlsConnection, error) {
	return &tlsConnection{servername: servername, addr: addr, tlsConfig: tlsconfig}, nil
}

func (c *tlsConnection) setupListener() (net.Listener, error) {
	listener, err := tls.Listen("tcp", (*c).addr, (*c).tlsConfig)
	if err != nil {
		return nil, fmt.Errorf("error starting server: %w", err)
	}
	log.Printf("%s listening on %s", (*c).servername, (*c).addr)
	return listener, nil
}

func (c *tlsConnection) shutdownListener(listener net.Listener) error {
	if listener == nil {
		return errors.New("server not running")
	}
	if err := listener.Close(); err != nil {
		return err
	}
	log.Printf("%s has been shut down", (*c).servername)
	return nil
}
