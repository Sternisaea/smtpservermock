package smtpservermock

import (
	"bufio"
	"crypto/tls"
	"net"
	"regexp"
	"strings"
)

var (
	endOfLine              = "\r\n"
	textAngleBracketsRegex = regexp.MustCompile(`<(.*?)>`)
)

type connectionType int

const (
	noType connectionType = iota
	heloType
	ehloType
	quitType
)

type transmission struct {
	address        string
	entryNo        int
	security       Security
	netConnection  net.Conn
	serverName     string
	starttlsConfig *tls.Config

	clientName     string
	connType       connectionType
	starttlsActive bool
	commands       []command
	msgStatus      messageStatus
	currentMessage *Message
	rawBuffer      []RawLine

	reader    *bufio.Reader
	writer    *bufio.Writer
	rawTextCh chan<- transmissionRawLines
	messageCh chan<- transmissionMessage
}

func newTransmission(address string, entryNo int, security Security, connection net.Conn, serverName string, rawCh chan<- transmissionRawLines, msgCh chan<- transmissionMessage) *transmission {
	return &transmission{
		address:       address,
		entryNo:       entryNo,
		security:      security,
		netConnection: connection,
		serverName:    serverName,

		reader: bufio.NewReader(connection),
		writer: bufio.NewWriter(connection),

		rawTextCh: rawCh,
		messageCh: msgCh,
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
		(*t).rawBuffer = append((*t).rawBuffer, RawLine{Direction: RequestDir, Text: line})
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
			(*t).flushRawBuffer()
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
	(*t).rawBuffer = append((*t).rawBuffer, RawLine{Direction: ResponseDir, Text: resp})
	(*t).writer.WriteString(resp)
	return (*t).writer.Flush()
}

func (t *transmission) initCurrentMessage() {
	(*t).flushRawBuffer()
	(*t).currentMessage = &Message{}
	(*t).msgStatus = emptyMessage
}

func (t *transmission) setCommands() {
	cmds := []command{&cmdEHLO{}, &cmdHELO{}, &cmdQuit{}, &cmdNOOP{}, &cmdHELP{}, &cmdRSET{}, &cmdVRFY{}}
	if (*t).security == StartTlsSec && !(*t).starttlsActive {
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

func (t *transmission) flushRawBuffer() {
	if len((*t).rawBuffer) == 0 {
		return
	}
	(*t).rawTextCh <- transmissionRawLines{address: (*t).address, entryNo: (*t).entryNo, lines: (*t).rawBuffer}
	(*t).rawBuffer = []RawLine{}
}

func (t *transmission) sendMessage() {
	(*t).messageCh <- transmissionMessage{address: (*t).address, entryNo: (*t).entryNo, message: *(*t).currentMessage}
	(*t).initCurrentMessage()

}
