package perf

import (
	"net"

	"github.com/qingw1230/plato/common/sdk"
)

var (
	TCPConnNum int32
)

func RunMain() {
	// table := make([]*sdk.Chat, 0, int(TCPConnNum))
	for i := 0; i < int(TCPConnNum); i++ {
		sdk.NewChat(net.ParseIP("127.0.0.1"), 8900, "im", "12345", "54321")
		// table = append(table, c)
	}

	// for i := 0; i < int(TCPConnNum); i++ {
	// 	table[i].Close()
	// }
}
