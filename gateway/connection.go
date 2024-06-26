package gateway

import "net"

// connection TCP 连接和对应 fd
type connection struct {
	fd   int
	conn *net.TCPConn
}

// RemoteAddr 获取对端地址
func (c *connection) RemoteAddr() string {
	return c.conn.RemoteAddr().String()
}

func (c *connection) Close() {
	err := c.conn.Close()
	panic(err)
}
