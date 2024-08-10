package state

import (
	"context"
	"sync"
	"time"

	"github.com/qingw1230/plato/common/timingwheel"
	"github.com/qingw1230/plato/state/rpc/client"
	"github.com/qingw1230/plato/state/rpc/service"
)

var cmdChannel chan *service.CmdContext
var connToStateTable sync.Map

type connState struct {
	sync.RWMutex
	heartTimer  *timingwheel.Timer
	reConnTimer *timingwheel.Timer
	connID      uint64
}

// reSetHeartTimer 重置定时器时间，触发时清除连接状态
func (c *connState) reSetHeartTimer() {
	c.Lock()
	defer c.Unlock()
	c.heartTimer.Stop()
	c.heartTimer = AfterFunc(5*time.Second, func() {
		clearState(c.connID)
	})
}

// clearState 为了实现重连，不要立即释放连接的状态，有 10s 的延迟时间
func clearState(connID uint64) {
	if data, ok := connToStateTable.Load(connID); ok {
		state, _ := data.(*connState)
		state.Lock()
		defer state.Unlock()
		state.reConnTimer = AfterFunc(10*time.Second, func() {
			ctx := context.TODO()
			client.DelConn(&ctx, connID, nil)
			connToStateTable.Delete(connID)
		})
	}
}
