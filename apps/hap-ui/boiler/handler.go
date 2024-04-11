package boiler

import (
	"github.com/brutella/hap/accessory"
	"github.com/eclipse/paho.mqtt.golang"
	"hap-ui/common"
	"net/http"
)

type Handler struct {
	HAPAccessory *HAPAccessory
	config       *Config
	mqttClient   mqtt.Client
	state        *StateGuard
}

func NewHandler(c *Config, mqttClient mqtt.Client) *Handler {
	h := Handler{
		HAPAccessory: NewHAPAccessory(accessory.Info{
			Name:         c.Name,
			Manufacturer: c.Manufacturer,
		}),
		config:     c,
		mqttClient: mqttClient,
	}

	h.setInitialState()

	if token := mqttClient.Subscribe(c.MQTT.StatusTopic, 0, h.handleMQTTMessage); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	service := h.HAPAccessory.Service

	service.On.OnValueRemoteUpdate(h.handleUpdateOn)
	service.On.ValueRequestFunc = h.handleFetchOn

	return &h
}

func (h *Handler) setInitialState() {
	h.state.M.Lock()
	defer h.state.M.Unlock()

	h.state.S.State = false
}

func (h *Handler) handleMQTTMessage(_ mqtt.Client, m mqtt.Message) {
	h.state.M.Lock()
	defer h.state.M.Unlock()

	DeserializeState(m.Payload(), h.state.S)

	h.HAPAccessory.Service.On.SetValue(h.state.S.State)
}

func (h *Handler) handleUpdateOn(v bool) {
	h.state.M.Lock()
	defer h.state.M.Unlock()

	h.state.S.State = v

	serialized := h.state.S.Serialize()
	h.mqttClient.Publish(h.config.MQTT.UpdateTopic, 0, false, serialized)
}

func (h *Handler) handleFetchOn(req *http.Request) (interface{}, int) {
	h.state.M.Lock()
	defer h.state.M.Unlock()

	return h.state.S.State, common.HAPOK
}
