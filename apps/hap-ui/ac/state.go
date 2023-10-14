package ac

import (
	"hap-ui/common"
)

const (
	Cool common.Mode = "cool"
	Heat             = "heat"
)

func NewState() *common.State {
	return &common.State{
		Mode: Cool,
	}
}
