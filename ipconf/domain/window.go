package domain

const (
	windowSize = 5
)

// stateWindow 机器资源状态窗口
type stateWindow struct {
	stateQueue []*Stat
	// statChan 用来接收机器最新状态
	statChan chan *Stat
	sumStat  *Stat
	idx      int64
}

func newStatWindow() *stateWindow {
	return &stateWindow{
		stateQueue: make([]*Stat, windowSize),
		statChan:   make(chan *Stat),
		sumStat:    &Stat{},
	}
}

// getStat 获取窗口内资源平均值
func (sw *stateWindow) getStat() *Stat {
	res := sw.sumStat.Clone()
	res.Avg(windowSize)
	return res
}

// appendStat 向窗口中追加新状态
func (sw *stateWindow) appendStat(s *Stat) {
	sw.sumStat.Sub(sw.stateQueue[sw.idx%windowSize])
	sw.stateQueue[sw.idx%windowSize] = s
	sw.sumStat.Add(s)
	sw.idx++
}
