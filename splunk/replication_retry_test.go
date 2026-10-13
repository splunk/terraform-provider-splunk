package splunk

import (
	"errors"
	"fmt"
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
		Status:     fmt.Sprintf("%d %s", statusCode, http.StatusText(statusCode)),
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
			if resp == nil {
				t.Fatal("resp = nil, want the last response")
			}
			resp.Body.Close()
		})
	}
}

func TestRetryOnStatusGivesUpAfterTimeout(t *testing.T) {
	resp, err := retryOnStatus(time.Second, http.StatusNotFound, func() (*http.Response, error) {
		resp := notFoundResponse()
		return resp, errors.New(resp.Status)
	})
	if err == nil || err.Error() != "404 Not Found" {
		t.Errorf("err = %v, want the last request error", err)
	}
	if resp == nil || resp.StatusCode != http.StatusNotFound {
		t.Fatalf("resp = %v, want the last 404", resp)
	}
	resp.Body.Close()
}

func TestRetryUpdateAclDoesNotRetryForbidden(t *testing.T) {
	calls := 0
	err := retryUpdateAcl(time.Second, func() error {
		calls++
		return errors.New("POST failed for endpoint /servicesNS/admin/search/saved/searches/alert/acl: 403 Forbidden: denied")
	})
	if err == nil || calls != 1 {
		t.Fatalf("err = %v, calls = %d, want one 403 and no retry", err, calls)
	}
}

func TestRetryUpdateAclRetriesNotFound(t *testing.T) {
	calls := 0
	err := retryUpdateAcl(time.Second, func() error {
		calls++
		if calls == 1 {
			return errors.New("POST failed for endpoint /acl: 404 Not Found: missing")
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if calls != 2 {
		t.Fatalf("calls = %d, want 2", calls)
	}
}

func useRefreshTimeout(t *testing.T, timeout time.Duration) {
	t.Helper()
	previous := replicationRefreshTimeout
	replicationRefreshTimeout = timeout
	t.Cleanup(func() { replicationRefreshTimeout = previous })
}
