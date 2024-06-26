package sdk

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/qingw1230/plato/common/tcp"
)

// connect 包含 IP 端口号，以及用于发送、接收消息的 chan *Message
type connect struct {
	conn        *net.TCPConn
	sendChan    chan *Message
	receiveChan chan *Message
}

func newConnect(ip net.IP, port int) *connect {
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

	go func() {
		for {
			// 不断从 TCP 连接中读取数据，解码后写入 clientConn.receiveChan
			data, err := tcp.ReadData(conn)
			if err != nil {
				fmt.Printf("ReadData err:%+v\n", err)
			}
			msg := &Message{}
			json.Unmarshal(data, msg)
			clientConn.receiveChan <- msg
		}
	}()

	return clientConn
}

// send 向 connect 中发送消息
func (c *connect) send(data *Message) {
	bytes, _ := json.Marshal(data)
	dataPkg := tcp.DataPkg{
		Len:  uint32(len(bytes)),
		Data: bytes,
	}
	xx := dataPkg.Marshal()
	c.conn.Write(xx)
}

// receive 从 connect 中获取消息
func (c *connect) receive() <-chan *Message {
	return c.receiveChan
}

func (c *connect) close() {
	c.conn.Close()
}
