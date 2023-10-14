package config_test

import (
	"github.com/stretchr/testify/assert"
	"hap-ui/ac"
	"hap-ui/config"
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
			Cooling: ac.Temperature{
				Min:  9.5,
				Max:  31.5,
				Step: 0.8,
			},
			Heating: ac.Temperature{
				Min:  0.5,
				Max:  31.5,
				Step: 0.8,
			},
			MQTT: ac.MQTT{
				UpdateTopic: "ac/update/LivingRoom",
				StatusTopic: "ac/status/LivingRoom",
			},
		},
	}, got)

}

func TestErroredForWrongConfig(t *testing.T) {
	got, err := config.LoadConfig("../testing_fixtures/invalid_config.toml")

	assert.NotNil(t, err)
	assert.Nil(t, got)
}
