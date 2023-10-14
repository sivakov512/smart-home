package heater

import (
	"github.com/brutella/hap/accessory"
	"github.com/brutella/hap/characteristic"
	"github.com/eclipse/paho.mqtt.golang"
	"hap-ui/common"
	"net/http"
	"sync"
)

const (
	HAPOK = 0
)

type Handler struct {
	HAPAccessory *HAPAccessory
	state        *common.StateGuard
	mqttClient   mqtt.Client
	config       *Config
}

func NewHandler(c *Config, mqttClient mqtt.Client) *Handler {
	h := Handler{
		HAPAccessory: NewHAPAccessory(accessory.Info{
			Name:         c.Name,
			Manufacturer: c.Manufacturer,
		}),
		state: &common.StateGuard{
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

	service := h.HAPAccessory.Service

	service.Active.OnValueRemoteUpdate(h.handleUpdateIsActive)
	service.Active.ValueRequestFunc = h.handleFetchIsActive

	service.CurrentHeaterCoolerState.ValueRequestFunc = h.handleFetchCurrentMode

	service.TargetHeaterCoolerState.ValueRequestFunc = func(_ *http.Request) (interface{}, int) {
		return characteristic.TargetHeaterCoolerStateHeat, HAPOK
	}

	service.CurrentTemperature.ValueRequestFunc = h.handleFetchCurrentTemperature

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
	h.state.S.Mode = Idle
	h.state.S.CurrentTemperature = h.config.Heating.Min
	h.state.S.TargetTemperature = h.config.Heating.Min
}

func (h *Handler) handleMQTTMessage(_ mqtt.Client, m mqtt.Message) {
	h.state.M.Lock()
	defer h.state.M.Unlock()

	common.DeserializeState(m.Payload(), h.state.S)

	service := h.HAPAccessory.Service

	service.Active.SetValue(func() int {
		v := characteristic.ActiveInactive
		if h.state.S.IsActive {
			v = characteristic.ActiveActive
		}

		return v

	}())

	service.CurrentHeaterCoolerState.SetValue(modeToCurrentState(h.state.S.Mode))
	service.TargetHeaterCoolerState.SetValue(characteristic.TargetHeaterCoolerStateHeat)

	service.CurrentTemperature.SetValue(h.state.S.CurrentTemperature)

	service.HeatingThresholdTemperature.SetValue(h.state.S.TargetTemperature)
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

func modeToCurrentState(mode common.Mode) int {
	var state int

	switch mode {
	case Idle:
		state = characteristic.CurrentHeaterCoolerStateIdle
	default:
		state = characteristic.CurrentHeaterCoolerStateHeating
	}

	return state
}
