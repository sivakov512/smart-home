package common

import (
	"github.com/brutella/hap/characteristic"
	"github.com/eclipse/paho.mqtt.golang"
	"net/http"
	"sync"
)

const (
	HAPOK = 0
)

type Handler struct {
	State      *StateGuard
	MQTTClient mqtt.Client
	MQTTConfig *MQTTConfig
}

func NewHandler(mqttConfig *MQTTConfig, mqttClient mqtt.Client) *Handler {
	return &Handler{
		State: &StateGuard{
			M: sync.Mutex{},
			S: &State{},
		},
		MQTTClient: mqttClient,
		MQTTConfig: mqttConfig,
	}
}

func (h *Handler) Publish2MQTT() {
	serialized := h.State.S.Serialize()
	h.MQTTClient.Publish(h.MQTTConfig.UpdateTopic, 0, false, serialized)
}

func (h *Handler) HandleUpdateIsActive(v int) {
	h.State.M.Lock()
	defer h.State.M.Unlock()

	h.State.S.IsActive = (v == characteristic.ActiveActive)

	h.Publish2MQTT()
}

func (h *Handler) HandleFetchIsActive(req *http.Request) (interface{}, int) {
	h.State.M.Lock()
	defer h.State.M.Unlock()

	v := characteristic.ActiveInactive
	if h.State.S.IsActive {
		v = characteristic.ActiveActive
	}

	return v, HAPOK
}

func (h *Handler) HandleFetchCurrentTemperature(req *http.Request) (interface{}, int) {
	h.State.M.Lock()
	defer h.State.M.Unlock()

	return h.State.S.CurrentTemperature, HAPOK
}

func (h *Handler) HandleUpdateTargetTemperature(v float64) {
	h.State.M.Lock()
	defer h.State.M.Unlock()

	h.State.S.TargetTemperature = v

	h.Publish2MQTT()
}

func (h *Handler) HandleFetchTargetTemperature(req *http.Request) (interface{}, int) {
	h.State.M.Lock()
	defer h.State.M.Unlock()

	return h.State.S.TargetTemperature, HAPOK
}
