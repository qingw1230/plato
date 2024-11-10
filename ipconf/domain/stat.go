package domain

import "math"

// Stat 表示 endpoint 对应机器资源剩余量
type Stat struct {
	// ConnectNum 持有的长连接数量的剩余值
	ConnectNum float64
	// MessageBytes 每秒收发消息字节数的剩余值
	MessageBytes float64
}

// CalcActiveScore 获取动态分
func (s *Stat) CalcActiveScore() float64 {
	return getGB(s.MessageBytes)
}

// CalcStaticScore 获取静态分
func (s *Stat) CalcStaticScore() float64 {
	return s.ConnectNum
}

// Avg 将机器资源除以指定值，用于获取窗口内机器资源平均值
func (s *Stat) Avg(num float64) {
	s.ConnectNum /= num
	s.MessageBytes /= num
}

// Add 增加机器资源
func (s *Stat) Add(st *Stat) {
	if st == nil {
		return
	}
	s.ConnectNum += st.ConnectNum
	s.MessageBytes += st.MessageBytes
}

// Sub 减少机器资源
func (s *Stat) Sub(st *Stat) {
	if st == nil {
		return
	}
	s.ConnectNum -= st.ConnectNum
	s.MessageBytes -= st.MessageBytes
}

// Clone 克隆资源信息
func (s *Stat) Clone() *Stat {
	newStat := &Stat{
		ConnectNum:   s.ConnectNum,
		MessageBytes: s.MessageBytes,
	}
	return newStat
}

// getGB 将字节转换为 GB，并保留 2 位小数
func getGB(m float64) float64 {
	return decimal(m / (1 << 30))
}

func decimal(value float64) float64 {
	return math.Trunc(value*1e2+0.5) * 1e-2
}

func min(a, b, c float64) float64 {
	m := func(k, j float64) float64 {
		if k > j {
			return j
		}
		return k
	}
	return m(a, m(b, c))
}
