/*
Copyright 2024 Chia Network Inc.
*/

package chianetwork

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMarshalNetworkOverride(t *testing.T) {
	tests := []struct {
		name     string
		data     interface{}
		expected string
		wantErr  bool
	}{
		{
			name:     "test1",
			data:     "value1",
			expected: `{"test1":"value1"}`,
		},
		{
			name:     "test2",
			data:     123,
			expected: `{"test2":123}`,
		},
		{
			name:     "test3",
			data:     map[string]interface{}{"key": "value"},
			expected: `{"test3":{"key":"value"}}`,
		},
	}

	for _, test := range tests {
		actual, err := marshalNetworkOverride(test.name, test.data)
		require.NoError(t, err)
		require.Equal(t, test.expected, actual)
	}
}
