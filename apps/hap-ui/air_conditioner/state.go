package airconditioner

import (
	"fmt"
	"sync"
)

type Mode uint

const (
	Cool Mode = iota
	Heat
)

type State struct {
	isActive           bool
	mode               Mode
	targetTemperature  float64
	currentTemperature float64

	mutex sync.Mutex
}

func NewState() *State {
	return &State{}
}

func (s *State) printState(method string) {
    _, err := fmt.Println(method, s.isActive, s.mode, s.targetTemperature, s.currentTemperature)
    if err != nil {
        panic(err)
    }
}

func (s *State) IsActive() bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()

    s.printState("IsActive")
	return s.isActive
}

func (s *State) SetActive(v bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.isActive = v
    s.printState("SetActive")
}

func (s *State) GetMode() Mode {
	s.mutex.Lock()
	defer s.mutex.Unlock()

    s.printState("GetMode")
	return s.mode
}

func (s *State) SetMode(mode Mode) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.mode = mode
    s.printState("SetMode")
}

func (s *State) SetTargetTemperature(targetTemperature float64) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.targetTemperature = targetTemperature
    s.printState("SetTargetTemperature")
}

func (s *State) GetTargetTemperature() float64 {
	s.mutex.Lock()
	defer s.mutex.Unlock()

    s.printState("GetTargetTemperature")
	return s.targetTemperature
}

func (s *State) SetCurrentTemperature(currentTemperature float64) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.currentTemperature = currentTemperature
    s.printState("SetCurrentTemperature")
}

func (s *State) GetCurrentTemperature() float64 {
	s.mutex.Lock()
	defer s.mutex.Unlock()

    s.printState("SetCurrentTemperature")
	return s.currentTemperature
}
