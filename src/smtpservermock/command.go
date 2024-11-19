package smtpservermock

import (
	"bufio"
	"encoding/base64"
	"strings"
)

type Command interface {
	GetPrefix() string
	Execute(reader *bufio.Reader, writer *bufio.Writer, line string) (bool, error)
}

type CmdEHLO struct{}

func (c *CmdEHLO) GetPrefix() string {
	return "EHLO"
}

func (c *CmdEHLO) Execute(reader *bufio.Reader, writer *bufio.Writer, line string) (bool, error) {
	writer.WriteString("250-Mock SMTP Server\r\n250-AUTH PLAIN\r\n250 OK\r\n")
	return false, writer.Flush()
}

type CmdAuth struct{}

func (c *CmdAuth) GetPrefix() string {
	return "AUTH"
}

func (c *CmdAuth) Execute(reader *bufio.Reader, writer *bufio.Writer, line string) (bool, error) {
	writer.WriteString("334 \r\n")
	writer.Flush()

	authLine, err := reader.ReadString('\n')
	if err != nil {
		return false, err
	}
	authLine = strings.TrimSpace(authLine)
	decoded, err := base64.StdEncoding.DecodeString(authLine)
	if err != nil {
		writer.WriteString("535 Authentication failed\r\n")
		writer.Flush()
		return false, err
	}

	parts := strings.SplitN(string(decoded), "\x00", 3)
	if len(parts) != 3 || parts[1] != "username" || parts[2] != "password" {
		writer.WriteString("535 Authentication failed\r\n")
		writer.Flush()
		return false, nil
	}

	writer.WriteString("235 Authentication successful\r\n")
	return false, writer.Flush()
}

type CmdQuit struct{}

func (c *CmdQuit) GetPrefix() string {
	return "QUIT"
}

func (c *CmdQuit) Execute(reader *bufio.Reader, writer *bufio.Writer, line string) (bool, error) {
	writer.WriteString("221 Bye\r\n")
	writer.Flush()
	return true, nil
}
