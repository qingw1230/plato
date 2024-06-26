package tcp

import (
	"bytes"
	"encoding/binary"
)

// DataPkg TCP 传输的数据包
type DataPkg struct {
	Len  uint32
	Data []byte
}

// Marshal 将 DataPkg 成员依次写入
func (d *DataPkg) Marshal() []byte {
	bytesBuffer := bytes.NewBuffer([]byte{})
	// 先将长度信息以大端方式写入
	binary.Write(bytesBuffer, binary.BigEndian, d.Len)
	return append(bytesBuffer.Bytes(), d.Data...)
}
