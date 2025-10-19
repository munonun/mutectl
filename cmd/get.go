package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get [CID]",
	Short: "Get file from MuteNet",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cid := args[0]
		// TODO: Pull from peers
		fmt.Printf("🔍 Trying to get file with CID: %s\n", cid)
	},
}
