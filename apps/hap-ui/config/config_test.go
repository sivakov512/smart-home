package config_test

import (
	"github.com/stretchr/testify/assert"
	"hap-ui/ac"
	"hap-ui/boiler"
	"hap-ui/common"
	"hap-ui/config"
	"hap-ui/heater"
	"testing"
)

func TestLoadsConfigCorrectly(t *testing.T) {
	got, err := config.LoadConfig("../testing_fixtures/valid_config.toml")

	assert.Nil(t, err)
	assert.Equal(t, &config.Config{
		Broker: "tcp://iot.eclipse.org:1883",
		PIN:    "11122333",
		AC: &ac.Config{
			Manufacturer: "Midea",
			Name:         "AC",
			Cooling: &common.TemperatureConfig{
				Min:  9.5,
				Max:  31.5,
				Step: 0.8,
			},
			Heating: &common.TemperatureConfig{
				Min:  0.5,
				Max:  31.5,
				Step: 0.8,
			},
			MQTT: &common.MQTTConfig{
				UpdateTopic: "ac/update/LivingRoom",
				StatusTopic: "ac/status/LivingRoom",
			},
		},
		Heater: &heater.Config{
			Manufacturer: "Midea",
			Name:         "Heater",
			Heating: &common.TemperatureConfig{
				Min:  0.5,
				Max:  31.5,
				Step: 0.8,
			},
			MQTT: &common.MQTTConfig{
				UpdateTopic: "heater/update/Kitchen",
				StatusTopic: "heater/status/Kitchen",
			},
		},
		Boiler: &boiler.Config{
			Manufacturer: "Nikita Sivakov",
			Name:         "Boiler",
			MQTT: &common.MQTTConfig{
				StatusTopic: "home/bathroom/boiler",
				UpdateTopic: "home/bathroom/boiler/set",
			},
		},
	}, got)

}

func TestErroredForWrongConfig(t *testing.T) {
	got, err := config.LoadConfig("../testing_fixtures/invalid_config.toml")

	assert.NotNil(t, err)
	assert.Nil(t, got)
}
