package tcp

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"time"
)

// ReadData 从 TCP 连接中读取数据包
func ReadData(conn *net.TCPConn) ([]byte, error) {
	var dataLen uint32
	dataLenBuf := make([]byte, 4)
	// 先从连接中读取 4 字节数据，即读取数据包的长度
	if err := readFixedData(conn, dataLenBuf); err != nil {
		return nil, err
	}
	buffer := bytes.NewBuffer(dataLenBuf)
	if err := binary.Read(buffer, binary.BigEndian, &dataLen); err != nil {
		return nil, fmt.Errorf("read headlen error:%s", err.Error())
	}
	if dataLen <= 0 {
		return nil, fmt.Errorf("wrong headlen: %d", dataLen)
	}
	dataBuf := make([]byte, dataLen)
	if err := readFixedData(conn, dataBuf); err != nil {
		return nil, fmt.Errorf("read tcp data error:%s", err.Error())
	}
	return dataBuf, nil
}

// readFixedData 从指定 TCP 连接读取 len(buf) 长度的数据
func readFixedData(conn *net.TCPConn, buf []byte) error {
	conn.SetDeadline(time.Now().Add(time.Duration(120) * time.Second))
	pos := 0
	totalSize := len(buf)
	for {
		c, err := conn.Read(buf[pos:])
		if err != nil {
			return err
		}
		pos += c
		if pos == totalSize {
			break
		}
	}
	return nil
}
