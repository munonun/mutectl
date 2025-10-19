package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check health of registered MuteNet nodes",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Load nodes.json and ping
		fmt.Println("📡 Checking all nodes...")
	},
}
