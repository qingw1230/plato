package sdk

// connect 包含 IP:Port、以及用于发送、接收消息的 chan
type connect struct {
	serverAddr  string
	sendChan    chan *Message
	receiveChan chan *Message
}

func newConnect(serverAddr string) *connect {
	return &connect{
		serverAddr:  serverAddr,
		sendChan:    make(chan *Message),
		receiveChan: make(chan *Message),
	}
}

// send 向 connect 发送消息
func (c *connect) send(data *Message) {
	c.receiveChan <- data
}

// receive 从 connect 获取消息
func (c *connect) receive() <-chan *Message {
	return c.receiveChan
}

func (c *connect) close() {
	close(c.sendChan)
	close(c.receiveChan)
}
