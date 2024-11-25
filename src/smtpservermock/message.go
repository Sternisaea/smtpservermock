package smtpservermock

type MessageStatus int

const (
	_ MessageStatus = iota
	EmptyMessage
	MailFromMessage
	ReceiptToMessage
	DataMessage
)

type Message struct {
	From string
	To   []string
	Data string
}

func NewMessage() *Message {
	return &Message{}
}
