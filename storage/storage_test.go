package storage

import "testing"

func TestGet(t *testing.T) {
	tests := []struct {
		name string
		data map[string]string
		key string
		value string
		exists bool
	}{
		{
			name: "key exists in mapping",
			data: map[string]string{"calgary": "flames"},
			key: "calgary",
			value: "flames",
			exists: true,
		},
		{
			name: "key does not exist in mapping",
			data: map[string]string{"edmonton": "oilers"},
			key: "calgary",
			value: "",
			exists: false,
		},
		{
			name: "edge case: empty string value",
			data: map[string]string{"vancouver": ""},
			key: "vancouver",
			value: "",
			exists: true,
		},
		{
			name: "edge case: empty mapping",
			data: map[string]string{},
			key: "utah",
			value: "",
			exists: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Store{
				data:  tt.data,
			}

			value, exists := s.Get(tt.key)

			if value != tt.value {
				t.Errorf("actual: %q, expected: %q\n", value, tt.value)
			}
			if exists != tt.exists {
				t.Errorf("actual: %v, expected: %v\n", exists, tt.exists)
			}
		})
	}
}