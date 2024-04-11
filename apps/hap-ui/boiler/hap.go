package boiler

import (
	"github.com/brutella/hap/accessory"
	"github.com/brutella/hap/characteristic"
	"github.com/brutella/hap/service"
)

type HAPService struct {
	*service.S

	On *characteristic.On
}

func NewHAPService() *HAPService {
	s := HAPService{}
	s.S = service.New(service.TypeSwitch)

	s.On = characteristic.NewOn()
	s.AddC(s.On.C)

	return &s
}

type HAPAccessory struct {
	*accessory.A
	Service *HAPService
}

func NewHAPAccessory(info accessory.Info) *HAPAccessory {
	a := HAPAccessory{}
	a.A = accessory.New(info, accessory.TypeSwitch)

	a.Service = NewHAPService()
	a.AddS(a.Service.S)

	return &a
}
