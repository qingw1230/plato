package domain

import (
	"sync/atomic"
	"unsafe"
)

// Endport 机器地址信息及资源状态
type Endport struct {
	IP          string  `json:"ip"`
	Port        string  `json:"port"`
	ActiveScore float64 `json:"-"`
	StaticScore float64 `json:"-"`
	// Stats 状态平均值
	Stats *Stat `json:"-"`
	// window 资源状态窗口
	window *stateWindow `json:"-"`
}

// NewEndport 创建一个机器信息，并启动协程更新状态
func NewEndport(ip, port string) *Endport {
	ed := &Endport{
		IP:   ip,
		Port: port,
	}
	ed.window = newStatWindow()
	ed.Stats = ed.window.getStat()

	go func() {
		// 在协程中不断更新机器状态
		for stat := range ed.window.statChan {
			ed.window.appendStat(stat)
			newStat := ed.window.getStat()
			atomic.SwapPointer((*unsafe.Pointer)((unsafe.Pointer)(ed.Stats)), unsafe.Pointer(newStat))
		}
	}()
	return ed
}

// UpdateStat 向机器 chan 添加新状态
func (ed *Endport) UpdateStat(s *Stat) {
	ed.window.statChan <- s
}

// CalculateScore 重新计算机器资源分
func (ed *Endport) CalculateScore(ctx *IPConfConext) {
	if ed.Stats != nil {
		ed.ActiveScore = ed.Stats.CalculateActiveScore()
		ed.StaticScore = ed.Stats.CalculateStaticScore()
	}
}
