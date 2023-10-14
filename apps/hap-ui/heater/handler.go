package heater

import (
	"github.com/brutella/hap/accessory"
	"github.com/brutella/hap/characteristic"
	"github.com/eclipse/paho.mqtt.golang"
	"hap-ui/common"
	"net/http"
)

const (
	Idle common.Mode = "idle"
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

	service.TargetHeaterCoolerState.ValueRequestFunc = func(_ *http.Request) (interface{}, int) {
		return characteristic.TargetHeaterCoolerStateHeat, common.HAPOK
	}

	service.CurrentTemperature.ValueRequestFunc = h.HandleFetchCurrentTemperature

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
	h.State.S.Mode = Idle
	h.State.S.CurrentTemperature = h.config.Heating.Min
	h.State.S.TargetTemperature = h.config.Heating.Min
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

	service.CurrentHeaterCoolerState.SetValue(modeToCurrentState(h.State.S.Mode))
	service.TargetHeaterCoolerState.SetValue(characteristic.TargetHeaterCoolerStateHeat)

	service.CurrentTemperature.SetValue(h.State.S.CurrentTemperature)

	service.HeatingThresholdTemperature.SetValue(h.State.S.TargetTemperature)
}

func (h *Handler) handleFetchCurrentMode(req *http.Request) (interface{}, int) {
	h.State.M.Lock()
	defer h.State.M.Unlock()

	return modeToCurrentState(h.State.S.Mode), common.HAPOK
}

func modeToCurrentState(mode common.Mode) int {
	var State int

	switch mode {
	case Idle:
		State = characteristic.CurrentHeaterCoolerStateIdle
	default:
		State = characteristic.CurrentHeaterCoolerStateHeating
	}

	return State
}
