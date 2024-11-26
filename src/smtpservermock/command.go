package smtpservermock

type Command interface {
	GetPrefix() string
	Execute(transmission *Transmission, arg string) error // error is only for critical errors, not for SMTP errors
}
