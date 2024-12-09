package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/Sternisaea/smtpservermock/src/smtpservermock"
)

func main() {
	// validCert, validKey, err := CreateCertificate("Domain Local", "smtp.intern.local")
	// if err != nil {
	// 	panic(err)
	// }

	validCert := "/tmp/cert.pem986103266"
	validKey := "/tmp/key.pem1133920650"
	fmt.Printf("%s\n%s\n", validCert, validKey)

	// smtpserv, err := smtpservermock.NewSmtpServer(smtpservermock.SslTlsSec, "Mock SMTP Server with SSL/TLS", "127.0.0.1:2526", validCert, validKey)
	// smtpserv, err := smtpservermock.NewSmtpServer(smtpservermock.StartTlsSec, "Mock SMTP Server with STARTTLS", "127.0.0.1:2526", validCert, validKey)
	smtpserv, err := smtpservermock.NewSmtpServer(smtpservermock.NoSecurity, "Mock SMTP Server", "127.0.0.1:2526", "", "")
	if err != nil {
		panic(err)
	}
	if err := smtpserv.ListenAndServe(); err != nil {
		panic(err)
	}
	defer smtpserv.Shutdown()

	time.Sleep(20 * time.Second)
	// time.Sleep(900 * time.Second)

	printResults(smtpserv)
}

func printResults(smtpserv *smtpservermock.SmtpServer) {
	addrs, err := smtpserv.GetConnectionAddresses()
	if err != nil {
		panic(err)
	}

	for _, a := range addrs {
		fmt.Printf("ADDRESS: %s\n", a)
		j := 0
		for {
			j++
			rls, err := smtpserv.GetResultRawText(a, j)
			if err != nil {
				if errors.Is(err, smtpservermock.ErrUnknownConnectionSequence) {
					break
				}
			}
			for _, rl := range rls {
				switch rl.Direction {
				case smtpservermock.RequestDir:
					fmt.Printf("CLIENT: ")
				case smtpservermock.ResponseDir:
					fmt.Printf("SERVER: ")
				}
				fmt.Printf("%s", rl.Text)
			}
			fmt.Println()

			i := 0
			for {
				i++
				msg, err := smtpserv.GetResultMessage(a, j, i)
				if err != nil {
					if errors.Is(err, smtpservermock.ErrUnkownMessageSequence) {
						break
					}
					panic(err)
				}
				fmt.Printf(">> FROM: %s\n", msg.From)
				fmt.Printf(">> TO: %v\n", msg.To)
				fmt.Println(msg.Data)
			}
		}
	}

}
