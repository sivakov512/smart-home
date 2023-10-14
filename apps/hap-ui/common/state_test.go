package common_test

import (
	"github.com/stretchr/testify/assert"
	"hap-ui/common"
	"testing"
)

var (
	serializationCases = []struct {
		de *common.State
		se string
	}{
		{
			de: &common.State{
				IsActive:           true,
				Mode:               "idle",
				TargetTemperature:  23.0,
				CurrentTemperature: 25.5,
			},
			se: "{\"is_active\":true,\"mode\":\"idle\",\"target_temperature\":23,\"current_temperature\":25.5}",
		},
		{
			de: &common.State{
				IsActive:           false,
				Mode:               "heat",
				TargetTemperature:  28.2,
				CurrentTemperature: 14.0,
			},
			se: "{\"is_active\":false,\"mode\":\"heat\",\"target_temperature\":28.2,\"current_temperature\":14}",
		},
	}
)

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
			var got common.State

			common.DeserializeState(input, &got)

			assert.Equal(t, tCase.de, &got)
		})
	}
}

