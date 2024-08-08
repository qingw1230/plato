package perf

import (
	"net"
	"syscall"

	"github.com/qingw1230/plato/common/sdk"
)

var (
	TCPConnNum int32
)

func RunMain() {
	var rLimit syscall.Rlimit
	if err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
		panic(err)
	}
	rLimit.Cur = rLimit.Max
	if err := syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
		panic(err)
	}

	for i := 0; i < int(TCPConnNum); i++ {
		sdk.NewChat(net.ParseIP("127.0.0.1"), 8900, "im", "12345", "54321")
	}
}
