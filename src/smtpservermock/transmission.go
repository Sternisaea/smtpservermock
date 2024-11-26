package smtpservermock

import (
	"bufio"
	"crypto/tls"
	"net"
	"regexp"
	"strings"

	"github.com/Sternisaea/smtpservermock/src/smtpconst"
)

var (
	endOfLine              = "\r\n"
	textAngleBracketsRegex = regexp.MustCompile(`<(.*?)>`)
)

//var emailRegex = regexp.MustCompile(`^([a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,})$`)

type ConnectionType int

const (
	NoType ConnectionType = iota
	HeloType
	EhloType
	QuitType
)

type Transmission struct {
	security         smtpconst.Security
	netConnection    net.Conn
	serverName       string
	starttlsRequired bool
	starttlsConfig   *tls.Config

	reader *bufio.Reader
	writer *bufio.Writer

	clientName     string
	connType       ConnectionType
	starttlsActive bool
	msgStatus      MessageStatus

	// RECOGNIZED COMMANDS
	commands       []Command
	messages       []*Message
	currentMessage *Message
}

func NewTransmission(security smtpconst.Security, connection net.Conn, serverName string) *Transmission {
	return &Transmission{
		security:      smtpconst.NoSecurity,
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

func (t *Transmission) Process() error {
	(*t).connType = NoType
	(*t).setCommands()
	if err := (*t).WriteResponse("220 " + (*t).serverName); err != nil {
		return err
	}
	for {
		line, err := (*t).reader.ReadString('\n')
		if err != nil {
			return err
		}
		line = strings.TrimSuffix(line, endOfLine)

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
			if err := (*t).WriteResponse("500 Command not recognized"); err != nil {
				return err
			}
			continue
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
	if !strings.HasSuffix(resp, endOfLine) {
		resp += endOfLine
	}
	(*t).writer.WriteString(resp)
	return (*t).writer.Flush()
}

func (t *Transmission) initCurrentMessage() {
	(*t).currentMessage = NewMessage()
	(*t).msgStatus = EmptyMessage
}

func (t *Transmission) setCommands() {
	cmds := []Command{&CmdEHLO{}, &CmdHELO{}, &CmdQuit{}, &CmdNOOP{}, &CmdHELP{}, &CmdRSET{}, &CmdVRFY{}}
	if (*t).security == smtpconst.StartTlsSec && !(*t).starttlsActive {
		cmds = append(cmds, []Command{&CmdSTARTTLS{}}...)
	}
	switch (*t).connType {
	case HeloType:
		cmds = append(cmds, []Command{&CmdMAILFROM{}, &CmdRCPTTO{}, &CmdDATA{}}...)
	case EhloType:
		cmds = append(cmds, []Command{&CmdMAILFROM{}, &CmdRCPTTO{}, &CmdDATA{}}...)
	}
	(*t).commands = cmds
}
