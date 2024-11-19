package smtpservermock

import (
	"errors"
	"fmt"
	"log"
	"net"
)

type SecurityNone struct {
	addr string
}

func NewSecurityNone(addr string) (*SecurityNone, error) {
	return &SecurityNone{addr: addr}, nil
}

func (s *SecurityNone) SetupListener() (net.Listener, error) {
	listener, err := net.Listen("tcp", (*s).addr)
	if err != nil {
		return nil, fmt.Errorf("error starting server: %w", err)
	}
	log.Printf("Mock SMTP server listening on %s", (*s).addr)
	return listener, nil
}

func (s *SecurityNone) ShutdownListener(listener net.Listener) error {
	if listener == nil {
		return errors.New("server not running")
	}
	if err := listener.Close(); err != nil {
		return err
	}
	log.Printf("Mock SMTP server has been shut down")
	return nil
}
