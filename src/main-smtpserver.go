package main

import (
	"fmt"
	"time"

	"github.com/Sternisaea/smtpservermock/src/smtpconst"
	"github.com/Sternisaea/smtpservermock/src/smtpservermock"
)

func main() {
	smtpserv, err := smtpservermock.NewSmtpServer(smtpconst.NoSecurity, "Mock SMTP Server", "127.0.0.1:2526", "", "")
	if err != nil {
		panic(err)
	}
	if err := smtpserv.ListenAndServe(); err != nil {
		panic(err)
	}
	defer smtpserv.Shutdown()

	// time.Sleep(900 * time.Second)
	time.Sleep(120 * time.Second)

	cms := smtpserv.GetConnectionMessages()
	fmt.Printf("%#v\n\n", cms)

	craws := smtpserv.GetRawText()
	for c, lns := range craws {
		for _, l := range lns {
			fmt.Printf("%s  %s", c, l)
		}
	}

}
