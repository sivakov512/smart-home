package boiler

import (
	"hap-ui/common"
)

type Config struct {
	Manufacturer string
	Name         string
	MQTT         *common.MQTTConfig
}
