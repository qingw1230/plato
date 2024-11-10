package cmd

import (
	"github.com/qingw1230/plato/ipconf"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(ipConfCmd)
}

var ipConfCmd = &cobra.Command{
	Use: "ipconf",
	Run: IPConfHandle,
}

func IPConfHandle(cmd *cobra.Command, args []string) {
	ipconf.RunMain(ConfigPath)
}
