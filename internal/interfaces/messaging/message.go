package messaging

// Message Структура сообщения для обработки.
type Message struct {
	Text            string
	Command         string
	Arguments       string
	IsCommand       bool
	UserID          int64
	UserName        string
	UserDisplayName string
	Date            int64
	IsCallback      bool
	CallbackMsgID   string
}
