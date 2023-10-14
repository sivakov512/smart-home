package heater

import (
	"hap-ui/common"
)

type Config struct {
	Manufacturer string
	Name         string
	Heating      *common.TemperatureConfig
	MQTT         *common.MQTTConfig
}
