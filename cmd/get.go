package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

type MetaInfo struct {
	CID      string `json:"cid"`
	Filename string `json:"filename"`
	Ext      string `json:"ext"`
	Mime     string `json:"mime"`
}

var getCmd = &cobra.Command{
	Use:   "get [CID]",
	Short: "Get file from MuteNet",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cid := args[0]

		// 1️⃣ 홈 디렉토리 경로
		homeDir, _ := os.UserHomeDir()
		cachePath := filepath.Join(homeDir, ".mutenet", "cache", cid)
		metaPath := filepath.Join(homeDir, ".mutenet", "meta", cid+".json")

		// 2️⃣ 캐시 확인
		if _, err := os.Stat(cachePath); os.IsNotExist(err) {
			fmt.Println("❌ File not found in cache.")
			return
		}

		// 3️⃣ 메타데이터 로드
		meta := MetaInfo{}
		metaFile, err := os.Open(metaPath)
		if err != nil {
			fmt.Println("⚠️  No metadata found, restoring as raw CID file.")
			meta.Filename = cid
		} else {
			defer metaFile.Close()
			if err := json.NewDecoder(metaFile).Decode(&meta); err != nil {
				fmt.Println("⚠️  Failed to read metadata, using CID as filename.")
				meta.Filename = cid
			}
		}

		// 4️⃣ 복원 파일 경로
		destPath := filepath.Join(".", meta.Filename)

		// 5️⃣ 복사
		src, err := os.Open(cachePath)
		if err != nil {
			fmt.Println("❌ Failed to open cached file:", err)
			return
		}
		defer src.Close()

		dst, err := os.Create(destPath)
		if err != nil {
			fmt.Println("❌ Failed to create destination file:", err)
			return
		}
		defer dst.Close()

		if _, err := io.Copy(dst, src); err != nil {
			fmt.Println("❌ Failed to copy file:", err)
			return
		}

		fmt.Printf("✅ File restored as: %s\n", destPath)
	},
}
