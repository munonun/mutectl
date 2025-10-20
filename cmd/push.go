package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"mutectl/utils"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var pushCmd = &cobra.Command{
	Use:   "push [filepath]",
	Short: "Push a file into MuteNet",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		file := args[0]

		// 1️⃣ CID 생성
		cid, err := utils.HashFile(file)
		if err != nil {
			fmt.Println("❌ Error hashing file:", err)
			return
		}

		// 2️⃣ 경로 설정
		homeDir, _ := os.UserHomeDir()
		cacheDir := filepath.Join(homeDir, ".mutenet", "cache")
		metaDir := filepath.Join(homeDir, ".mutenet", "meta")

		// 3️⃣ 디렉토리 없으면 생성
		if err := os.MkdirAll(cacheDir, os.ModePerm); err != nil {
			fmt.Println("❌ Failed to create cache directory:", err)
			return
		}
		if err := os.MkdirAll(metaDir, os.ModePerm); err != nil {
			fmt.Println("❌ Failed to create meta directory:", err)
			return
		}

		// 4️⃣ 캐시에 파일 저장
		src, err := os.Open(file)
		if err != nil {
			fmt.Println("❌ Failed to open source file:", err)
			return
		}
		defer src.Close()

		destPath := filepath.Join(cacheDir, cid)
		dst, err := os.Create(destPath)
		if err != nil {
			fmt.Println("❌ Failed to create cached file:", err)
			return
		}
		defer dst.Close()

		if _, err := io.Copy(dst, src); err != nil {
			fmt.Println("❌ Failed to copy file:", err)
			return
		}

		// 5️⃣ 메타데이터 저장
		filename := filepath.Base(file)
		ext := filepath.Ext(file)
		mimeType := mime.TypeByExtension(ext)
		if mimeType == "" {
			mimeType = "application/octet-stream" // fallback
		}

		meta := map[string]string{
			"cid":      cid,
			"filename": filename,
			"ext":      ext,
			"mime":     mimeType,
		}

		metaPath := filepath.Join(metaDir, cid+".json")
		metaFile, err := os.Create(metaPath)
		if err != nil {
			fmt.Println("❌ Failed to create meta file:", err)
			return
		}
		defer metaFile.Close()

		if err := json.NewEncoder(metaFile).Encode(meta); err != nil {
			fmt.Println("❌ Failed to write meta file:", err)
			return
		}

		// ✅ 출력
		fmt.Printf("✅ File pushed with CID: %s\n", cid)
		fmt.Printf("📦 Stored at: %s\n", destPath)
		fmt.Printf("🧠 Metadata saved: %s\n", metaPath)
	},
}
