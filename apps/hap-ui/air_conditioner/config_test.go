package airconditioner_test

import (
	"github.com/stretchr/testify/assert"
	"hap-ui/air_conditioner"
	"testing"
)

func TestLoadsConfigCorrectly(t *testing.T) {
	got, err := airconditioner.LoadConfig("../testing_fixtures/valid_config.toml")

	assert.Nil(t, err)
	assert.Equal(t, got, &airconditioner.Config{
		Manufacturer: "Midea",
		Name:         "AC",
		PIN:          "11122333",
		Temperature: airconditioner.Temperature{
			Min:  9.5,
			Max:  31.5,
			Step: 0.8,
		},
	})
}

func TestErroredForWrongConfig(t *testing.T) {
	got, err := airconditioner.LoadConfig("../testing_fixtures/invalid_config.toml")

	assert.NotNil(t, err)
	assert.Nil(t, got)
}
