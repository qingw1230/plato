package gateway

import (
	"net"
	"sync/atomic"
)

var nextConnID uint64

// connection TCP 连接和对应 fd
type connection struct {
	id   uint64 // 进程级的生命周期
	fd   int
	e    *epoller
	conn *net.TCPConn
}

func NewConnection(conn *net.TCPConn) *connection {
	connID := atomic.AddUint64(&nextConnID, 1)
	return &connection{
		id:   connID,
		fd:   socketFD(conn),
		conn: conn,
	}
}

func (c *connection) Close() {
	ep.tables.Delete(c.id)
	if c.e != nil {
		c.e.fdToConnTable.Delete(c.fd)
	}
	err := c.conn.Close()
	panic(err)
}

// RemoteAddr 获取对端地址
func (c *connection) RemoteAddr() string {
	return c.conn.RemoteAddr().String()
}

func (c *connection) BindEpoller(e *epoller) {
	c.e = e
}
