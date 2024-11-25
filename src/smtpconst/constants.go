package smtpconst

type Security string

const (
	NoSecurity  Security = ""
	StartTlsSec Security = "starttls"
	SslTlsSec   Security = "ssl/tls"
)

type AuthenticationMethod string

const (
	NoAuthentication AuthenticationMethod = ""
	PlainAuth        AuthenticationMethod = "plain"
	CramMd5Auth      AuthenticationMethod = "cram-md5"
)

type Status int

const (
	VoidStatus Status = iota
	HeloStatus
	EhloStatus
	StartTlsStatus
	MailFromStatus
	ReceiptToStatus
	DataStatus
	QuitStatus
)
