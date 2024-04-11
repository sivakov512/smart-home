package boiler

import (
	"encoding/json"
	"sync"
)

type Mode string

type State struct {
	State bool `json:"state"`
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
