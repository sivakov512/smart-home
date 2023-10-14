package common

import (
	"encoding/json"
	"sync"
)

type Mode string

type State struct {
	IsActive           bool    `json:"is_active"`
	Mode               Mode    `json:"mode"`
	TargetTemperature  float64 `json:"target_temperature"`
	CurrentTemperature float64 `json:"current_temperature"`
}

type StateGuard struct {
	M sync.Mutex
	S *State
}

func (s *State) Serialize() []byte {
	c, err := json.Marshal(s)
	if err != nil {
		panic(err)
	}

	return c
}

func DeserializeState(i []byte, s *State) {
	json.Unmarshal(i, &s)
}
