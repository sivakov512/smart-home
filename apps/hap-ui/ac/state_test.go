package ac_test

import (
	"github.com/stretchr/testify/assert"
	"hap-ui/ac"
	"testing"
)

var (
	serializationCases = []struct {
		de *ac.State
		se string
	}{
		{
			de: &ac.State{
				IsActive:           true,
				Mode:               ac.Cool,
				TargetTemperature:  23.0,
				CurrentTemperature: 25.5,
			},
			se: "{\"is_active\":true,\"mode\":\"cool\",\"target_temperature\":23,\"current_temperature\":25.5}",
		},
		{
			de: &ac.State{
				IsActive:           false,
				Mode:               ac.Heat,
				TargetTemperature:  28.2,
				CurrentTemperature: 14.0,
			},
			se: "{\"is_active\":false,\"mode\":\"heat\",\"target_temperature\":28.2,\"current_temperature\":14}",
		},
	}
)

func TestNewStateReturnsExpected(t *testing.T) {
	got := ac.NewState()

	assert.Equal(t, &ac.State{
		IsActive:           false,
		Mode:               ac.Cool,
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
			var got ac.State

			ac.DeserializeState(input, &got)

			assert.Equal(t, tCase.de, &got)
		})
	}
}
