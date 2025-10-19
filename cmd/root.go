package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "mutectl",
	Short: "MuteNet CLI tool for pushing, fetching, and checking nodes",
	Long:  `Control your MuteNet CDN with commands like push, get, check, and more.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}

func init() {
	rootCmd.AddCommand(pushCmd)
	rootCmd.AddCommand(getCmd)
	rootCmd.AddCommand(checkCmd)
}
