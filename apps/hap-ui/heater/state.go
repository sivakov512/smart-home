package heater

import (
	"hap-ui/common"
)

const (
	Idle common.Mode = "idle"
	Heat             = "heat"
)

func NewState() *common.State {
	return &common.State{
		Mode: Idle,
	}
}
