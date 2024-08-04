package source

import (
	"fmt"

	"github.com/qingw1230/plato/common/discovery"
)

// eventChan 存储机器资源变更事件
var eventChan chan *Event

// EventChan 获取机器资源变更 Event 的 chan
func EventChan() <-chan *Event {
	return eventChan
}

type EventType string

const (
	AddNodeEvent EventType = "addNode"
	DelNodeEvent EventType = "delNode"
)

// Event 机器资源信息事件
type Event struct {
	Type         EventType
	IP           string
	Port         string
	ConnectNum   float64
	MessageBytes float64
}

func NewEvent(ed *discovery.EndpointInfo) *Event {
	if ed == nil || ed.MetaData == nil {
		return nil
	}
	var connNum, msgBytes float64
	if data, ok := ed.MetaData["connect_num"]; ok {
		connNum = data.(float64)
	}
	if data, ok := ed.MetaData["message_bytes"]; ok {
		msgBytes = data.(float64)
	}
	return &Event{
		Type:         AddNodeEvent,
		IP:           ed.IP,
		Port:         ed.Port,
		ConnectNum:   connNum,
		MessageBytes: msgBytes,
	}
}

// Key 获取当前 Event 的地址信息 IP:Port
func (e *Event) Key() string {
	return fmt.Sprintf("%s:%s", e.IP, e.Port)
}
