package airconditioner

import (
	"encoding/json"
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
	return &State{
		inner: &Inner{},
		mutex: sync.Mutex{},
	}
}

func (s *State) Serialize() []byte {
	s.mutex.Lock()
	defer s.mutex.Unlock()

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
}

func (s *State) IsActive() bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return s.inner.IsActive
}

func (s *State) SetActive(v bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.inner.IsActive = v
}

func (s *State) GetMode() Mode {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return s.inner.Mode
}

func (s *State) SetMode(mode Mode) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.inner.Mode = mode
}

func (s *State) SetTargetTemperature(targetTemperature float64) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.inner.TargetTemperature = targetTemperature
}

func (s *State) GetTargetTemperature() float64 {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return s.inner.TargetTemperature
}

func (s *State) SetCurrentTemperature(currentTemperature float64) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.inner.CurrentTemperature = currentTemperature
}

func (s *State) GetCurrentTemperature() float64 {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return s.inner.CurrentTemperature
}
