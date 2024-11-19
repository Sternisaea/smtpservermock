package smtpservermock

import (
	"bufio"
	"errors"
	"log"
	"net"
)

type SmtpServer struct {
	security SmtpSecurity
	listener net.Listener
}

func NewSmtpServer(security SmtpSecurity) *SmtpServer {
	return &SmtpServer{security: security}
}

func (s *SmtpServer) ListenAndServe() error {
	var err error
	(*s).listener, err = (*s).security.SetupListener()
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
	return (*s).security.ShutdownListener((*s).listener)
}

func (s *SmtpServer) handle(conn net.Conn) {
	defer conn.Close()

	//	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	writer.WriteString("220 Mock SMTP Server\r\n")
	writer.Flush()

	// Process Commands

}

// #### STARTTLS ####
// if strings.HasPrefix(line, "EHLO") {
// 	writer.WriteString("250-Mock SMTP Server\r\n250-STARTTLS\r\n250 OK\r\n")
// 	writer.Flush()
// } else if strings.HasPrefix(line, "STARTTLS") {
// 	writer.WriteString("220 Ready to start TLS\r\n")
// 	writer.Flush()

// 	tlsConn := tls.Server(conn, s.tlsConfig)
// 	if err := tlsConn.Handshake(); err != nil {
// 		log.Printf("TLS handshake failed: %s", err)
// 		return
// 	}

// 	reader = bufio.NewReader(tlsConn)
// 	writer = bufio.NewWriter(tlsConn)
// } else if strings.HasPrefix(line, "AUTH PLAIN") {
// 	writer.WriteString("235 Authentication successful\r\n")
// 	writer.Flush()
// } else {
// 	writer.WriteString("500 Unrecognized command\r\n")
// 	writer.Flush()
// }

// #### SSL-TLS / None ####
// if strings.HasPrefix(line, "EHLO") {
// 	writer.WriteString("250-Mock SMTP Server\r\n250-AUTH PLAIN\r\n250 OK\r\n")
// 	writer.Flush()
// } else if strings.HasPrefix(line, "AUTH PLAIN") {
// 	writer.WriteString("235 Authentication successful\r\n")
// 	writer.Flush()
// } else {
// 	writer.WriteString("500 Unrecognized command\r\n")
// 	writer.Flush()
// }
