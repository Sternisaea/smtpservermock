package smtpservermock

type command interface {
	getPrefix() string
	execute(transmission *transmission, arg string) error // error is only for critical errors, not for SMTP errors
}
