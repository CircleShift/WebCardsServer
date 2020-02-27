package msg

import (
	"card"
	"encoding/json"
)

type packInfo struct {
	Name  string   `json:"name"`
	Suits []string `json:"suits"`
	Cards int      `json:"cards"`
}

func createPList(p []card.Pack) (Message, error) {
	infos := []packInfo{}

	for _, pack := range p {
		infos = append(infos, packInfo{Name: pack.Name, Suits: pack.GetSuits(), Cards: pack.CardCount()})
	}

	dat, err := json.Marshal(infos)

	if err != nil {
		return Message{}, err
	}

	return Message{Type: "plist", Data: string(dat)}, nil
}
