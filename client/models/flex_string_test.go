package models

import (
	"encoding/json"
	"testing"

	"github.com/google/go-querystring/query"
)

func TestFlexString_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    FlexString
		wantErr bool
	}{
		{"string number", `"3"`, "3", false},
		{"bare number", `3`, "3", false},
		{"empty string", `""`, "", false},
		{"bool false", `false`, "", false},
		{"bool true", `true`, "true", false},
		{"null keeps value", `null`, "keep", false},
		{"splunk token", `"$result.severity$"`, "$result.severity$", false},
		{"fractional number", `1.5`, "", true},
		{"object", `{"a":1}`, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := FlexString("keep")
			err := json.Unmarshal([]byte(tt.input), &f)
			if (err != nil) != tt.wantErr {
				t.Fatalf("UnmarshalJSON(%s) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
			if !tt.wantErr && f != tt.want {
				t.Errorf("UnmarshalJSON(%s) = %q, want %q", tt.input, f, tt.want)
			}
		})
	}
}

// TestSavedSearchObject_SnowSeverityEncoding verifies that a Splunk token is
// posted unchanged, "0" is posted, and an empty severity is omitted.
func TestSavedSearchObject_SnowSeverityEncoding(t *testing.T) {
	const key = "action.snow_event.param.severity"
	tests := []struct {
		severity FlexString
		want     string
		present  bool
	}{
		{"$result.severity$", "$result.severity$", true},
		{"3", "3", true},
		{"0", "0", true},
		{"", "", false},
	}
	for _, tt := range tests {
		values, err := query.Values(&SavedSearchObject{ActionSnowEventParamSeverity: tt.severity})
		if err != nil {
			t.Fatalf("query.Values: %v", err)
		}
		if _, ok := values[key]; ok != tt.present {
			t.Errorf("severity %q: key present = %v, want %v", tt.severity, ok, tt.present)
		}
		if got := values.Get(key); got != tt.want {
			t.Errorf("severity %q: encoded %q, want %q", tt.severity, got, tt.want)
		}
	}
}

func TestSavedSearchObject_SnowSeverityToken(t *testing.T) {
	var obj SavedSearchObject
	payload := `{"action.snow_event.param.severity": "$result.severity$", "search": "index=main"}`
	if err := json.Unmarshal([]byte(payload), &obj); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if obj.ActionSnowEventParamSeverity != "$result.severity$" {
		t.Errorf("ActionSnowEventParamSeverity: got %q", obj.ActionSnowEventParamSeverity)
	}
	if obj.Search != "index=main" {
		t.Errorf("Search: got %q", obj.Search)
	}
}
