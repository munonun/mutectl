// dht.go
package cmd

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mutectl/utils"
	"net"
	"os"
	"sync"

	"golang.org/x/crypto/chacha20poly1305"
)

// -------------------- STRUCTS --------------------

type DHT struct {
	Nodes map[string]utils.Node
	Mutex sync.RWMutex
	Key   []byte
}

// -------------------- INIT --------------------
func NewDHT(secretKey []byte) *DHT {
	return &DHT{
		Nodes: make(map[string]utils.Node),
		Key:   secretKey,
	}
}

// -------------------- DHT CORE --------------------
func (d *DHT) Join(node utils.Node) {
	d.Mutex.Lock()
	defer d.Mutex.Unlock()
	d.Nodes[node.ID] = node
}

func (d *DHT) Find(id string) (utils.Node, bool) {
	d.Mutex.RLock()
	defer d.Mutex.RUnlock()
	n, ok := d.Nodes[id]
	return n, ok
}

func (d *DHT) AllNodes() []utils.Node {
	d.Mutex.RLock()
	defer d.Mutex.RUnlock()
	list := []utils.Node{}
	for _, n := range d.Nodes {
		list = append(list, n)
	}
	return list
}

// -------------------- FILE IO --------------------
func (d *DHT) SaveToFile(filepath string) error {
	d.Mutex.RLock()
	defer d.Mutex.RUnlock()

	aead, err := chacha20poly1305.NewX(d.Key)
	if err != nil {
		return err
	}

	plaintext, err := json.Marshal(d.Nodes)
	if err != nil {
		return err
	}

	nonce := make([]byte, chacha20poly1305.NonceSizeX)
	_, err = rand.Read(nonce)
	if err != nil {
		return err
	}

	ciphertext := aead.Seal(nonce, nonce, plaintext, nil)
	return ioutil.WriteFile(filepath, ciphertext, 0600)
}

func (d *DHT) LoadFromFile(filepath string) error {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // silently ignore missing file
		}
		return err
	}
	aead, err := chacha20poly1305.NewX(d.Key)
	if err != nil {
		return err
	}

	if len(data) < chacha20poly1305.NonceSizeX {
		return fmt.Errorf("invalid file format")
	}

	nonce := data[:chacha20poly1305.NonceSizeX]
	ciphertext := data[chacha20poly1305.NonceSizeX:]
	plaintext, err := aead.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return err
	}

	return json.Unmarshal(plaintext, &d.Nodes)
}

// -------------------- EXAMPLE USAGE --------------------
func DHTUsage() {
	key := make([]byte, chacha20poly1305.KeySize)
	rand.Read(key) // 임시 키. 실환경에서는 PQC + XChaCha20 조합 필요

	dht := NewDHT(key)
	dht.LoadFromFile("nodes.json")

	self := utils.Node{
		ID:   "abc123",
		IP:   getLocalIP(),
		Port: 7777,
	}
	dht.Join(self)
	dht.SaveToFile("nodes.json")

	fmt.Println("DHT Nodes:", dht.AllNodes())
}

// -------------------- UTILS --------------------
func getLocalIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "127.0.0.1"
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}
