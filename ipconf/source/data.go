package source

import (
	"context"

	"github.com/bytedance/gopkg/util/logger"
	"github.com/qingw1230/plato/common/config"
	"github.com/qingw1230/plato/common/discovery"
)

func Init() {
	eventChan = make(chan *Event)
	ctx := context.Background()
	go DataHandler(&ctx)

	if config.IsDebug() {
		ctx := context.Background()
		testServiceRegister(&ctx, "7896", "node1")
		testServiceRegister(&ctx, "7897", "node2")
		testServiceRegister(&ctx, "7898", "node3")
	}
}

// DataHandler 监听 etcd 中以 /im/ip_dispatcher 为前缀的键值对，设置发生变化时的处理函数
func DataHandler(ctx *context.Context) {
	dis := discovery.NewServiceDiscovery(ctx)
	defer dis.Close()

	// setFunc 返还指定机器的资源
	setFunc := func(key, value string) {
		if ed, err := discovery.UnMarshal([]byte(value)); err == nil {
			if event := NewEvent(ed); event != nil {
				event.Type = AddNodeEvent
				// 将添加事件写入 chan，由 dispatcher 处理
				eventChan <- event
			} else {
				logger.CtxErrorf(*ctx, "DataHandler.setFunc.err: %s", err.Error())
			}
		}
	}

	// delFunc 消耗指定机器的资源
	delFunc := func(key, value string) {
		if ed, err := discovery.UnMarshal([]byte(value)); err == nil {
			if event := NewEvent(ed); event != nil {
				event.Type = DelNodeEvent
				// 将删除事件写入 chan，由 dispatcher 处理
				eventChan <- event
			} else {
				logger.CtxErrorf(*ctx, "DataHandler.delFunc.err: %s", err.Error())
			}
		}
	}

	err := dis.WatchService(config.GetServicePathForIPConf(), setFunc, delFunc)
	if err != nil {
		panic(err)
	}
}
