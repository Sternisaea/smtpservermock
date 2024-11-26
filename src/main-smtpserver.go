package main

import (
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

	time.Sleep(900 * time.Second)
}
