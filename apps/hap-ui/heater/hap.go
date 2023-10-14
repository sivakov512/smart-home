package heater

import (
	"github.com/brutella/hap/accessory"
	"github.com/brutella/hap/characteristic"
	"github.com/brutella/hap/service"
)

type HAPService struct {
	*service.S

	Active                      *characteristic.Active
	CurrentHeaterCoolerState    *characteristic.CurrentHeaterCoolerState
	TargetHeaterCoolerState     *characteristic.TargetHeaterCoolerState
	CurrentTemperature          *characteristic.CurrentTemperature
	HeatingThresholdTemperature *characteristic.HeatingThresholdTemperature
}

func NewHAPService() *HAPService {
	s := HAPService{}
	s.S = service.New(service.TypeHeaterCooler)

	s.Active = characteristic.NewActive()
	s.AddC(s.Active.C)

	s.CurrentHeaterCoolerState = characteristic.NewCurrentHeaterCoolerState()
	s.CurrentHeaterCoolerState.ValidVals = []int{
		characteristic.CurrentHeaterCoolerStateInactive,
		characteristic.CurrentHeaterCoolerStateIdle,
		characteristic.CurrentHeaterCoolerStateCooling,
	}
	s.AddC(s.CurrentHeaterCoolerState.C)

	s.TargetHeaterCoolerState = characteristic.NewTargetHeaterCoolerState()
	s.TargetHeaterCoolerState.ValidVals = []int{
		characteristic.TargetHeaterCoolerStateHeat,
	}
	s.AddC(s.TargetHeaterCoolerState.C)

	s.CurrentTemperature = characteristic.NewCurrentTemperature()
	s.AddC(s.CurrentTemperature.C)

	s.HeatingThresholdTemperature = characteristic.NewHeatingThresholdTemperature()
	s.AddC(s.HeatingThresholdTemperature.C)

	return &s
}

type HAPAccessory struct {
	*accessory.A
	Service *HAPService
}

func NewHAPAccessory(info accessory.Info) *HAPAccessory {
	a := HAPAccessory{}
	a.A = accessory.New(info, accessory.TypeAirConditioner)

	a.Service = NewHAPService()
	a.AddS(a.Service.S)

	return &a
}

