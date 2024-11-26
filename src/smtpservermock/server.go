package smtpservermock

import (
	"crypto/tls"
	"errors"
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
		log.Printf("Connection error: %s", err)
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

func (s *SmtpServer) GetConnectionMessages() ConnectionMessages {
	return (*s).connectionMessages
}
