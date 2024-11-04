package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var ConfigPath string

func init() {
	rootCmd.PersistentFlags().StringVar(
		&ConfigPath,
		"config",
		"./im.yaml",
		"config file (default is ./im.yaml)",
	)
}

var rootCmd = &cobra.Command{
	Use:   "plato",
	Short: "这是一个分布式即时通信系统。",
	Run:   IM,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func IM(cmd *cobra.Command, args []string) {
	fmt.Println("call IM")
	fmt.Println("args", args)
	fmt.Println(ConfigPath)
}
