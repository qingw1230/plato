package sdk

import (
	"net"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/qingw1230/plato/common/idl/message"
	"github.com/qingw1230/plato/common/tcp"
)

const (
	MsgTypeText      = "text"
	MsgTypeAck       = "ack"
	MsgTypeReConn    = "reConn"
	MsgTypeHeartbeat = "heartbeat"
	MsgLogin         = "loginMsg"
)

// Chat 聊天
type Chat struct {
	Nick      string
	UserID    string
	SessionID string
	conn      *connect
	closeChan chan struct{}
}

// Message 聊天时使用的消息
type Message struct {
	Type       string // 信令类型
	Name       string
	FromUserID string
	ToUserID   string
	Content    string
	Session    string
}

func NewChat(ip net.IP, port int, nick, userID, sessionID string, connID uint64, isReConn bool) *Chat {
	chat := &Chat{
		Nick:      nick,
		UserID:    userID,
		SessionID: sessionID,
		conn:      newConnect(ip, port, connID),
		closeChan: make(chan struct{}, 0),
	}
	go chat.loop()
	if isReConn {
		chat.reConn(connID)
	} else {
		chat.login()
	}
	go chat.heartbeat()
	return chat
}

// Send 发送消息
func (c *Chat) Send(msg *Message) {
	c.conn.receiveChan <- msg
}

func (chat *Chat) GetConnID() uint64 {
	return chat.conn.connID
}

// Receive 接收消息
func (c *Chat) Receive() <-chan *Message {
	return c.conn.receive()
}

func (c *Chat) Close() {
	c.conn.close()
	close(c.closeChan)
	close(c.conn.receiveChan)
	close(c.conn.sendChan)
}

func (c *Chat) loop() {
	for {
		select {
		case <-c.closeChan:
			return
		default:
			data, err := tcp.ReadData(c.conn.conn)
			if err != nil {
				return
			}
			mc := &message.MsgCmd{}
			err = proto.Unmarshal(data, mc)
			if err != nil {
				panic(err)
			}
			var msg *Message
			switch mc.Type {
			case message.CmdType_ACK:
				msg = handAckMsg(c.conn, mc.Payload)
			}
			c.conn.receiveChan <- msg
		}
	}
}

func (c *Chat) login() {
	loginMsg := message.LoginMsg{
		Head: &message.LoginMsgHead{
			DeviceID: 123,
		},
	}
	payload, err := proto.Marshal(&loginMsg)
	if err != nil {
		panic(err)
	}
	c.conn.send(message.CmdType_Login, payload)
}

func (c *Chat) reConn(connID uint64) {
	reConn := message.ReConnMsg{
		Head: &message.ReConnMsgHead{
			ConnID: connID,
		},
	}
	payload, err := proto.Marshal(&reConn)
	if err != nil {
		panic(err)
	}
	c.conn.send(message.CmdType_ReConn, payload)
}

func (c *Chat) heartbeat() {
	t := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-c.closeChan:
			return
		case <-t.C:
			hb := message.HeartbeatMsg{
				Head: &message.HeartbeatMsgHead{},
			}
			payload, err := proto.Marshal(&hb)
			if err != nil {
				panic(err)
			}
			c.conn.send(message.CmdType_Heartbeat, payload)
		}
	}
}
