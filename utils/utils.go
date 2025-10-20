package utils

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func GetConfigDir() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".mutenet")
}

func IsDevMode() bool {
	path := filepath.Join(GetConfigDir(), "config.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return false
	}
	var cfg map[string]interface{}
	json.Unmarshal(data, &cfg)
	return cfg["dev"] == true
}

func GetMyIP() string {
	// 외부 IP 조회 → 추후 로컬 테스트 시 "127.0.0.1" 써도 무방
	resp, err := http.Get("https://api.ipify.org")
	if err != nil {
		return "127.0.0.1"
	}
	defer resp.Body.Close()
	ip, _ := io.ReadAll(resp.Body)
	return string(ip)
}
