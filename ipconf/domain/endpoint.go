package domain

import (
	"sync/atomic"
	"unsafe"

	"github.com/gin-gonic/gin"
)

// Endpoint 机器地址信息及资源状态
type Endpoint struct {
	IP          string  `json:"ip"`
	Port        string  `json:"port"`
	ActiveScore float64 `json:"active"`
	StaticScore float64 `json:"static"`
	// Stats 状态平均值
	Stats *Stat `json:"-"`
	// window 资源状态窗口
	window *statWindow `json:"-"`
}

// NewEndpoint 创建一个机器信息，并启动协程不断更新状态
func NewEndpoint(ip, port string) *Endpoint {
	ed := &Endpoint{
		IP:   ip,
		Port: port,
	}
	ed.window = newStatWindow()
	ed.Stats = ed.window.getStat()

	go func() {
		for stat := range ed.window.statChan {
			ed.window.appendStat(stat)
			newStat := ed.window.getStat()
			atomic.SwapPointer((*unsafe.Pointer)((unsafe.Pointer)(&ed.Stats)), unsafe.Pointer(newStat))
		}
	}()
	return ed
}

// UpdateStat 添加机器最新状态
func (e *Endpoint) UpdateStat(s *Stat) {
	e.window.statChan <- s
}

func (e *Endpoint) CalcScore(ctx *gin.Context) {
	e.Stats = e.window.getStat()
	if e.Stats != nil {
		e.ActiveScore = e.Stats.CalcActiveScore()
		e.StaticScore = e.Stats.CalcStaticScore()
	}
}
