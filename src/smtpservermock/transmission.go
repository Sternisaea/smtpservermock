package smtpservermock

import (
	"bufio"
	"crypto/tls"
	"fmt"
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

type connectionType int

const (
	noType connectionType = iota
	heloType
	ehloType
	quitType
)

type transmission struct {
	security       smtpconst.Security
	netConnection  net.Conn
	serverName     string
	starttlsConfig *tls.Config

	reader    *bufio.Reader
	writer    *bufio.Writer
	id        string
	messageCh chan<- connMessage
	rawTextCh chan<- connRaw

	clientName     string
	connType       connectionType
	starttlsActive bool
	msgStatus      messageStatus

	commands []command

	messages       []*message
	currentMessage *message
}

func newTransmission(security smtpconst.Security, connection net.Conn, serverName string, id string, msgCh chan<- connMessage, rawCh chan<- connRaw) *transmission {
	return &transmission{
		security:      security,
		netConnection: connection,
		serverName:    serverName,
		reader:        bufio.NewReader(connection),
		writer:        bufio.NewWriter(connection),
		id:            id,
		messageCh:     msgCh,
		rawTextCh:     rawCh,
	}
}

func (t *transmission) SetStartTLSConfig(config *tls.Config) {
	(*t).starttlsActive = false
	(*t).starttlsConfig = config
}

func (t *transmission) Process() error {
	(*t).connType = noType
	(*t).setCommands()
	if err := (*t).writeResponse("220 " + (*t).serverName); err != nil {
		return err
	}

	for {
		line, err := (*t).reader.ReadString('\n')
		if err != nil {
			return err
		}
		(*t).writeRaw(line)
		line = strings.TrimSuffix(line, endOfLine)

		found := false
		for _, c := range (*t).commands {
			if arg, ok := checkPrefix(c, line); ok {
				if err := c.execute(t, arg); err != nil {
					return err
				}
				found = true
				break
			}
		}
		if !found {
			cmd := ""
			if words := strings.Fields(line); len(words) > 0 {
				cmd = words[0]
			}
			if err := (*t).writeResponse("500 Command " + cmd + " not recognized"); err != nil {
				return err
			}
			continue
		}
		if (*t).connType == quitType {
			return nil
		}
	}
}

func checkPrefix(c command, line string) (string, bool) {
	prefix := c.getPrefix()
	if len(line) < len(prefix) || strings.ToUpper(line[:len(prefix)]) != prefix {
		return "", false
	}
	return strings.TrimLeft(line[len(prefix):], " "), true
}

func (t *transmission) writeResponse(resp string) error {
	if !strings.HasSuffix(resp, endOfLine) {
		resp += endOfLine
	}
	(*t).writeRaw(resp)
	(*t).writer.WriteString(resp)
	return (*t).writer.Flush()
}

func (t *transmission) writeRaw(rawtext string) {
	(*t).rawTextCh <- connRaw{id: (*t).id, rawtext: rawtext}
}

func (t *transmission) initCurrentMessage() {
	(*t).currentMessage = newMessage()
	(*t).msgStatus = emptyMessage
}

func (t *transmission) setCommands() {
	cmds := []command{&cmdEHLO{}, &cmdHELO{}, &cmdQuit{}, &cmdNOOP{}, &cmdHELP{}, &cmdRSET{}, &cmdVRFY{}}
	if (*t).security == smtpconst.StartTlsSec && !(*t).starttlsActive {
		cmds = append(cmds, []command{&cmdSTARTTLS{}}...)
	}
	switch (*t).connType {
	case heloType:
		cmds = append(cmds, []command{&cmdMAILFROM{}, &cmdRCPTTO{}, &cmdDATA{}}...)
	case ehloType:
		cmds = append(cmds, []command{&cmdMAILFROM{}, &cmdRCPTTO{}, &cmdDATA{}, &cmdAUTH{}}...)
	}
	(*t).commands = cmds
}
