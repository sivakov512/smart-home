package airconditioner_test

import (
	"github.com/stretchr/testify/assert"
	"hap-ui/airconditioner"
	"testing"
)

var (
	serializationCases = []struct {
		de *airconditioner.State
		se string
	}{
		{
			de: &airconditioner.State{
				IsActive:           true,
				Mode:               airconditioner.Cool,
				TargetTemperature:  23.0,
				CurrentTemperature: 25.5,
			},
			se: "{\"is_active\":true,\"mode\":\"cool\",\"target_temperature\":23,\"current_temperature\":25.5}",
		},
		{
			de: &airconditioner.State{
				IsActive:           false,
				Mode:               airconditioner.Heat,
				TargetTemperature:  28.2,
				CurrentTemperature: 14.0,
			},
			se: "{\"is_active\":false,\"mode\":\"heat\",\"target_temperature\":28.2,\"current_temperature\":14}",
		},
	}
)

func TestNewStateReturnsExpected(t *testing.T) {
	got := airconditioner.NewState()

	assert.Equal(t, &airconditioner.State{
		IsActive:           false,
		Mode:               airconditioner.Cool,
		TargetTemperature:  0.0,
		CurrentTemperature: 0.0,
	}, got)
}

func TestStateSerializesCorrectly(t *testing.T) {
	for _, tCase := range serializationCases {
		t.Run(tCase.se, func(t *testing.T) {
			got := tCase.de.Serialize()

			assert.Equal(t, tCase.se, string(got))
		})
	}
}

func TestStateDeserializesCorrectly(t *testing.T) {
	for _, tCase := range serializationCases {
		t.Run(tCase.se, func(t *testing.T) {
			input := []byte(tCase.se)
			var got airconditioner.State

			airconditioner.DeserializeState(input, &got)

			assert.Equal(t, tCase.de, &got)
		})
	}
}
