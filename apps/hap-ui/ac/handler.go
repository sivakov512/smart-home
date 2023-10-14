package ac

import (
	"github.com/brutella/hap/accessory"
	"github.com/brutella/hap/characteristic"
	"github.com/eclipse/paho.mqtt.golang"
	"hap-ui/common"
	"net/http"
)

const (
	Cool common.Mode = "cool"
	Heat             = "heat"
)

type Handler struct {
	HAPAccessory *HAPAccessory
	config       *Config
	*common.Handler
}

func NewHandler(c *Config, mqttClient mqtt.Client) *Handler {
	h := Handler{
		HAPAccessory: NewHAPAccessory(accessory.Info{
			Name:         c.Name,
			Manufacturer: c.Manufacturer,
		}),
		config:  c,
		Handler: common.NewHandler(c.MQTT, mqttClient),
	}

	h.setInitialState()

	if token := mqttClient.Subscribe(c.MQTT.StatusTopic, 0, h.handleMQTTMessage); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	service := h.HAPAccessory.Service

	service.Active.OnValueRemoteUpdate(h.HandleUpdateIsActive)
	service.Active.ValueRequestFunc = h.HandleFetchIsActive

	service.CurrentHeaterCoolerState.ValueRequestFunc = h.handleFetchCurrentMode

	service.TargetHeaterCoolerState.OnValueRemoteUpdate(h.handleUpdateTargetMode)
	service.TargetHeaterCoolerState.ValueRequestFunc = h.handleFetchTargetMode

	service.CurrentTemperature.ValueRequestFunc = h.HandleFetchCurrentTemperature

	service.CoolingThresholdTemperature.SetMinValue(c.Cooling.Min)
	service.CoolingThresholdTemperature.SetMaxValue(c.Cooling.Max)
	service.CoolingThresholdTemperature.SetStepValue(c.Cooling.Step)
	service.CoolingThresholdTemperature.OnValueRemoteUpdate(h.HandleUpdateTargetTemperature)
	service.CoolingThresholdTemperature.ValueRequestFunc = h.HandleFetchTargetTemperature

	service.HeatingThresholdTemperature.SetMinValue(c.Heating.Min)
	service.HeatingThresholdTemperature.SetMaxValue(c.Heating.Max)
	service.HeatingThresholdTemperature.SetStepValue(c.Heating.Step)
	service.HeatingThresholdTemperature.OnValueRemoteUpdate(h.HandleUpdateTargetTemperature)
	service.HeatingThresholdTemperature.ValueRequestFunc = h.HandleFetchTargetTemperature

	return &h
}

func (h *Handler) setInitialState() {
	h.State.M.Lock()
	defer h.State.M.Unlock()

	h.State.S.IsActive = true
	h.State.S.Mode = Cool
	h.State.S.CurrentTemperature = h.config.Cooling.Min
	h.State.S.TargetTemperature = h.config.Cooling.Min
}

func (h *Handler) handleMQTTMessage(_ mqtt.Client, m mqtt.Message) {
	h.State.M.Lock()
	defer h.State.M.Unlock()

	common.DeserializeState(m.Payload(), h.State.S)

	service := h.HAPAccessory.Service

	service.Active.SetValue(func() int {
		v := characteristic.ActiveInactive
		if h.State.S.IsActive {
			v = characteristic.ActiveActive
		}

		return v

	}())

	mode := h.State.S.Mode
	service.CurrentHeaterCoolerState.SetValue(modeToCurrentState(mode))
	service.TargetHeaterCoolerState.SetValue(modeToTargetState(mode))

	service.CurrentTemperature.SetValue(h.State.S.CurrentTemperature)

	targetTemperature := h.State.S.TargetTemperature
	service.CoolingThresholdTemperature.SetValue(targetTemperature)
	service.HeatingThresholdTemperature.SetValue(targetTemperature)
}

func (h *Handler) handleFetchCurrentMode(req *http.Request) (interface{}, int) {
	h.State.M.Lock()
	defer h.State.M.Unlock()

	return modeToCurrentState(h.State.S.Mode), common.HAPOK
}

func (h *Handler) handleUpdateTargetMode(v int) {
	h.State.M.Lock()
	defer h.State.M.Unlock()

	var mode common.Mode
	switch v {
	case characteristic.TargetHeaterCoolerStateHeat:
		mode = Heat
	case characteristic.TargetHeaterCoolerStateCool:
		mode = Cool
	}

	h.State.S.Mode = mode

	h.Publish2MQTT()
}

func (h *Handler) handleFetchTargetMode(req *http.Request) (interface{}, int) {
	h.State.M.Lock()
	defer h.State.M.Unlock()

	return modeToTargetState(h.State.S.Mode), common.HAPOK
}

func modeToCurrentState(mode common.Mode) int {
	var State int

	switch mode {
	case Heat:
		State = characteristic.CurrentHeaterCoolerStateHeating
	default:
		State = characteristic.CurrentHeaterCoolerStateCooling
	}

	return State
}

func modeToTargetState(mode common.Mode) int {
	var State int

	switch mode {
	case Heat:
		State = characteristic.TargetHeaterCoolerStateHeat
	case Cool:
		State = characteristic.TargetHeaterCoolerStateCool
	}

	return State
}
