package cmd

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/quic-go/quic-go"
	"github.com/spf13/cobra"
)

type Node struct {
	IP      string `json:"ip"`
	Port    int    `json:"port"`
	Country string `json:"country"`
	Name    string `json:"name"`
}

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check health of registered MuteNet nodes via QUIC",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("📡 Checking all nodes...")

		// 1. Load nodes.json
		home, _ := os.UserHomeDir()
		nodesPath := filepath.Join(home, ".mutenet", "nodes.json")
		data, err := os.ReadFile(nodesPath)
		if err != nil {
			fmt.Println("❌ Failed to read nodes.json:", err)
			return
		}

		var nodes []Node
		if err := json.Unmarshal(data, &nodes); err != nil {
			fmt.Println("❌ Failed to parse nodes.json:", err)
			return
		}

		// 2. Check each node using QUIC
		for _, node := range nodes {
			addr := fmt.Sprintf("%s:%d", node.IP, node.Port)
			tlsConf := &tls.Config{
				InsecureSkipVerify: true,
				NextProtos:         []string{"mutenet"},
			}
			start := time.Now()
			session, err := quic.DialAddr(context.Background(), addr, tlsConf, &quic.Config{})
			if err != nil {
				if err != nil {
					fmt.Printf("❌ [%s] %s (%s) - Unreachable\n", node.Country, node.Name, addr)
					continue
				}
				defer session.CloseWithError(0, "check done")
				elapsed := time.Since(start).Milliseconds()

				fmt.Printf("✅ [%s] %s (%s) - %dms\n", node.Country, node.Name, addr, elapsed)
			}
		}
	},
}
