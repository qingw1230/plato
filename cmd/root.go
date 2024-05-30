package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "plato",
	Short: "这是一个 IM 即时通信系统",
	Run:   IM,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func IM(cmd *cobra.Command, args []string) {
	fmt.Println("call IM")
}
