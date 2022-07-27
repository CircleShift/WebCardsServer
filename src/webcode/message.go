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
	Data interface{} `json:"data"`
}

// ChatMessage facilitates sending a chat message to a player
type ChatMessage struct {
	UserName string `json:"user"`
	Text string `json:"text"`
	Color string `json:"color"`
	ChatID string `json:"channel"`
	Server bool `json:"server"`
}

// AddChatMessage notifies a user that they have access to a chat
type AddChatMessage struct {
	Name string `json:"name"`
	ChatID string `json:"id"`
	Follow bool `json:"follow"`
}
