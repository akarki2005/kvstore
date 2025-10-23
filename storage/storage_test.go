package storage

import "testing"

// To mock the logfile

type mockLogFile struct{}

func (mockLogFile) WriteString(s string) (int, error) {
    return len(s), nil
}

func (mockLogFile) Sync() error {
    return nil
}

func (mockLogFile) Close() error {
	return nil
}

// Tests

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

func TestSet(t *testing.T) {
	tests := []struct {
		name string
		key string
		value string
	}{
		{
			name: "set new key",
			key: "connor",
			value: "mcdavid",
		},
		{
			name: "overwrite existing key",
			key: "connor",
			value: "bedard",
		},
		{
			name: "edge case: set empty string value",
			key: "ovi",
			value: "",
		},
	}

	s := &Store{
				data: map[string]string{},
				logFile: mockLogFile{},
			}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			err := s.Set(tt.key, tt.value)
			if err != nil {
				t.Errorf("Error calling Set(): %v\n", err)
			}

			value, exists := s.data[tt.key]
			if !exists {
				t.Errorf("Error: %q not in store.\n", tt.key)
			}
			if value != tt.value {
				t.Errorf("Error: Expected key %q to have value %q, but got %q\n", tt.key, tt.value, value)
			}

		})
	}
}