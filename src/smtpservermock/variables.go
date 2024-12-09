package smtpservermock

import "time"

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

type messageStatus int

const (
	_ messageStatus = iota
	emptyMessage
	mailFromMessage
	receiptToMessage
	dataMessage
)

type Direction int

const (
	RequestDir  Direction = 0
	ResponseDir Direction = 1
)

type RawLine struct {
	Direction Direction
	Text      string
}

type Message struct {
	From string
	To   []string
	Data string
}

type Result struct {
	EntryNo int
	Start   time.Time
	End     time.Time

	Raw      []RawLine
	Messages []Message
}

type transmissionRawLines struct {
	address string
	entryNo int
	lines   []RawLine
}

type transmissionMessage struct {
	address string
	entryNo int
	message Message
}
