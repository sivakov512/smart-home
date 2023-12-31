package ac

import (
	"hap-ui/common"
)

type Config struct {
	Manufacturer string
	Name         string
	Cooling      *common.TemperatureConfig
	Heating      *common.TemperatureConfig
	MQTT         *common.MQTTConfig
}
