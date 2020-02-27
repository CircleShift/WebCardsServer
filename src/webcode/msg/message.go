package msg

// This file includes message types for talking to clients.

// Message represents a json message from a websocket client
type Message struct {
	Type string `json:"type"`
	Data string `json:"data"`
}
