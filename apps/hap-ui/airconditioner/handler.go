package airconditioner

import (
	"github.com/brutella/hap/accessory"
	"github.com/brutella/hap/characteristic"
	"github.com/eclipse/paho.mqtt.golang"
	"net/http"
	"sync"
)

const (
	HAPOK = 0
)

type Handler struct {
	HAPAccessory *HAPAccessory
	state        *StateGuard
	mqttClient   mqtt.Client
	config       *Config
}

type StateGuard struct {
	M sync.Mutex
	S *State
}

func NewHandler(c *Config, mqttClient mqtt.Client) *Handler {
	h := Handler{
		HAPAccessory: NewHAPAccessory(accessory.Info{
			Name:         c.Name,
			Manufacturer: c.Manufacturer,
		}),
		state: &StateGuard{
			M: sync.Mutex{},
			S: NewState(),
		},
		mqttClient: mqttClient,
		config:     c,
	}

	h.setInitialState()

	if token := mqttClient.Subscribe(c.MQTT.StatusTopic, 0, h.handleMQTTMessage); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	service := h.HAPAccessory.HAPService

	service.Active.OnValueRemoteUpdate(h.handleUpdateIsActive)
	service.Active.ValueRequestFunc = h.handleFetchIsActive

	service.CurrentHeaterCoolerState.ValueRequestFunc = h.handleFetchCurrentMode

	service.TargetHeaterCoolerState.OnValueRemoteUpdate(h.handleUpdateTargetMode)
	service.TargetHeaterCoolerState.ValueRequestFunc = h.handleFetchTargetMode

	service.CurrentTemperature.ValueRequestFunc = h.handleFetchCurrentTemperature

	service.CoolingThresholdTemperature.SetMinValue(c.Cooling.Min)
	service.CoolingThresholdTemperature.SetMaxValue(c.Cooling.Max)
	service.CoolingThresholdTemperature.SetStepValue(c.Cooling.Step)
	service.CoolingThresholdTemperature.OnValueRemoteUpdate(h.handleUpdateTargetTemperature)
	service.CoolingThresholdTemperature.ValueRequestFunc = h.handleFetchTargetTemperature

	service.HeatingThresholdTemperature.SetMinValue(c.Heating.Min)
	service.HeatingThresholdTemperature.SetMaxValue(c.Heating.Max)
	service.HeatingThresholdTemperature.SetStepValue(c.Heating.Step)
	service.HeatingThresholdTemperature.OnValueRemoteUpdate(h.handleUpdateTargetTemperature)
	service.HeatingThresholdTemperature.ValueRequestFunc = h.handleFetchTargetTemperature

	return &h
}

func (h *Handler) publish2MQTT() {
	serialized := h.state.S.Serialize()
	h.mqttClient.Publish(h.config.MQTT.UpdateTopic, 0, false, serialized)
}

func (h *Handler) setInitialState() {
	h.state.M.Lock()
	defer h.state.M.Unlock()

	h.state.S.IsActive = true
	h.state.S.Mode = Cool
	h.state.S.CurrentTemperature = h.config.Cooling.Min
	h.state.S.TargetTemperature = h.config.Cooling.Min
}

func (h *Handler) handleMQTTMessage(_ mqtt.Client, m mqtt.Message) {
	h.state.M.Lock()
	defer h.state.M.Unlock()

	DeserializeState(m.Payload(), h.state.S)

	service := h.HAPAccessory.HAPService

	service.Active.SetValue(func() int {
		v := characteristic.ActiveInactive
		if h.state.S.IsActive {
			v = characteristic.ActiveActive
		}

		return v

	}())

	mode := h.state.S.Mode
	service.CurrentHeaterCoolerState.SetValue(modeToCurrentState(mode))
	service.TargetHeaterCoolerState.SetValue(modeToTargetState(mode))

	service.CurrentTemperature.SetValue(h.state.S.CurrentTemperature)

	targetTemperature := h.state.S.TargetTemperature
	service.CoolingThresholdTemperature.SetValue(targetTemperature)
	service.HeatingThresholdTemperature.SetValue(targetTemperature)
}

func (h *Handler) handleUpdateIsActive(v int) {
	h.state.M.Lock()
	defer h.state.M.Unlock()

	h.state.S.IsActive = (v == characteristic.ActiveActive)

	h.publish2MQTT()
}

func (h *Handler) handleFetchIsActive(req *http.Request) (interface{}, int) {
	h.state.M.Lock()
	defer h.state.M.Unlock()

	v := characteristic.ActiveInactive
	if h.state.S.IsActive {
		v = characteristic.ActiveActive
	}

	return v, HAPOK
}

func (h *Handler) handleFetchCurrentMode(req *http.Request) (interface{}, int) {
	h.state.M.Lock()
	defer h.state.M.Unlock()

	return modeToCurrentState(h.state.S.Mode), HAPOK
}

func (h *Handler) handleUpdateTargetMode(v int) {
	h.state.M.Lock()
	defer h.state.M.Unlock()

	var mode Mode
	switch v {
	case characteristic.TargetHeaterCoolerStateHeat:
		mode = Heat
	case characteristic.TargetHeaterCoolerStateCool:
		mode = Cool
	}

	h.state.S.Mode = mode

	h.publish2MQTT()
}

func (h *Handler) handleFetchTargetMode(req *http.Request) (interface{}, int) {
	h.state.M.Lock()
	defer h.state.M.Unlock()

	return modeToTargetState(h.state.S.Mode), HAPOK
}

func (h *Handler) handleFetchCurrentTemperature(req *http.Request) (interface{}, int) {
	h.state.M.Lock()
	defer h.state.M.Unlock()

	return h.state.S.CurrentTemperature, HAPOK
}

func (h *Handler) handleUpdateTargetTemperature(v float64) {
	h.state.M.Lock()
	defer h.state.M.Unlock()

	h.state.S.TargetTemperature = v

	h.publish2MQTT()
}

func (h *Handler) handleFetchTargetTemperature(req *http.Request) (interface{}, int) {
	h.state.M.Lock()
	defer h.state.M.Unlock()

	return h.state.S.TargetTemperature, HAPOK
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
