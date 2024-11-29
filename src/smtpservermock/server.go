package smtpservermock

import (
	"crypto/tls"
	"errors"
	"io"
	"log"
	"net"

	"github.com/Sternisaea/smtpservermock/src/smtpconst"
	"github.com/google/uuid"
)

type SmtpServer struct {
	name       string
	security   smtpconst.Security
	address    string
	connection smtpConnection
	tlsconfig  *tls.Config

	listener           net.Listener
	connectionMessages ConnectionMessages
}

type ConnectionMessages map[string][]CompletedMessage

type CompletedMessage struct {
	From string
	To   []string
	Data string
}

type connMessage struct {
	id      string
	message message
}

// NewSmtpServer creates a new instance of SmtpServer
// - sec smtpconst.Security = type of security (e.g. No security, SSL-TLS, STARTTLS)
// - servername string      = name of the server
// - certFile string        = path to PEM encoded public key (leave empty if no security)
// - keyFile string         = path to PEM encoded privat key (leave empty if no security)
// An error is returned for an unknown security type or invalid keys
func NewSmtpServer(sec smtpconst.Security, servername, addr, certFile, keyFile string) (*SmtpServer, error) {
	tlsconfig, err := getTLSConfig(sec, certFile, keyFile)
	if err != nil {
		return nil, err
	}
	smtpconn, err := getSmtpConnection(sec, servername, addr, tlsconfig)
	if err != nil {
		return nil, err
	}
	conMsgs := make(ConnectionMessages)
	return &SmtpServer{name: servername, security: sec, address: addr, connection: smtpconn, tlsconfig: tlsconfig, connectionMessages: conMsgs}, nil
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

// ListenAndServe starts the SMTP server and begins listening for incoming connections.
// An error is returned if there was an issue setting up the listener.
func (s *SmtpServer) ListenAndServe() error {
	var err error
	(*s).listener, err = (*s).connection.setupListener()
	if err != nil {
		return err
	}

	msgCh := make(chan connMessage)
	go s.handleMessages(msgCh)
	go s.listening(msgCh)
	return nil
}

func (s *SmtpServer) listening(msgCh chan<- connMessage) {
	defer close(msgCh)
	for {
		conn, err := (*s).listener.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				return
			}
			log.Printf("Connection error: %s", err)
			return
		}
		go s.handle(conn, msgCh)
	}
}

// Shutdown gracefully shuts down the SMTP server
// It returns an error if there was an issue shutting down the server.
func (s *SmtpServer) Shutdown() error {
	return (*s).connection.shutdownListener((*s).listener)
}

func (s *SmtpServer) handle(conn net.Conn, msgCh chan<- connMessage) {
	defer conn.Close()

	id := uuid.New().String()
	trsm := newTransmission((*s).security, conn, (*s).name, id, msgCh)
	if (*s).security == smtpconst.StartTlsSec {
		trsm.SetStartTLSConfig((*s).tlsconfig)
	}
	if err := trsm.Process(); err != nil {
		if err == io.EOF {
			log.Printf("Connection closed by client (EOF)")
		} else {
			log.Printf("Connection error: %s", err)
		}
	}
}

func (s *SmtpServer) handleMessages(msgCh <-chan connMessage) {
	for connMsg := range msgCh {
		msg := CompletedMessage{
			From: connMsg.message.from,
			To:   connMsg.message.to,
			Data: connMsg.message.data,
		}
		(*s).connectionMessages[connMsg.id] = append((*s).connectionMessages[connMsg.id], msg)
	}
}

// GetConnectionMessages returns the e-mail messages received by the SMTP Server
// For every connection an unique GUID is created
func (s *SmtpServer) GetConnectionMessages() ConnectionMessages {
	return (*s).connectionMessages
}
