package models

import (
	"encoding/json"
	"testing"
)

func TestSplunkBoolInt_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected SplunkBoolInt
		wantErr  bool
	}{
		// Splunk Cloud returns boolean strings
		{name: "string false", input: `"false"`, expected: 0},
		{name: "string true", input: `"true"`, expected: 1},
		// Splunk Enterprise returns integer strings
		{name: "string 0", input: `"0"`, expected: 0},
		{name: "string 1", input: `"1"`, expected: 1},
		// Direct integer values
		{name: "int 0", input: `0`, expected: 0},
		{name: "int 1", input: `1`, expected: 1},
		// Invalid
		{name: "invalid string", input: `"foo"`, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var b SplunkBoolInt
			err := json.Unmarshal([]byte(tt.input), &b)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON(%s) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if !tt.wantErr && b != tt.expected {
				t.Errorf("UnmarshalJSON(%s) = %d, want %d", tt.input, b, tt.expected)
			}
		})
	}
}

func TestSplunkBoolInt_MarshalJSON(t *testing.T) {
	tests := []struct {
		input    SplunkBoolInt
		expected string
	}{
		{input: 0, expected: `"0"`},
		{input: 1, expected: `"1"`},
	}

	for _, tt := range tests {
		data, err := json.Marshal(tt.input)
		if err != nil {
			t.Errorf("MarshalJSON(%d) error = %v", tt.input, err)
			continue
		}
		if string(data) != tt.expected {
			t.Errorf("MarshalJSON(%d) = %s, want %s", tt.input, data, tt.expected)
		}
	}
}

func TestHttpEventCollectorObject_UnmarshalUseACK(t *testing.T) {
	// Simulates the Splunk Cloud API response that was causing the bug
	splunkCloudResponse := `{"useACK": "false", "index": "main"}`

	var obj HttpEventCollectorObject
	if err := json.Unmarshal([]byte(splunkCloudResponse), &obj); err != nil {
		t.Fatalf("failed to unmarshal Splunk Cloud response: %v", err)
	}
	if obj.UseACK != 0 {
		t.Errorf("UseACK = %d, want 0", obj.UseACK)
	}
}
