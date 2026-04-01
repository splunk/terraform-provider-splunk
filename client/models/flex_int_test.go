package models

import (
	"encoding/json"
	"testing"
)

func TestFlexInt_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    FlexInt
		wantErr bool
	}{
		{"string 1", `"1"`, 1, false},
		{"string 0", `"0"`, 0, false},
		{"bare int 1", `1`, 1, false},
		{"bare int 0", `0`, 0, false},
		{"bool true", `true`, 1, false},
		{"bool false", `false`, 0, false},
		{"invalid string", `"abc"`, 0, true},
		{"null", `null`, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var f FlexInt
			err := json.Unmarshal([]byte(tt.input), &f)
			if (err != nil) != tt.wantErr {
				t.Fatalf("UnmarshalJSON(%s) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
			if !tt.wantErr && f != tt.want {
				t.Errorf("UnmarshalJSON(%s) = %d, want %d", tt.input, f, tt.want)
			}
		})
	}
}

// TestSavedSearchObject_CloudBooleanFields verifies that SavedSearchObject
// correctly deserialises the boolean values Splunk Cloud returns for
// integer fields that Splunk Enterprise returns as string-encoded integers.
// Regression test for https://github.com/splunk/terraform-provider-splunk/issues/130
func TestSavedSearchObject_CloudBooleanFields(t *testing.T) {
	// Splunk Cloud returns booleans for these fields instead of "0"/"1" strings.
	cloudJSON := `{
		"action.email.include.results_link": false,
		"action.email.include.search": false,
		"action.email.include.trigger": true,
		"action.email.include.trigger_time": true,
		"action.email.include.view_link": false,
		"action.email.sendcsv": false,
		"action.snow_event.param.severity": false
	}`

	var obj SavedSearchObject
	if err := json.Unmarshal([]byte(cloudJSON), &obj); err != nil {
		t.Fatalf("Unmarshal cloud response: %v", err)
	}

	if obj.ActionEmailIncludeResultsLink != 0 {
		t.Errorf("ActionEmailIncludeResultsLink: got %d, want 0", obj.ActionEmailIncludeResultsLink)
	}
	if obj.ActionEmailIncludeSearch != 0 {
		t.Errorf("ActionEmailIncludeSearch: got %d, want 0", obj.ActionEmailIncludeSearch)
	}
	if obj.ActionEmailIncludeTrigger != 1 {
		t.Errorf("ActionEmailIncludeTrigger: got %d, want 1", obj.ActionEmailIncludeTrigger)
	}
	if obj.ActionEmailIncludeTriggerTime != 1 {
		t.Errorf("ActionEmailIncludeTriggerTime: got %d, want 1", obj.ActionEmailIncludeTriggerTime)
	}
	if obj.ActionEmailIncludeViewLink != 0 {
		t.Errorf("ActionEmailIncludeViewLink: got %d, want 0", obj.ActionEmailIncludeViewLink)
	}
	if obj.ActionEmailSendCSV != 0 {
		t.Errorf("ActionEmailSendCSV: got %d, want 0", obj.ActionEmailSendCSV)
	}
	if obj.ActionSnowEventParamSeverity != 0 {
		t.Errorf("ActionSnowEventParamSeverity: got %d, want 0", obj.ActionSnowEventParamSeverity)
	}
}

// TestSavedSearchObject_EnterpriseStringFields verifies that SavedSearchObject
// correctly deserialises the string-encoded integers Splunk Enterprise returns.
func TestSavedSearchObject_EnterpriseStringFields(t *testing.T) {
	enterpriseJSON := `{
		"action.email.include.results_link": "1",
		"action.email.include.search": "1",
		"action.email.include.trigger": "1",
		"action.email.include.trigger_time": "1",
		"action.email.include.view_link": "1",
		"action.email.sendcsv": "0",
		"action.snow_event.param.severity": "3"
	}`

	var obj SavedSearchObject
	if err := json.Unmarshal([]byte(enterpriseJSON), &obj); err != nil {
		t.Fatalf("Unmarshal enterprise response: %v", err)
	}

	if obj.ActionEmailIncludeResultsLink != 1 {
		t.Errorf("ActionEmailIncludeResultsLink: got %d, want 1", obj.ActionEmailIncludeResultsLink)
	}
	if obj.ActionEmailIncludeSearch != 1 {
		t.Errorf("ActionEmailIncludeSearch: got %d, want 1", obj.ActionEmailIncludeSearch)
	}
	if obj.ActionEmailIncludeTrigger != 1 {
		t.Errorf("ActionEmailIncludeTrigger: got %d, want 1", obj.ActionEmailIncludeTrigger)
	}
	if obj.ActionEmailIncludeTriggerTime != 1 {
		t.Errorf("ActionEmailIncludeTriggerTime: got %d, want 1", obj.ActionEmailIncludeTriggerTime)
	}
	if obj.ActionEmailIncludeViewLink != 1 {
		t.Errorf("ActionEmailIncludeViewLink: got %d, want 1", obj.ActionEmailIncludeViewLink)
	}
	if obj.ActionEmailSendCSV != 0 {
		t.Errorf("ActionEmailSendCSV: got %d, want 0", obj.ActionEmailSendCSV)
	}
	if obj.ActionSnowEventParamSeverity != 3 {
		t.Errorf("ActionSnowEventParamSeverity: got %d, want 3", obj.ActionSnowEventParamSeverity)
	}
}
