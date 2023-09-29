package airconditioner

import (
	"github.com/brutella/hap/accessory"
	"github.com/brutella/hap/characteristic"
	"github.com/eclipse/paho.mqtt.golang"
	"net/http"
)

type Handler struct {
	HAPAccessory *HAPAccessory
	State        *State
	MQTTClient   mqtt.Client
}

func NewHandler(config *Config) *Handler {
	mqttOpts := mqtt.NewClientOptions().AddBroker(config.MQTT.Broker)
	mqttClient := mqtt.NewClient(mqttOpts)
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	handler := Handler{
		HAPAccessory: NewHAPAccessory(accessory.Info{
			Name:         config.Name,
			Manufacturer: config.Manufacturer,
		}),
		State:      NewState(),
		MQTTClient: mqttClient,
	}

	handler.setupStates(config)
	handler.setupCurrentTemperature(config)
	handler.setupThresholds(config)

    handler.setupMQTTSubsriber(config)

	return &handler
}

func (h *Handler) setupStates(config *Config) {
	active := h.HAPAccessory.HAPService.Active
	active.OnValueRemoteUpdate(func(v int) {
		h.State.SetActive(v == characteristic.ActiveActive)
		h.MQTTClient.Publish(config.MQTT.UpdateTopic, 0, false, h.State.Serialize())
	})
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
		h.MQTTClient.Publish(config.MQTT.UpdateTopic, 0, false, h.State.Serialize())
	})
	currentState.ValueRequestFunc = func(req *http.Request) (interface{}, int) {
		return modeToCurrentState(h.State.GetMode()), 0
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
		h.MQTTClient.Publish(config.MQTT.UpdateTopic, 0, false, h.State.Serialize())
	})
	targetState.ValueRequestFunc = func(req *http.Request) (interface{}, int) {
		return modeToTargetState(h.State.GetMode()), 0
	}
}

func (h *Handler) setupCurrentTemperature(config *Config) {
	currentTemperature := h.HAPAccessory.HAPService.CurrentTemperature

	currentTemperature.OnValueRemoteUpdate(func(v float64) {
		h.State.SetCurrentTemperature(v)
		h.MQTTClient.Publish(config.MQTT.UpdateTopic, 0, false, h.State.Serialize())
	})
	currentTemperature.ValueRequestFunc = func(req *http.Request) (interface{}, int) {
		return h.State.GetCurrentTemperature(), 0
	}
}

func (h *Handler) setupThresholds(config *Config) {
	ValueRequestFunc := func(req *http.Request) (interface{}, int) {
		return h.State.GetTargetTemperature(), 0
	}
	OnValueRemoteUpdate := func(v float64) {
		h.State.SetTargetTemperature(v)
		h.MQTTClient.Publish(config.MQTT.UpdateTopic, 0, false, h.State.Serialize())
	}

	CoolingThreshold := h.HAPAccessory.HAPService.CoolingThresholdTemperature

	CoolingThreshold.SetMinValue(config.Temperature.Min)
	CoolingThreshold.SetMaxValue(config.Temperature.Max)
	CoolingThreshold.SetStepValue(config.Temperature.Step)
	CoolingThreshold.OnValueRemoteUpdate(OnValueRemoteUpdate)
	CoolingThreshold.ValueRequestFunc = ValueRequestFunc

	HeatingThreshold := h.HAPAccessory.HAPService.HeatingThresholdTemperature
	HeatingThreshold.SetMinValue(config.Temperature.Min)
	HeatingThreshold.SetMaxValue(config.Temperature.Max)
	HeatingThreshold.SetStepValue(config.Temperature.Step)
	HeatingThreshold.OnValueRemoteUpdate(OnValueRemoteUpdate)
	HeatingThreshold.ValueRequestFunc = ValueRequestFunc
}

func (h *Handler) setupMQTTSubsriber(config *Config) {
	handler := func(c mqtt.Client, m mqtt.Message) {
		h.State.ReplaceState(m.Payload())

		service := h.HAPAccessory.HAPService

		service.Active.SetValue(func() int {
			val := characteristic.ActiveInactive
			if h.State.IsActive() {
				val = characteristic.ActiveActive
			}

			return val

		}())

		mode := h.State.GetMode()
		service.CurrentHeaterCoolerState.SetValue(modeToCurrentState(mode))
		service.TargetHeaterCoolerState.SetValue(modeToTargetState(mode))

        service.CurrentTemperature.SetValue(h.State.GetCurrentTemperature())

		targetTemperature := h.State.GetTargetTemperature()
		service.CoolingThresholdTemperature.SetValue(targetTemperature)
		service.HeatingThresholdTemperature.SetValue(targetTemperature)
	}
	if token := h.MQTTClient.Subscribe(config.MQTT.StatusTopic, 0, handler); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
}

func modeToCurrentState(mode Mode) int {
	var state int
	switch mode {
	case Heat:
		state = characteristic.CurrentHeaterCoolerStateHeating
	case Cool:
		state = characteristic.CurrentHeaterCoolerStateCooling
	default:
		state = characteristic.CurrentHeaterCoolerStateCooling
	}

	return state
}

func modeToTargetState(mode Mode) int {
	var state int
	switch mode {
	case Heat:
		state = characteristic.TargetHeaterCoolerStateHeat
	case Cool:
		state = characteristic.TargetHeaterCoolerStateCool
	}

	return state
}
