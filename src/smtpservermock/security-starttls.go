package smtpservermock

import (
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net"
)

type SecuritySTARTTLS struct {
	addr      string
	tlsConfig *tls.Config
}

func NewSecurritySTARTTLS(addr, certFile, keyFile string) (*SecuritySTARTTLS, error) {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}
	config := &tls.Config{Certificates: []tls.Certificate{cert}}
	return &SecuritySTARTTLS{addr: addr, tlsConfig: config}, nil
}

func (s *SecuritySTARTTLS) SetupListener() (net.Listener, error) {
	listener, err := net.Listen("tcp", (*s).addr)
	if err != nil {
		return nil, fmt.Errorf("error starting server: %w", err)
	}
	log.Printf("Mock SMTP server listening on %s", (*s).addr)

	// Handle command STARTTLS

	return listener, nil
}

func (s *SecuritySTARTTLS) ShutdownListener(listener net.Listener) error {
	if listener == nil {
		return errors.New("server not running")
	}
	if err := listener.Close(); err != nil {
		return err
	}
	log.Printf("Mock SMTP server has been shut down")
	return nil
}
