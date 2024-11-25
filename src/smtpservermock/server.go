package smtpservermock

import (
	"crypto/tls"
	"errors"
	"log"
	"net"

	"github.com/Sternisaea/smtpservermock/src/smtpconst"
)

type SmtpServer struct {
	name       string
	security   smtpconst.Security
	address    string
	connection SmtpConnection
	tlsconfig  *tls.Config

	listener net.Listener
}

func NewSmtpServer(sec smtpconst.Security, servername, addr, certFile, keyFile string) (*SmtpServer, error) {
	tlsconfig, err := getTLSConfig(sec, certFile, keyFile)
	if err != nil {
		return nil, err
	}
	smtpconn, err := getSmtpConnection(sec, servername, addr, tlsconfig)
	if err != nil {
		return nil, err
	}
	return &SmtpServer{name: servername, security: sec, address: addr, connection: smtpconn, tlsconfig: tlsconfig}, nil
}

func getTLSConfig(sec smtpconst.Security, certFile, keyFile string) (*tls.Config, error) {
	switch sec {
	case smtpconst.StartTlsSec, smtpconst.SslTlsSec:
		cert, err := tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			return nil, err
		}
		return &tls.Config{Certificates: []tls.Certificate{cert}}, nil
	default:
		return nil, nil
	}
}

func (s *SmtpServer) ListenAndServe() error {
	var err error
	(*s).listener, err = (*s).connection.SetupListener()
	if err != nil {
		return err
	}
	go s.listening()
	return nil
}

func (s *SmtpServer) listening() {
	for {
		conn, err := (*s).listener.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				return
			}
			log.Printf("Connection error: %s", err)
			return
		}
		go s.handle(conn)
	}
}

func (s *SmtpServer) Shutdown() error {
	return (*s).connection.ShutdownListener((*s).listener)
}

func (s *SmtpServer) handle(conn net.Conn) {
	defer conn.Close()

	trsm := NewTransmission(conn, (*s).name)
	if (*s).security == smtpconst.StartTlsSec {
		trsm.SetStartTLSConfig((*s).tlsconfig)
	}
	trsm.SetCommands([]Command{&CmdEHLO{}, &CmdHELO{}, &CmdQuit{}, &CmdNOOP{}, &CmdMAILFROM{}, &CmdRCPTTO{}, &CmdSTARTTLS{}})
	if err := trsm.Process(); err != nil {
		log.Printf("Connection error: %s", err)
	}
}
