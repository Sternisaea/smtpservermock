package smtpservermock

import (
	"bufio"
	"crypto/tls"
	"net"
	"regexp"
	"strings"
)

var emailAngleBracketsRegex = regexp.MustCompile(`^<([a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,})>$`)

type ConnectionType int

const (
	NoType ConnectionType = iota
	HeloType
	EhloType
	QuitType
)

type Transmission struct {
	netConnection    net.Conn
	serverName       string
	starttlsRequired bool
	starttlsConfig   *tls.Config

	reader     *bufio.Reader
	writer     *bufio.Writer
	clientName string
	msgStatus  MessageStatus

	connType       ConnectionType
	starttlsActive bool

	// RECOGNIZED COMMANDS
	commands       []Command
	messages       []*Message
	currentMessage *Message
}

func NewTransmission(connection net.Conn, serverName string) *Transmission {
	return &Transmission{
		netConnection: connection,
		serverName:    serverName,
		reader:        bufio.NewReader(connection),
		writer:        bufio.NewWriter(connection),
	}
}

func (t *Transmission) SetStartTLSConfig(config *tls.Config) {
	(*t).starttlsRequired = true
	(*t).starttlsActive = false
	(*t).starttlsConfig = config
}

func (t *Transmission) SetCommands(cmds []Command) {
	(*t).commands = cmds
}

func (t *Transmission) Process() error {
	(*t).connType = NoType
	(*t).WriteResponse("220 " + (*t).serverName)
	for {
		line, err := (*t).reader.ReadString('\n')
		if err != nil {
			return err
		}
		line = strings.TrimSuffix(line, "\r\n")

		found := false
		for _, c := range (*t).commands {
			if arg, ok := checkPrefix(c, line); ok {
				if err := c.Execute(t, arg); err != nil {
					return err
				}
				found = true
				break
			}
		}
		if !found {
			(*t).WriteResponse("500 Command not recognized")
		}
		if (*t).connType == QuitType {
			return nil
		}
	}
}

func checkPrefix(c Command, line string) (string, bool) {
	prefix := c.GetPrefix()
	if len(line) < len(prefix) || line[:len(prefix)] != prefix {
		return "", false
	}
	return strings.TrimLeft(line[len(prefix):], " "), true
}

func (t *Transmission) WriteResponse(resp string) error {
	if !strings.HasSuffix(resp, "\r\n") {
		resp += "\r\n"
	}
	(*t).writer.WriteString(resp)
	return (*t).writer.Flush()
}

func (t *Transmission) initCurrentMessage() {
	(*t).currentMessage = NewMessage()
	(*t).msgStatus = EmptyMessage
}
