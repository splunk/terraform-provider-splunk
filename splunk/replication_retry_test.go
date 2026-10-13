package splunk

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/splunk/terraform-provider-splunk/client"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

func newRoundTripProvider(t *testing.T, transport roundTripFunc) *SplunkProvider {
	t.Helper()
	c, err := client.NewSplunkdClient("", [2]string{}, "local:8089", "", &http.Client{Transport: transport})
	if err != nil {
		t.Fatal(err)
	}
	return &SplunkProvider{Client: c}
}

func jsonResponse(statusCode int, body string) *http.Response {
	return &http.Response{
		StatusCode: statusCode,
		Status:     http.StatusText(statusCode),
		Body:       io.NopCloser(strings.NewReader(body)),
		Header:     make(http.Header),
	}
}

func notFoundResponse() *http.Response {
	return jsonResponse(http.StatusNotFound, `{"messages":[{"type":"ERROR","text":"Could not find object"}]}`)
}

func TestRetryOnStatus(t *testing.T) {
	tests := []struct {
		name         string
		statuses     []int
		wantErr      bool
		wantRequests int
	}{
		{"succeeds first time", []int{http.StatusOK}, false, 1},
		{"retries the status until it succeeds", []int{http.StatusNotFound, http.StatusNotFound, http.StatusOK}, false, 3},
		{"does not retry other errors", []int{http.StatusForbidden, http.StatusOK}, true, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requests := 0
			resp, err := retryOnStatus(time.Minute, http.StatusNotFound, func() (*http.Response, error) {
				status := tt.statuses[min(requests, len(tt.statuses)-1)]
				requests++
				resp := jsonResponse(status, "{}")
				if status != http.StatusOK {
					return resp, errors.New(resp.Status)
				}
				return resp, nil
			})
			if (err != nil) != tt.wantErr {
				t.Fatalf("err = %v, wantErr %v", err, tt.wantErr)
			}
			if requests != tt.wantRequests {
				t.Errorf("requests = %d, want %d", requests, tt.wantRequests)
			}
			if (resp == nil) != tt.wantErr {
				t.Errorf("resp = %v, want a response only on success", resp)
			}
		})
	}
}

func TestRetryOnStatusGivesUpAfterTimeout(t *testing.T) {
	_, err := retryOnStatus(time.Second, http.StatusNotFound, func() (*http.Response, error) {
		resp := notFoundResponse()
		return resp, errors.New(resp.Status)
	})
	if err == nil || err.Error() != "Not Found" {
		t.Errorf("err = %v, want the last request error", err)
	}
}
