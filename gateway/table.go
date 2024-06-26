package gateway

import "sync"

var tables table

type table struct {
	//  device id 与连接的映射
	did2conn sync.Map
}

func InitTables() {
	tables = table{
		did2conn: sync.Map{},
	}
}
