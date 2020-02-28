package webcode

// This file includes message types for talking to clients.

// RecieveMessage represents a json message from a websocket client
type RecieveMessage struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

// SendMessage represents a json message to the websocket client
type SendMessage struct {
	Type string `json:"type"`
	Data {}interface `json:"data"`
}
