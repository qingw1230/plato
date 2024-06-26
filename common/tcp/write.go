package tcp

import "net"

// SendData 将 data 写入指定 TCP 连接
func SendData(conn *net.TCPConn, data []byte) error {
	totalLen := len(data)
	writeLen := 0
	for {
		len, err := conn.Write(data[writeLen:])
		if err != nil {
			return err
		}
		writeLen += len
		if writeLen >= totalLen {
			break
		}
	}
	return nil
}
