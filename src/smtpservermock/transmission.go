package smtpservermock

import (
	"bufio"
	"crypto/tls"
	"net"
	"strings"

	"github.com/Sternisaea/smtpservermock/src/smtpconst"
)

type Transmission struct {
	netConnection    net.Conn
	serverName       string
	starttlsRequired bool
	starttlsConfig   *tls.Config

	reader         *bufio.Reader
	writer         *bufio.Writer
	clientName     string
	status         smtpconst.Status
	starttlsActive bool

	// RECOGNIZED COMMANDS
	commands []Command
	messages []Message
}

type Message struct {
	From string
	To   []string
	Data string
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
	(*t).status = smtpconst.VoidStatus
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
		if (*t).status == smtpconst.QuitStatus {
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
