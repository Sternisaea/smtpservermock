package smtpservermock

import (
	"bufio"
	"fmt"
	"log"
	"strings"
)

type Commands []Command

func (cs *Commands) ProcesLines(reader *bufio.Reader, writer *bufio.Writer) error {
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return err
		}

		line = strings.TrimSpace(line)
		log.Printf("Received: %s\n", line)

		if err := cs.processCommands(reader, writer, line); err != nil {
			return err
		}
	}
}

func (cs *Commands) processCommands(reader *bufio.Reader, writer *bufio.Writer, line string) error {
	for _, cmd := range *cs {
		if strings.HasPrefix(line, cmd.GetPrefix()) {
			if quit, err := cmd.Execute(reader, writer, line); quit || err != nil {
				if quit {
					return nil
				}
				return fmt.Errorf("error processing command %s: %w", cmd.GetPrefix(), err)
			}
			return nil
		}
	}
	return fmt.Errorf("unkown command %s", line)
}
