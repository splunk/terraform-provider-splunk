package models

import (
	"encoding/json"
	"testing"
)

func TestHttpEventCollectorObject_UnmarshalUseACK(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  FlexInt
	}{
		{name: "cloud boolean string false", input: `{"useACK": "false", "index": "main"}`, want: 0},
		{name: "cloud boolean string true", input: `{"useACK": "true", "index": "main"}`, want: 1},
		{name: "cloud boolean false", input: `{"useACK": false, "index": "main"}`, want: 0},
		{name: "cloud boolean true", input: `{"useACK": true, "index": "main"}`, want: 1},
		{name: "enterprise string 0", input: `{"useACK": "0", "index": "main"}`, want: 0},
		{name: "enterprise string 1", input: `{"useACK": "1", "index": "main"}`, want: 1},
		{name: "empty string", input: `{"useACK": "", "index": "main"}`, want: 0},
		{name: "blank string", input: `{"useACK": " ", "index": "main"}`, want: 0},
		{name: "null", input: `{"useACK": null, "index": "main"}`, want: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var obj HttpEventCollectorObject
			if err := json.Unmarshal([]byte(tt.input), &obj); err != nil {
				t.Fatalf("unmarshal %s: %v", tt.input, err)
			}
			if obj.UseACK != tt.want {
				t.Errorf("UseACK = %d, want %d", obj.UseACK, tt.want)
			}
		})
	}
}
