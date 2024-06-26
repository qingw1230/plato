package cmd

import (
	"github.com/qingw1230/plato/perf"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(perfCmd)
	perfCmd.PersistentFlags().Int32Var(&perf.TCPConnNum, "tcp_conn_num", 15000, "tcp 连接数默认 10000")
}

var perfCmd = &cobra.Command{
	Use: "perf",
	Run: PerfHandle,
}

func PerfHandle(cmd *cobra.Command, args []string) {
	perf.RunMain()
}
