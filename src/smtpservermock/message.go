package smtpservermock

type messageStatus int

const (
	_ messageStatus = iota
	emptyMessage
	mailFromMessage
	receiptToMessage
	dataMessage
)

type message struct {
	from string
	to   []string
	data string
}

func newMessage() *message {
	return &message{}
}
