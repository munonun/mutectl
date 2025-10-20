package cmd

import (
	"encoding/json"
	"fmt"
	"mutectl/utils"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var joinOriginCmd = &cobra.Command{
	Use:   "join-origin",
	Short: "Register this node as a root (only before release)",
	Run: func(cmd *cobra.Command, args []string) {
		if !utils.IsDevMode() {
			fmt.Println("🚫 This command is only allowed in dev mode.")
			return
		}

		ip := utils.GetMyIP() // 네 IP 알아오는 함수
		port := 8989          // 기본 포트 (필요하면 바꿔도 됨)

		node := utils.Node{
			IP:   ip,
			Port: port,
		}

		nodes := []utils.Node{node}
		data, _ := json.MarshalIndent(nodes, "", "  ")

		configDir := utils.GetConfigDir()
		_ = os.MkdirAll(configDir, 0700)
		path := filepath.Join(configDir, "nodes.json")
		err := os.WriteFile(path, data, 0644)
		if err != nil {
			fmt.Println("❌ Failed to save root node:", err)
			return
		}

		fmt.Printf("✅ Registered current node (%s:%d) as root.\n", ip, port)
	},
}
