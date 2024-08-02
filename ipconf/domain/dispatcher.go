package domain

import (
	"sort"
	"sync"

	"github.com/qingw1230/plato/ipconf/source"
)

// Dispatcher
type Dispatcher struct {
	// candidateTable 候选机器列表 key 为 IP:Port
	candidateTable map[string]*Endpoint
	sync.RWMutex
}

var d *Dispatcher

func Init() {
	d = &Dispatcher{}
	d.candidateTable = make(map[string]*Endpoint)
	go func() {
		for event := range source.EventChan() {
			switch event.Type {
			case source.AddNodeEvent:
				d.addNode(event)
			case source.DelNodeEvent:
				d.delNode(event)
			default:
				panic("ipconf/domain/dispatcher.go/Init()")
			}
		}
	}()
}

// Dispatch 派发机器，按降序返回机器列表
func Dispatch(ctx *IPConfConext) []*Endpoint {
	eds := d.getCandidateEndpoint(ctx)
	// 计算各机器得分
	for _, ed := range eds {
		ed.CalculateScore(ctx)
	}
	// 根据得分为机器排序
	sort.Slice(eds, func(i, j int) bool {
		// 优先基于动态分进行排序
		if eds[i].ActiveScore > eds[j].ActiveScore {
			return true
		}
		if eds[i].ActiveScore == eds[j].ActiveScore {
			return eds[i].StaticScore > eds[j].StaticScore
		}
		return false
	})
	return eds
}

// getCandidateEndpoint 获取候选机器列表
func (dp *Dispatcher) getCandidateEndpoint(_ *IPConfConext) []*Endpoint {
	dp.RLock()
	defer dp.RUnlock()
	candidateList := make([]*Endpoint, 0, len(dp.candidateTable))
	for _, ed := range dp.candidateTable {
		candidateList = append(candidateList, ed)
	}
	return candidateList
}

// delNode 从派发器中删除指定机器
func (dp *Dispatcher) delNode(event *source.Event) {
	dp.Lock()
	defer dp.Unlock()
	delete(dp.candidateTable, event.Key())
}

// addNode 向派发器中添加机器，若已存在则更新状态
func (dp *Dispatcher) addNode(event *source.Event) {
	dp.Lock()
	defer dp.Unlock()
	var (
		ed *Endpoint
		ok bool
	)
	if ed, ok = dp.candidateTable[event.Key()]; !ok {
		ed = NewEndpoint(event.IP, event.Port)
		dp.candidateTable[event.Key()] = ed
	}
	ed.UpdateStat(&Stat{
		ConnectNum:   event.ConnectNum,
		MessageBytes: event.MessageBytes,
	})
}
