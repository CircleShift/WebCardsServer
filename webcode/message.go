package webcode

// This file includes message types for talking to clients.

//********************
//* Generic Wrappers *
//********************

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

//**********************
//* Messages for chats *
//**********************

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

//**************************
//* Messages for the Lobby *
//**************************

// StatMessage update's the user's server stat view
type StatsMessage struct {
	UsersOnline int `json:"online"`
	UsersInGame int `json:"ingame"`
	PublicGames int `json:"pubgame"`
}

// Message to let the user know of a public game
type AddGameMessage struct {
	Name string `json:"name"`
	Packs int `json:"packs"`
	GameID string `json:"id"`
	Pass bool `json:"pass"`
}

// Initial message to let the user view all the public games
type GameListMessage struct {
	Game string `json:"game"`
	Packs int `json:"packs"`
	Games []AddGameMessage `json:"games"`
}

//**********************
//* Messages for games *
//**********************

type JoinGameMessage struct {
	GameID string `json:"id"`
	Password string `json:"pass"`
}

type DeckOptions struct {
	Mode string `json:"mode"`
	SelectMode string `json:"smode"`
	SelectCount int `json:"sct"`
	Position [4]float64 `json:"pos"`
}

type NewDeckMessage struct {
	DeckID string `json:"id"`
	Options DeckOptions `json:"options"`
}

type NewCardMessage struct {
	CardID string `json:"id"`
	DeckID string `json:"deck"`
	Data interface{} `json:"data"`
}

type MoveCardMessage struct {
	CardID string `type:"cardID"`
	DeckID string `type:"deckID"`
	Index string `type:"index"`
}

type SwapCardMessage struct {
	CardID string `type:"cardID"`
	NewID string `type:"newID"`
	Data interface{} `type:"data"`
}