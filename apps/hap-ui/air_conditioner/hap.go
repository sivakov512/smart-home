package airconditioner
import (
	"github.com/brutella/hap/characteristic"
	"github.com/brutella/hap/accessory"
	"github.com/brutella/hap/service"
)

type HAPService struct {
	*service.S

	Active                   *characteristic.Active
	CurrentHeaterCoolerState *characteristic.CurrentHeaterCoolerState
	TargetHeaterCoolerState  *characteristic.TargetHeaterCoolerState
	CurrentTemperature       *characteristic.CurrentTemperature
	CoolingThresholdTemperature *characteristic.CoolingThresholdTemperature
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
		characteristic.CurrentHeaterCoolerStateHeating,
	}
	s.AddC(s.CurrentHeaterCoolerState.C)

	s.TargetHeaterCoolerState = characteristic.NewTargetHeaterCoolerState()
	s.TargetHeaterCoolerState.ValidVals = []int{
		characteristic.TargetHeaterCoolerStateCool,
		characteristic.TargetHeaterCoolerStateHeat,
	}
	s.AddC(s.TargetHeaterCoolerState.C)

	s.CurrentTemperature = characteristic.NewCurrentTemperature()
	s.AddC(s.CurrentTemperature.C)

	s.CoolingThresholdTemperature = characteristic.NewCoolingThresholdTemperature()
	s.AddC(s.CoolingThresholdTemperature.C)

	s.HeatingThresholdTemperature = characteristic.NewHeatingThresholdTemperature()
	s.AddC(s.HeatingThresholdTemperature.C)

	return &s
}

type HAPAccessory struct {
    *accessory.A
    HAPService *HAPService
}

func NewHAPAccessory(info accessory.Info) *HAPAccessory {
    a := HAPAccessory{}
    a.A = accessory.New(info, accessory.TypeAirConditioner)

    a.HAPService = NewHAPService()
    a.AddS(a.HAPService.S)

    return &a
}
