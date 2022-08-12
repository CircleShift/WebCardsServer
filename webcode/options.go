package webcode

// Put the name of your game here
const GAME_NAME = "No, U"

// Represents an option
// Don't modify
type Option struct {
	Type string `json:"type"`
	Title string `json:"title"`
	Args []interface{} `json:"args"`
}

// GOptionsMsg represents the game options the user can tweak when creating a new game (Sent to user)
type GOptionsMsg struct {
	Name Option
	Password Option
	Hidden Option
	UsePassword Option
}

// GOptions represents the user's current game option values
// There must be parity in naming between GOptions and GOptionsMsg
type GOptions struct {
	Name string
	Password string
	Hidden bool
	UsePassword bool
}

// UOptionsMsg represents the user's options (sent to user)
// must have Name and Color at minimum
type UOptionsMsg struct {
	Name Option
	Color Option
}

// UOptions represents the user's current option values
// There must be parity in naming between UOptions and UOptionsMsg
type UOptions struct {
	Name string
	Color string
}

// InitSettingsMsg defines the initial settings
// Don't modify
type InitSettingsMsg struct {
	User UOptionsMsg `json:"user"`
	Game GOptionsMsg `json:"game"`
}

// Default settings to send to new clients
var DefaultSettingsMsg InitSettingsMsg = InitSettingsMsg{
	UOptionsMsg{
		Option{
			"text",
			"Player Name",
			[]interface{}{"Generic Player", ""},
		},
		Option{
			"color",
			"Player Color",
			[]interface{}{"#ff0000"},
		},
	},
	GOptionsMsg{
		Option{
			"text",
			"Room Name",
			[]interface{}{"", ""},
		},
		Option{
			"text",
			"Room Password",
			[]interface{}{"", ""},
		},
		Option{
			"checkbox",
			"Hidden Game",
			[]interface{}{},
		},
		Option{
			"checkbox",
			"Use Password",
			[]interface{}{},
		},
	},
}

var DefaultUserSettings UOptions = UOptions{
	"Generic Player",
	"#ff0000",
}