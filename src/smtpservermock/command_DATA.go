package smtpservermock

import (
	"strings"
)

type cmdDATA struct{}

func (c *cmdDATA) getPrefix() string {
	return "DATA"
}

func (c *cmdDATA) execute(t *transmission, arg string) error {
	if err := (*t).writeResponse("354 End message with ."); err != nil {
		return err
	}
	for {
		line, err := (*t).reader.ReadString('\n')
		if err != nil {
			return err
		}
		(*t).rawBuffer = append((*t).rawBuffer, RawLine{Direction: RequestDir, Text: line})
		if strings.TrimSuffix(line, endOfLine) == "." {
			break
		}
		(*t).currentMessage.Data += line
	}
	err := (*t).writeResponse("250 Mail accepted")
	(*t).sendMessage()
	return err
}
