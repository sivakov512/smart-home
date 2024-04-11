package boiler_test

import (
	"github.com/stretchr/testify/assert"
	"hap-ui/boiler"
	"testing"
)

var (
	serializationCases = []struct {
		de *boiler.State
		se string
	}{
		{
			de: &boiler.State{State: true},
			se: "{\"state\":true}",
		},
		{
			de: &boiler.State{State: false},
			se: "{\"state\":false}",
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
			var got boiler.State

			boiler.DeserializeState(input, &got)

			assert.Equal(t, tCase.de, &got)
		})
	}
}
