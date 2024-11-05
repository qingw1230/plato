package sdk

const (
	MsgTypeText = "text"
)

// Message 聊天时使用的消息
type Message struct {
	Type       string // 信令类型
	Name       string
	FromUserID string
	ToUserID   string
	Content    string
	Session    string
}

type Chat struct {
	Nick      string
	UserID    string
	SessionID string
	conn      *connect
}

func NewChat(serverAddr, nick, userID, sessionID string) *Chat {
	return &Chat{
		Nick:      nick,
		UserID:    userID,
		SessionID: sessionID,
		conn:      newConnect(serverAddr),
	}
}

func (c *Chat) Send(msg *Message) {
	c.conn.send(msg)
}

func (c *Chat) Receive() <-chan *Message {
	return c.conn.receive()
}

func (c *Chat) Close() {
	c.conn.close()
}
