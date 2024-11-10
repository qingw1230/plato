package domain

import (
	"sort"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/qingw1230/plato/ipconf/source"
)

type Dispatcher struct {
	sync.RWMutex
	// candidateTable 候选机器列表 key 为 IP:Port
	candidateTable map[string]*Endpoint
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
				panic("ipconf/domain/dispacter.go/Init()")
			}
		}
	}()
}

// Dispactch 派发机器，按降序返回机器列表
func Dispactch(ctx *gin.Context) []*Endpoint {
	eds := d.getCandidateEndpoint(ctx)
	// 计算各机器得分
	for _, ed := range eds {
		ed.CalcScore(ctx)
	}
	sort.Slice(eds, func(i, j int) bool {
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

// getCandidateEndpoint 获取所有候选机器列表
func (dp *Dispatcher) getCandidateEndpoint(_ *gin.Context) []*Endpoint {
	dp.Lock()
	defer dp.Unlock()
	candidateList := make([]*Endpoint, 0, len(dp.candidateTable))
	for _, ed := range dp.candidateTable {
		candidateList = append(candidateList, ed)
	}
	return candidateList
}

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

func (dp *Dispatcher) delNode(event *source.Event) {
	dp.Lock()
	defer dp.Unlock()
	delete(dp.candidateTable, event.Key())
}
