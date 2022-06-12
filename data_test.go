package main

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFailureType_MarshalJSON(t *testing.T) {
	var (
		unknown = Unknown
		down    = Down
		_error  = Error
	)

	tests := map[string]struct {
		t        *FailureType
		expected []byte
	}{
		"Nil FailureType": {
			t:        nil,
			expected: []byte(`null`),
		},
		"Unknown": {
			t:        &unknown,
			expected: []byte(`"unknown"`),
		},
		"Down": {
			t:        &down,
			expected: []byte(`"down"`),
		},
		"Error": {
			t:        &_error,
			expected: []byte(`"error"`),
		},
	}

	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			result, err := json.Marshal(tc.t)
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestFailureType_UnmarshalJSON(t *testing.T) {
	var (
		unknown = Unknown
		down    = Down
		_error  = Error
		zero    = FailureType(0)
	)

	tests := map[string]struct {
		input    []byte
		expected *FailureType
	}{
		"Null": {
			input:    []byte("null"),
			expected: &zero,
		},
		"Unknown": {
			input:    []byte(`"unknown"`),
			expected: &unknown,
		},
		"Down": {
			input:    []byte(`"down"`),
			expected: &down,
		},
		"Error": {
			input:    []byte(`"error"`),
			expected: &_error,
		},
	}

	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			f := new(FailureType)
			assert.NoError(t, json.Unmarshal(tc.input, f))
			assert.Equal(t, *tc.expected, *f)
		})
	}
}
