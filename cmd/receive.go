package cmd

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/quic-go/quic-go"
)

const (
	listenAddr = ":8787"
	cacheDir   = "~/.mutenet/cache"
	metaDir    = "~/.mutenet/meta"
)

func ensureDirs() {
	must(os.MkdirAll(expandPath(cacheDir), 0700))
	must(os.MkdirAll(expandPath(metaDir), 0700))
}

func expandPath(path string) string {
	home, err := os.UserHomeDir()
	must(err)
	return filepath.Clean(filepath.Join(home, path[2:]))
}

func main() {
	ensureDirs()

	listener, err := quic.ListenAddr(listenAddr, generateTLSConfig(), nil)
	must(err)
	fmt.Printf("📥 Listening for push files on %s...\n", listenAddr)

	for {
		conn, err := listener.Accept(context.Background())
		if err != nil {
			log.Printf("Accept failed: %v", err)
			continue
		}
		go handleSession(conn)
	}
}

func handleSession(conn *quic.Conn) {
	stream, err := conn.AcceptStream(context.Background())
	if err != nil {
		log.Printf("Stream accept failed: %v", err)
		return
	}
	defer stream.Close()

	fileTypeBuf := make([]byte, 4)
	_, err = io.ReadFull(stream, fileTypeBuf)
	if err != nil {
		log.Printf("Header read error: %v", err)
		return
	}

	var dstDir string
	switch string(fileTypeBuf) {
	case "META":
		dstDir = expandPath(metaDir)
	case "DATA":
		dstDir = expandPath(cacheDir)
	default:
		log.Printf("Unknown file type: %s", fileTypeBuf)
		return
	}

	nameLenBuf := make([]byte, 1)
	_, err = io.ReadFull(stream, nameLenBuf)
	if err != nil {
		log.Printf("Name len read error: %v", err)
		return
	}
	nameLen := int(nameLenBuf[0])

	nameBuf := make([]byte, nameLen)
	_, err = io.ReadFull(stream, nameBuf)
	if err != nil {
		log.Printf("Name read error: %v", err)
		return
	}

	dstPath := filepath.Join(dstDir, string(nameBuf))
	f, err := os.Create(dstPath)
	if err != nil {
		log.Printf("File create error: %v", err)
		return
	}
	defer f.Close()

	_, err = io.Copy(f, stream)
	if err != nil {
		log.Printf("File write error: %v", err)
		return
	}

	fmt.Printf("✅ Received and saved %s\n", dstPath)
}

func generateTLSConfig() *tls.Config {
	cert, err := tls.X509KeyPair(testCert, testKey)
	must(err)
	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		NextProtos:   []string{"mutenet-transfer"},
	}
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

var testCert = []byte(`-----BEGIN CERTIFICATE-----
MIIBYTCCAQegAwIBAgIQNKnZ4XcCwvhHVqGzWbGx6DAKBggqhkjOPQQDAjASMRAw
DgYDVQQKEwdNdXRlTmV0MB4XDTI1MTAxNzAwMDAwMFoXDTM1MTAxNDAwMDAwMFow
EjEQMA4GA1UEChMHTXV0ZU5ldDBZMBMGByqGSM49AgEGCCqGSM49AwEHA0IABKsa
OD23+bfYCP4Q8bqEHHCrWT7yoEvfhZ/NQxAN2uZhrFVjXZ+YkZFBTk7+PS+FV3Wq
AduTtBql5tsgNY5uwm+jQjBAMA4GA1UdDwEB/wQEAwIHgDAMBgNVHRMBAf8EAjAA
MAoGCCqGSM49BAMCA0gAMEUCIDUtFnAhzE2rSCbbHY5bNJYyWZ1eG6IPmJXvmFZf
G4Y/AiEAtBtXDojErRA/gVnNjnbXjSVTzIRUR3a0JDRQW0JKUMI=
-----END CERTIFICATE-----`)

var testKey = []byte(`-----BEGIN EC PRIVATE KEY-----
MHcCAQEEIHWe2+lKf2BtDGti1EdXWumZK3iRHCeOAfKlBl8ZFRl/oAoGCCqGSM49
AwEHoUQDQgAEqxo4Pbf5t9gI/hDxuocccKtZPvKgS9+Fn81DEA3a5mGsVWNdn5iR
kUFORv49L4VXdaoB25O0GqXm2yA1jm7Cbw==
-----END EC PRIVATE KEY-----`)
