package smtpservermock

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"
)

type SmtpServer struct {
	name       string
	security   Security
	address    string
	connection smtpConnection
	tlsconfig  *tls.Config

	listener  net.Listener
	rawTextCh chan transmissionRawLines
	messageCh chan transmissionMessage

	lockConnectResults sync.Mutex
	connectionResults  map[string][]Result
}

var (
	timeout = 2000 * time.Millisecond

	ErrUnknownConnectionAddress  = errors.New("unknown address")
	ErrUnknownConnectionSequence = errors.New("unkown connection sequence")
	ErrUnkownMessageSequence     = errors.New("unknown message sequence")
	ErrTimeout                   = errors.New("timeout waiting for clear channel buffer")
)

// NewSmtpServer creates a new instance of SmtpServer. Parameters are the type of security (e.g. No security, SSL-TLS, STARTTLS), a servername,
// the server address (mail.example.com:587), path to PEM encoded public key and path to PEM encoded privat key.
// Leave the paths empty if no security is applied. An error is returned for an unknown security type or invalid keys
func NewSmtpServer(sec Security, servername, addr, certFile, keyFile string) (*SmtpServer, error) {
	tlsconfig, err := getTLSConfig(sec, certFile, keyFile)
	if err != nil {
		return nil, err
	}
	smtpconn, err := getSmtpConnection(sec, servername, addr, tlsconfig)
	if err != nil {
		return nil, err
	}
	return &SmtpServer{name: servername, security: sec, address: addr, connection: smtpconn, tlsconfig: tlsconfig, connectionResults: make(map[string][]Result)}, nil
}

func getTLSConfig(sec Security, certFile, keyFile string) (*tls.Config, error) {
	switch sec {
	case StartTlsSec, SslTlsSec:
		cert, err := tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			return nil, err
		}
		return &tls.Config{
			Certificates: []tls.Certificate{cert},
		}, nil
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

	(*s).rawTextCh = make(chan transmissionRawLines)
	(*s).messageCh = make(chan transmissionMessage)
	go func() {
		for trsRaw := range (*s).rawTextCh {
			(*s).lockConnectResults.Lock()
			(*s).connectionResults[trsRaw.address][trsRaw.entryNo-1].Raw = append((*s).connectionResults[trsRaw.address][trsRaw.entryNo-1].Raw, trsRaw.lines...)
			(*s).lockConnectResults.Unlock()
		}
	}()

	go func() {
		for trsMsg := range (*s).messageCh {
			(*s).lockConnectResults.Lock()
			(*s).connectionResults[trsMsg.address][trsMsg.entryNo-1].Messages = append((*s).connectionResults[trsMsg.address][trsMsg.entryNo-1].Messages, trsMsg.message)
			(*s).lockConnectResults.Unlock()
		}
	}()
	go s.listening((*s).rawTextCh, (*s).messageCh)
	return nil
}

func (s *SmtpServer) listening(rawCh chan<- transmissionRawLines, msgCh chan<- transmissionMessage) {
	defer close(rawCh)
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

		addr := conn.RemoteAddr().String()
		(*s).lockConnectResults.Lock()
		entryNo := len((*s).connectionResults[addr]) + 1
		(*s).connectionResults[addr] = append((*s).connectionResults[addr], Result{EntryNo: entryNo, Start: time.Now()})
		(*s).lockConnectResults.Unlock()
		go s.handle(conn, addr, entryNo, rawCh, msgCh)
	}
}

func (s *SmtpServer) handle(conn net.Conn, address string, entryNo int, rawCh chan<- transmissionRawLines, msgCh chan<- transmissionMessage) {
	defer conn.Close()
	trm := newTransmission(address, entryNo, (*s).security, conn, (*s).name, rawCh, msgCh)
	if (*s).security == StartTlsSec {
		trm.SetStartTLSConfig((*s).tlsconfig)
	}
	if err := trm.Process(); err != nil {
		if err != io.EOF {
			log.Printf("Connection error: %s", err)
		}
	}
	(*s).lockConnectResults.Lock()
	(*s).connectionResults[address][entryNo-1].End = time.Now()
	(*s).lockConnectResults.Unlock()
}

// Shutdown gracefully shuts down the SMTP server
// It returns an error if there was an issue shutting down the server.
func (s *SmtpServer) Shutdown() error {
	return (*s).connection.shutdownListener((*s).listener)
}

// GetConnectionAddresses returns all connection addresses that made a connection to the SMTP server.
func (s *SmtpServer) GetConnectionAddresses() ([]string, error) {
	addrs := make([]string, 0, len((*s).connectionResults))
	for a := range (*s).connectionResults {
		addrs = append(addrs, a)
	}
	return addrs, nil
}

// GetResultMessage returns the mail message received by the SMTP server for a given connection address, connection sequence number
// and message sequence number. The connection sequence number is usually 1, but can be increased if a subsequent connection would use
// the same TCP port (which is very unlikely). The message sequence number starts with 1 and is increased by 1 for every new message
// within the same connection.
func (s *SmtpServer) GetResultMessage(connectionAddress string, connectionSequenceNo int, messageSequenceNo int) (*Message, error) {
	_, ok := (*s).connectionResults[connectionAddress]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrUnknownConnectionAddress, connectionAddress)
	}
	if len((*s).connectionResults[connectionAddress]) < connectionSequenceNo {
		return nil, fmt.Errorf("%w: %d", ErrUnknownConnectionSequence, connectionSequenceNo)
	}

	if len((*s).connectionResults[connectionAddress][connectionSequenceNo-1].Messages) < messageSequenceNo {
		t := time.Now()
		for {
			time.Sleep(10 * time.Millisecond)
			if len((*s).connectionResults[connectionAddress][connectionSequenceNo-1].Messages) >= messageSequenceNo {
				break
			}
			if time.Since(t) > timeout {
				return nil, fmt.Errorf("%w: %d", ErrUnkownMessageSequence, messageSequenceNo)
			}
		}
	}
	return &(*s).connectionResults[connectionAddress][connectionSequenceNo-1].Messages[messageSequenceNo-1], nil
}

// GetResultRawText returns the raw text received by the SMTP server for a given connection address and connection sequence number.
// The connection sequence number is usually 1, but can be increased if a subsequent connection would use the same TCP port (which
// is very unlikely).
func (s *SmtpServer) GetResultRawText(connectionAddress string, connectionSequenceNo int) ([]RawLine, error) {
	crs, ok := (*s).connectionResults[connectionAddress]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrUnknownConnectionAddress, connectionAddress)
	}
	if len(crs) < connectionSequenceNo {
		return nil, fmt.Errorf("%w: %d", ErrUnknownConnectionSequence, connectionSequenceNo)
	}
	cr := crs[connectionSequenceNo-1]

	t := time.Now()
	for {
		if len((*s).rawTextCh) == 0 {
			break
		}
		if time.Since(t) > timeout {
			return nil, ErrTimeout
		}
	}
	return cr.Raw, nil
}
