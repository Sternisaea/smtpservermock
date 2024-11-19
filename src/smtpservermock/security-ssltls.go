package smtpservermock

import (
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net"
)

type SecurityTLS struct {
	addr      string
	tlsConfig *tls.Config
}

func NewSecurityTLS(addr, certFile, keyFile string) (*SecurityTLS, error) {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}
	config := &tls.Config{Certificates: []tls.Certificate{cert}}
	return &SecurityTLS{addr: addr, tlsConfig: config}, nil
}

func (s *SecurityTLS) SetupListener() (net.Listener, error) {
	listener, err := tls.Listen("tcp", (*s).addr, (*s).tlsConfig)
	if err != nil {
		return nil, fmt.Errorf("error starting server: %w", err)
	}
	log.Printf("Mock SMTP server listening on %s", (*s).addr)
	return listener, nil
}

func (s *SecurityTLS) ShutdownListener(listener net.Listener) error {
	if listener == nil {
		return errors.New("server not running")
	}
	if err := listener.Close(); err != nil {
		return err
	}
	log.Printf("Mock SMTP server has been shut down")
	return nil
}
