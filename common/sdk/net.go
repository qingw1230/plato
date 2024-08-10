package sdk

import (
	"fmt"
	"net"

	"github.com/golang/protobuf/proto"
	"github.com/qingw1230/plato/common/idl/message"
	"github.com/qingw1230/plato/common/tcp"
)

// connect 连接结构啼
type connect struct {
	conn        *net.TCPConn // TCP 连接
	sendChan    chan *Message
	receiveChan chan *Message
	connID      uint64
}

func newConnect(ip net.IP, port int, connID uint64) *connect {
	clientConn := &connect{
		sendChan:    make(chan *Message),
		receiveChan: make(chan *Message),
	}
	addr := &net.TCPAddr{
		IP:   ip,
		Port: port,
	}
	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		fmt.Printf("DialTCP.err=%+v\n", err)
		return nil
	}

	clientConn.conn = conn
	if connID != 0 {
		clientConn.connID = connID
	}
	return clientConn
}

// handAckMsg 处理 ACK 消息
func handAckMsg(c *connect, data []byte) *Message {
	ackMsg := &message.ACKMsg{}
	proto.Unmarshal(data, ackMsg)
	switch ackMsg.Type {
	case message.CmdType_Login:
		c.connID = ackMsg.ConnID
	}
	return &Message{
		Type:       MsgTypeAck,
		Name:       "im",
		FromUserID: "1234",
		ToUserID:   "123456",
		Content:    ackMsg.Msg,
	}
}

// send 向 connect 中发送消息
func (c *connect) send(ct message.CmdType, payload []byte) {
	msgCmd := message.MsgCmd{
		Type:    ct,
		Payload: payload,
	}
	msg, err := proto.Marshal(&msgCmd)
	if err != nil {
		panic(err)
	}
	dataPkg := tcp.DataPkg{
		Len:  uint32(len(msg)),
		Data: msg,
	}
	c.conn.Write(dataPkg.Marshal())
}

// receive 从 connect 中获取消息
func (c *connect) receive() <-chan *Message {
	return c.receiveChan
}

func (c *connect) close() {
	c.conn.Close()
}
