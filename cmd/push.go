package cmd

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"mutectl/utils"
	"os"
	"path/filepath"

	"github.com/quic-go/quic-go"
	"github.com/spf13/cobra"
)

var rootIP = "<Root IP>" // 루트 부트스트랩 IP 지정
var rootPort = 8787

var pushCmd = &cobra.Command{
	Use:   "push [filepath]",
	Short: "Push a file into MuteNet (and broadcast to all nodes)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		file := args[0]

		cid, err := utils.HashFile(file)
		if err != nil {
			fmt.Println("❌ Error hashing file:", err)
			return
		}

		homeDir, _ := os.UserHomeDir()
		cacheDir := filepath.Join(homeDir, ".mutenet", "cache")
		metaDir := filepath.Join(homeDir, ".mutenet", "meta")

		os.MkdirAll(cacheDir, 0755)
		os.MkdirAll(metaDir, 0755)

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
		io.Copy(dst, src)
		dst.Close()

		filename := filepath.Base(file)
		ext := filepath.Ext(file)
		mimeType := mime.TypeByExtension(ext)
		if mimeType == "" {
			mimeType = "application/octet-stream"
		}
		meta := map[string]string{
			"cid":      cid,
			"filename": filename,
			"ext":      ext,
			"mime":     mimeType,
		}
		metaPath := filepath.Join(metaDir, cid+".json")
		metaFile, _ := os.Create(metaPath)
		json.NewEncoder(metaFile).Encode(meta)
		metaFile.Close()

		fmt.Printf("✅ File pushed with CID: %s\n", cid)
		fmt.Printf("📦 Stored locally at: %s\n", destPath)

		nodes, err := loadNodes()
		if err != nil || len(nodes) == 0 {
			fmt.Println("⚙️ nodes.json not found or empty, using bootstrap node.")
			nodes = []Node{{IP: rootIP, Port: rootPort}}
			writeNodes(nodes)
		}

		for _, n := range nodes {
			go sendFileOverQUIC(n.IP, n.Port, "META", metaPath)
			go sendFileOverQUIC(n.IP, n.Port, "DATA", destPath)
		}
	},
}

func loadNodes() ([]Node, error) {
	homeDir, _ := os.UserHomeDir()
	path := filepath.Join(homeDir, ".mutenet", "nodes.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var nodes []Node
	err = json.Unmarshal(data, &nodes)
	return nodes, err
}

func writeNodes(nodes []Node) {
	homeDir, _ := os.UserHomeDir()
	path := filepath.Join(homeDir, ".mutenet", "nodes.json")
	data, _ := json.MarshalIndent(nodes, "", "  ")
	os.WriteFile(path, data, 0644)
}

func sendFileOverQUIC(ip string, port int, fileType string, filePath string) error {
	addr := fmt.Sprintf("%s:%d", ip, port)
	conf := &tls.Config{InsecureSkipVerify: true, NextProtos: []string{"mutenet-transfer"}}
	sess, err := quic.DialAddr(context.Background(), addr, conf, nil)
	if err != nil {
		return fmt.Errorf("dial failed: %v", err)
	}
	defer sess.CloseWithError(0, "done")

	stream, err := sess.OpenStreamSync(context.Background())
	if err != nil {
		return fmt.Errorf("stream open failed: %v", err)
	}
	defer stream.Close()

	if _, err := stream.Write([]byte(fileType)); err != nil {
		return err
	}

	name := filepath.Base(filePath)
	if _, err := stream.Write([]byte{byte(len(name))}); err != nil {
		return err
	}
	if _, err := stream.Write([]byte(name)); err != nil {
		return err
	}

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = io.Copy(stream, file)
	if err != nil {
		return err
	}

	fmt.Printf("📤 Sent %s to %s (%s)\n", fileType, addr, name)
	return nil
}
