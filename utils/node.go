package utils

type Node struct {
	IP      string `json:"ip"`      // Required
	Port    int    `json:"port"`    // Required
	Country string `json:"country"` // Optional
	ID      string `json:"id"`      // Optional
	Name    string `json:"name"`    // Optional
}
