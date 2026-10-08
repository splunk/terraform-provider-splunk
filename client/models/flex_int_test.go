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
		{"string true", `"true"`, 1, false},
		{"string false", `"false"`, 0, false},
		{"string TRUE", `"TRUE"`, 1, false},
		{"string False", `"False"`, 0, false},
		{"bare int 1", `1`, 1, false},
		{"bare int 0", `0`, 0, false},
		{"bool true", `true`, 1, false},
		{"bool false", `false`, 0, false},
		{"empty string", `""`, 0, false},
		{"blank string", `" "`, 0, false},
		{"invalid string", `"abc"`, 0, true},
		{"null", `null`, 0, false},
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

// TestFlexInt_UnmarshalJSONExistingValue verifies that null keeps the current
// value while an empty string sets it to 0.
func TestFlexInt_UnmarshalJSONExistingValue(t *testing.T) {
	for input, want := range map[string]FlexInt{`null`: 3, `""`: 0} {
		f := FlexInt(3)
		if err := json.Unmarshal([]byte(input), &f); err != nil {
			t.Fatalf("UnmarshalJSON(%s) error = %v", input, err)
		}
		if f != want {
			t.Errorf("UnmarshalJSON(%s) = %d, want %d", input, f, want)
		}
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

// TestSavedSearchObject_EmptySeverity verifies that an empty severity, which
// Splunk returns for saved searches when the ServiceNow add-on is installed,
// does not stop decoding of the fields that follow it in the response.
func TestSavedSearchObject_EmptySeverity(t *testing.T) {
	enterpriseJSON := `{
		"action.snow_event.param.severity": "",
		"search": "index=_internal | head 1"
	}`

	var obj SavedSearchObject
	if err := json.Unmarshal([]byte(enterpriseJSON), &obj); err != nil {
		t.Fatalf("Unmarshal enterprise response: %v", err)
	}
	if obj.ActionSnowEventParamSeverity != 0 {
		t.Errorf("ActionSnowEventParamSeverity: got %d, want 0", obj.ActionSnowEventParamSeverity)
	}
	if obj.Search != "index=_internal | head 1" {
		t.Errorf("Search: got %q, want %q", obj.Search, "index=_internal | head 1")
	}
}

// TestSavedSearchObject_Splunk10StringFields verifies that a Splunk 10
// saved search response decodes when boolean and time fields are strings.
// reportIncludeSplunkLogo is "1" or "0". Action maxtime values are durations
// such as "5m". auto_summarize.max_time is a numeric string. An empty
// dispatch.indexedRealtimeOffset decodes as 0 so later fields still decode.
func TestSavedSearchObject_Splunk10StringFields(t *testing.T) {
	enterpriseJSON := `{
		"action.email.reportIncludeSplunkLogo": "1",
		"action.populate_lookup.maxtime": "5m",
		"action.rss.maxtime": "1m",
		"action.script.maxtime": "5m",
		"action.summary_index.maxtime": "5m",
		"auto_summarize.max_time": "3600",
		"dispatch.indexedRealtimeOffset": "",
		"dispatch.indexedRealtimeMinSpan": "",
		"alert.suppress": null,
		"dispatch.indexedRealtime": null,
		"search": "index=_internal | head 1"
	}`

	var obj SavedSearchObject
	if err := json.Unmarshal([]byte(enterpriseJSON), &obj); err != nil {
		t.Fatalf("Unmarshal Splunk 10 response: %v", err)
	}
	if !obj.ActionEmailReportIncludeSplunkLogo {
		t.Errorf("ActionEmailReportIncludeSplunkLogo: got false, want true")
	}
	if obj.ActionPopulateLookupMaxTime != 300 {
		t.Errorf("ActionPopulateLookupMaxTime: got %d, want 300", obj.ActionPopulateLookupMaxTime)
	}
	if obj.ActionRSSMaxTime != 60 {
		t.Errorf("ActionRSSMaxTime: got %d, want 60", obj.ActionRSSMaxTime)
	}
	if obj.ActionScriptMaxTime != 300 {
		t.Errorf("ActionScriptMaxTime: got %d, want 300", obj.ActionScriptMaxTime)
	}
	if obj.ActionSummaryIndexMaxTime != 300 {
		t.Errorf("ActionSummaryIndexMaxTime: got %d, want 300", obj.ActionSummaryIndexMaxTime)
	}
	if obj.AutoSummarizeMaxTime != 3600 {
		t.Errorf("AutoSummarizeMaxTime: got %d, want 3600", obj.AutoSummarizeMaxTime)
	}
	if obj.DispatchIndexedRealtimeOffset != 0 {
		t.Errorf("DispatchIndexedRealtimeOffset: got %d, want 0", obj.DispatchIndexedRealtimeOffset)
	}
	if obj.DispatchIndexedRealtimeMinspan != 0 {
		t.Errorf("DispatchIndexedRealtimeMinspan: got %d, want 0", obj.DispatchIndexedRealtimeMinspan)
	}
	if obj.Search != "index=_internal | head 1" {
		t.Errorf("Search: got %q, want %q", obj.Search, "index=_internal | head 1")
	}
}

func TestFlexDuration_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		input   string
		want    FlexDuration
		wantErr bool
	}{
		{`"5m"`, 300, false},
		{`"1m"`, 60, false},
		{`"30s"`, 30, false},
		{`"2h"`, 7200, false},
		{`"1d"`, 86400, false},
		{`"3600"`, 3600, false},
		{`""`, 0, false},
		{`300`, 300, false},
		{`"5x"`, 0, true},
		{`null`, 7, false},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			f := FlexDuration(7)
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

func TestSavedSearchObject_ReportIncludeSplunkLogoBool(t *testing.T) {
	for _, input := range []string{`false`, `"0"`, `"false"`} {
		var obj SavedSearchObject
		payload := `{"action.email.reportIncludeSplunkLogo": ` + input + `}`
		if err := json.Unmarshal([]byte(payload), &obj); err != nil {
			t.Fatalf("Unmarshal %s: %v", input, err)
		}
		if obj.ActionEmailReportIncludeSplunkLogo {
			t.Errorf("Unmarshal %s: got true, want false", input)
		}
	}
}
