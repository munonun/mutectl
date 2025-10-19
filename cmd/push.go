package cmd

import (
	"fmt"
	"mutectl/utils"

	"github.com/spf13/cobra"
)

var pushCmd = &cobra.Command{
	Use:   "push [filepath]",
	Short: "Push a file into MuteNet",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		file := args[0]
		cid, err := utils.HashFile(file)
		if err != nil {
			fmt.Println("❌ Error hashing file:", err)
			return
		}
		fmt.Printf("✅ File pushed with CID: %s\n", cid)
		// TODO: Send to P2P cache later
	},
}
