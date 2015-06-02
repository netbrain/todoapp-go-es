package common

import (
	"encoding/json"
)

//CommandMessage is a WS command message
type CommandMessage struct {
	Name string           `json:"name"`
	Data *json.RawMessage `json:"data"`
}

//ErrorMessage is a generic WS error message
type ErrorMessage struct {
	Reason string `json:"reason"`
}

//EventMessage is a WS event message
type EventMessage struct {
	Name    string           `json:"name"`
	Data    *json.RawMessage `json:"data"`
	Version int              `json:"version"`
}
