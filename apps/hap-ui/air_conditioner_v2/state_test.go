package airconditionerv2_test

import (
	"github.com/stretchr/testify/assert"
	"hap-ui/air_conditioner_v2"
	"sync"
	"testing"
)

func TestNewStateReturnsExpected(t *testing.T) {
	got := airconditionerv2.NewState()

	assert.Equal(t, &airconditionerv2.State{
		Mutex:              sync.Mutex{},
		IsActive:           false,
		Mode:               airconditionerv2.Cool,
		TargetTemperature:  0.0,
		CurrentTemperature: 0.0,
	}, got)
}

func TestStateSerializesCorrectly(t *testing.T) {
	tCases := []struct {
		input    *airconditionerv2.State
		expected string
	}{
		{
			input: &airconditionerv2.State{
				Mutex:              sync.Mutex{},
				IsActive:           true,
				Mode:               airconditionerv2.Cool,
				TargetTemperature:  23.0,
				CurrentTemperature: 25.5,
			},
			expected: "{\"is_active\":true,\"mode\":\"cool\",\"target_temperature\":23,\"current_temperature\":25.5}",
		},
		{
			input: &airconditionerv2.State{
				Mutex:              sync.Mutex{},
				IsActive:           false,
				Mode:               airconditionerv2.Heat,
				TargetTemperature:  28.2,
				CurrentTemperature: 14.0,
			},
			expected: "{\"is_active\":false,\"mode\":\"heat\",\"target_temperature\":28.2,\"current_temperature\":14}",
		},
	}

	for _, tCase := range tCases {
		t.Run(tCase.expected, func(t *testing.T) {
			got := tCase.input.Serialize()

			assert.Equal(t, tCase.expected, string(got))
		})
	}
}
