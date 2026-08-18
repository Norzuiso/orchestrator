package models

import "encoding/json"

type ClientInfo struct {
	HasOpenStream bool            `json:"hasOpenStream"`
	Client        json.RawMessage `json:"client"`
	Connections   json.RawMessage `json:"connections"`
}
