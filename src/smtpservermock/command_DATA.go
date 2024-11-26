package smtpservermock

import "strings"

type CmdDATA struct{}

func (c *CmdDATA) GetPrefix() string {
	return "DATA"
}

func (c *CmdDATA) Execute(t *Transmission, arg string) error {
	if err := (*t).WriteResponse("354 End message with ."); err != nil {
		return err
	}
	for {
		line, err := (*t).reader.ReadString('\n')
		if err != nil {
			return err
		}
		if strings.TrimSuffix(line, endOfLine) == "." {
			break
		}
		(*t).currentMessage.Data += line
	}
	(*t).messages = append((*t).messages, (*t).currentMessage)
	(*t).initCurrentMessage()
	return (*t).WriteResponse("250 Mail accepted")
}
