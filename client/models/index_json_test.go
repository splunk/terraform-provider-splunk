package models

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/google/go-querystring/query"
)

func TestIndexObject_MixedJSON(t *testing.T) {
	tests := []struct {
		name string
		body string
		hint string
		hot  FlexHotBuckets
		rep  string
	}{
		{
			name: "numbers and numeric strings",
			body: `{
				"bucketRebuildMemoryHint": 16777216,
				"maxHotBuckets": "6",
				"repFactor": 0
			}`,
			hint: "16777216",
			hot:  6,
			rep:  "0",
		},
		{
			name: "auto and size suffix",
			body: `{
				"bucketRebuildMemoryHint": "16MB",
				"maxHotBuckets": "auto",
				"repFactor": "auto"
			}`,
			hint: "16MB",
			hot:  3,
			rep:  "auto",
		},
		{
			name: "legacy integer maxHotBuckets",
			body: `{
				"bucketRebuildMemoryHint": "auto",
				"maxHotBuckets": 3,
				"repFactor": "0"
			}`,
			hint: "auto",
			hot:  3,
			rep:  "0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var obj IndexObject
			if err := json.Unmarshal([]byte(tt.body), &obj); err != nil {
				t.Fatalf("Unmarshal: %v", err)
			}
			if got := string(obj.BucketRebuildMemoryHint); got != tt.hint {
				t.Errorf("bucketRebuildMemoryHint = %q, want %q", got, tt.hint)
			}
			if obj.MaxHotBuckets != tt.hot {
				t.Errorf("maxHotBuckets = %d, want %d", obj.MaxHotBuckets, tt.hot)
			}
			if got := string(obj.RepFactor); got != tt.rep {
				t.Errorf("repFactor = %q, want %q", got, tt.rep)
			}
		})
	}
}

func TestIndexObject_UnmarshalErrorNamesField(t *testing.T) {
	tests := []struct {
		body  string
		field string
	}{
		{`{"maxHotBuckets":{"unexpected":true}}`, "maxHotBuckets"},
		{`{"bucketRebuildMemoryHint":{"unexpected":true}}`, "bucketRebuildMemoryHint"},
		{`{"repFactor":{"unexpected":true}}`, "repFactor"},
	}
	for _, tt := range tests {
		var obj IndexObject
		err := json.Unmarshal([]byte(tt.body), &obj)
		if err == nil {
			t.Fatalf("Unmarshal(%s) error = nil, want a field error", tt.body)
		}
		if got := err.Error(); !strings.Contains(got, tt.field) {
			t.Fatalf("Unmarshal error = %q, want it to name %s", got, tt.field)
		}
	}
}

func TestIndexObject_FormEncodingKeepsStringValues(t *testing.T) {
	obj := &IndexObject{
		BucketRebuildMemoryHint: "16MB",
		MaxHotBuckets:           6,
		RepFactor:               "0",
	}
	values, err := query.Values(obj)
	if err != nil {
		t.Fatalf("query.Values: %v", err)
	}
	// UpdateIndexObject passes a pointer to the pointer.
	updateValues, err := query.Values(&obj)
	if err != nil {
		t.Fatalf("query.Values(&obj): %v", err)
	}
	for _, key := range []string{"bucketRebuildMemoryHint", "maxHotBuckets", "repFactor"} {
		if values.Get(key) != updateValues.Get(key) {
			t.Fatalf("create %s=%q, update %s=%q", key, values.Get(key), key, updateValues.Get(key))
		}
	}
	if got := values.Get("bucketRebuildMemoryHint"); got != "16MB" {
		t.Errorf("bucketRebuildMemoryHint = %q, want 16MB", got)
	}
	if got := values.Get("maxHotBuckets"); got != "6" {
		t.Errorf("maxHotBuckets = %q, want 6", got)
	}
	if got := values.Get("repFactor"); got != "0" {
		t.Errorf("repFactor = %q, want 0", got)
	}

	obj.RepFactor = "auto"
	values, err = query.Values(obj)
	if err != nil {
		t.Fatalf("query.Values auto: %v", err)
	}
	if got := values.Get("repFactor"); got != "auto" {
		t.Errorf("repFactor = %q, want auto", got)
	}
}
