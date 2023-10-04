package airconditionerv2_test

import (
	"github.com/stretchr/testify/assert"
	"hap-ui/air_conditioner_v2"
	"testing"
)

var (
	serializationCases = []struct {
		de *airconditionerv2.State
		se string
	}{
		{
			de: &airconditionerv2.State{
				IsActive:           true,
				Mode:               airconditionerv2.Cool,
				TargetTemperature:  23.0,
				CurrentTemperature: 25.5,
			},
			se: "{\"is_active\":true,\"mode\":\"cool\",\"target_temperature\":23,\"current_temperature\":25.5}",
		},
		{
			de: &airconditionerv2.State{
				IsActive:           false,
				Mode:               airconditionerv2.Heat,
				TargetTemperature:  28.2,
				CurrentTemperature: 14.0,
			},
			se: "{\"is_active\":false,\"mode\":\"heat\",\"target_temperature\":28.2,\"current_temperature\":14}",
		},
	}
)

func TestNewStateReturnsExpected(t *testing.T) {
	got := airconditionerv2.NewState()

	assert.Equal(t, &airconditionerv2.State{
		IsActive:           false,
		Mode:               airconditionerv2.Cool,
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
			var got airconditionerv2.State

			airconditionerv2.DeserializeState(input, &got)

			assert.Equal(t, tCase.de, &got)
		})
	}
}
