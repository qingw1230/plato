package domain

const (
	windowSize = 5
)

type statWindow struct {
	// statChan 用来接收机器最新状态
	statChan  chan *Stat
	sumStat   *Stat
	statQueue []*Stat
	idx       int64
}

func newStatWindow() *statWindow {
	return &statWindow{
		statChan:  make(chan *Stat),
		statQueue: make([]*Stat, windowSize),
		sumStat:   &Stat{},
	}
}

// getStat 获取窗口内资源平均值
func (sw *statWindow) getStat() *Stat {
	ans := sw.sumStat.Clone()
	ans.Avg(windowSize)
	return ans
}

// appendStat 向窗口内追加新状态
func (sw *statWindow) appendStat(s *Stat) {
	sw.sumStat.Sub(sw.statQueue[sw.idx%windowSize])
	sw.statQueue[sw.idx%windowSize] = s
	sw.sumStat.Add(s)
	sw.idx++
}
