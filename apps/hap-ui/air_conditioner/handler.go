package airconditioner

import (
	"github.com/brutella/hap/accessory"
	"github.com/brutella/hap/characteristic"
	"net/http"
)

type Handler struct {
	HAPAccessory *HAPAccessory
	State        *State
}

func NewHandler(config *Config) *Handler {
	handler := Handler{
		HAPAccessory: NewHAPAccessory(accessory.Info{
			Name:         config.Name,
			Manufacturer: config.Manufacturer,
		}),
		State: NewState(),
	}

	handler.setupStates()
	handler.setupCurrentTemperature()
	handler.setupThresholds(config)

	return &handler
}

func (h *Handler) setupStates() {
	active := h.HAPAccessory.HAPService.Active
	active.OnValueRemoteUpdate(func(v int) { h.State.SetActive(v == characteristic.ActiveActive) })
	active.ValueRequestFunc = func(req *http.Request) (interface{}, int) {
		val := characteristic.ActiveInactive
		if h.State.IsActive() {
			val = characteristic.ActiveActive
		}

		return val, 0
	}

	currentState := h.HAPAccessory.HAPService.CurrentHeaterCoolerState
	currentState.OnValueRemoteUpdate(func(v int) {
		var mode Mode
		switch v {
		case characteristic.CurrentHeaterCoolerStateHeating:
			mode = Heat
		case characteristic.CurrentHeaterCoolerStateCooling:
			mode = Cool
		}

		h.State.SetMode(mode)
	})
	currentState.ValueRequestFunc = func(req *http.Request) (interface{}, int) {
		var state int
		switch h.State.GetMode() {
		case Heat:
			state = characteristic.CurrentHeaterCoolerStateHeating
		case Cool:
			state = characteristic.CurrentHeaterCoolerStateCooling
		default:
			state = characteristic.CurrentHeaterCoolerStateIdle
		}

		return state, 0
	}

	targetState := h.HAPAccessory.HAPService.TargetHeaterCoolerState
	targetState.OnValueRemoteUpdate(func(v int) {
		var mode Mode
		switch v {
		case characteristic.TargetHeaterCoolerStateHeat:
			mode = Heat
		case characteristic.TargetHeaterCoolerStateCool:
			mode = Cool
		}

		h.State.SetMode(mode)
	})
	targetState.ValueRequestFunc = func(req *http.Request) (interface{}, int) {
		var state int
		switch h.State.GetMode() {
		case Heat:
			state = characteristic.TargetHeaterCoolerStateHeat
		case Cool:
			state = characteristic.TargetHeaterCoolerStateCool
		}

		return state, 0
	}
}

func (h *Handler) setupCurrentTemperature() {
	currentTemperature := h.HAPAccessory.HAPService.CurrentTemperature

	currentTemperature.OnValueRemoteUpdate(h.State.SetCurrentTemperature)
	currentTemperature.ValueRequestFunc = func(req *http.Request) (interface{}, int) {
		return h.State.GetCurrentTemperature(), 0
	}
}

func (h *Handler) setupThresholds(config *Config) {
	ValueRequestFunc := func(req *http.Request) (interface{}, int) {
		return h.State.GetTargetTemperature(), 0
	}

	CoolingThreshold := h.HAPAccessory.HAPService.CoolingThresholdTemperature

	CoolingThreshold.SetMinValue(config.Temperature.Min)
	CoolingThreshold.SetMaxValue(config.Temperature.Max)
	CoolingThreshold.SetStepValue(config.Temperature.Step)
	CoolingThreshold.OnValueRemoteUpdate(h.State.SetTargetTemperature)
	CoolingThreshold.ValueRequestFunc = ValueRequestFunc

	HeatingThreshold := h.HAPAccessory.HAPService.HeatingThresholdTemperature
	HeatingThreshold.SetMinValue(config.Temperature.Min)
	HeatingThreshold.SetMaxValue(config.Temperature.Max)
	HeatingThreshold.SetStepValue(config.Temperature.Step)
	HeatingThreshold.OnValueRemoteUpdate(h.State.SetTargetTemperature)
	HeatingThreshold.ValueRequestFunc = ValueRequestFunc
}
