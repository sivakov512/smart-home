package airconditioner

import (
	"encoding/json"
	"fmt"
	"sync"
)

type Mode string

const (
	Cool Mode = "cool"
	Heat      = "heat"
)

type State struct {
	inner *Inner
	mutex sync.Mutex
}

type Inner struct {
	IsActive           bool    `json:"is_active"`
	Mode               Mode    `json:"mode"`
	TargetTemperature  float64 `json:"target_temperature"`
	CurrentTemperature float64 `json:"current_temperature"`
}

func NewState() *State {
    return &State{inner: &Inner{IsActive: false, Mode: Cool, TargetTemperature: 10, CurrentTemperature: 20}, mutex: sync.Mutex{}}
}

func (s *State) printState(method string) {
	_, err := fmt.Println(method, string(s.Serialize()))
	if err != nil {
		panic(err)
	}
}

func (s *State) Serialize() []byte {
	// s.mutex.Lock()
	// defer s.mutex.Unlock()

	c, err := json.Marshal(s.inner)
	if err != nil {
		panic(err)
	}

	return c
}

func (s *State) ReplaceState(c []byte) {
    s.mutex.Lock()
    defer s.mutex.Unlock()

    var state Inner
    json.Unmarshal(c, &state)
    s.inner = &state
    s.printState("ReplaceState")
}

func (s *State) IsActive() bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.printState("IsActive")
	return s.inner.IsActive
}

func (s *State) SetActive(v bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.inner.IsActive = v
	s.printState("SetActive")
}

func (s *State) GetMode() Mode {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.printState("GetMode")
	return s.inner.Mode
}

func (s *State) SetMode(mode Mode) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.inner.Mode = mode
	s.printState("SetMode")
}

func (s *State) SetTargetTemperature(targetTemperature float64) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.inner.TargetTemperature = targetTemperature
	s.printState("SetTargetTemperature")
}

func (s *State) GetTargetTemperature() float64 {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.printState("GetTargetTemperature")
	return s.inner.TargetTemperature
}

func (s *State) SetCurrentTemperature(currentTemperature float64) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.inner.CurrentTemperature = currentTemperature
	s.printState("SetCurrentTemperature")
}

func (s *State) GetCurrentTemperature() float64 {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.printState("SetCurrentTemperature")
	return s.inner.CurrentTemperature
}
