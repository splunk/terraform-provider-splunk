package models

import (
	"testing"

	"github.com/google/go-querystring/query"
)

func TestHttpEventCollectorObject_EnterpriseFormEncoding(t *testing.T) {
	for _, ack := range []FlexInt{0, 1} {
		obj := &HttpEventCollectorObject{
			Index:  "main",
			UseACK: ack,
		}
		values, err := query.Values(obj)
		if err != nil {
			t.Fatalf("query.Values(ack=%d): %v", ack, err)
		}
		// UpdateHttpEventCollectorObject passes a pointer to the pointer.
		updateValues, err := query.Values(&obj)
		if err != nil {
			t.Fatalf("query.Values(&obj) ack=%d: %v", ack, err)
		}
		got := values.Get("useACK")
		if got != updateValues.Get("useACK") {
			t.Fatalf("create useACK=%q, update useACK=%q", got, updateValues.Get("useACK"))
		}
		want := "0"
		if ack == 1 {
			want = "1"
		}
		if got != want {
			t.Errorf("useACK form value = %q, want %q", got, want)
		}
	}
}
