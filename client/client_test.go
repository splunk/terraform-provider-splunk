package client

import (
	"bytes"
	"io"
	"log"
	"math"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/splunk/terraform-provider-splunk/client/models"
)

const (
	testSessionKey = "123"
	testURL        = "https://test:8089/test"
)

func TestBuildSplunkURLNoURLPath(t *testing.T) {
	client, err := NewDefaultSplunkdClient()
	if err != nil {
		t.Error(err)
	}
	url := client.BuildSplunkURL(nil, "")

	if got, want := url.Hostname(), "localhost"; got != want {
		t.Errorf("hostname invalid, got %s, want %s", got, want)
	}
	if got, want := url.Scheme, defaultScheme; got != want {
		t.Errorf("scheme invalid, got %s, want %s", got, want)
	}
	if got, want := url.Port(), "8089"; got != want {
		t.Errorf("port invalid, got %s, want %s", got, want)
	}
	if got, want := url.Path, ""; got != want {
		t.Errorf("path invalid, got %s, want %s", got, want)
	}
	if got, want := url.Fragment, ""; got != want {
		t.Errorf("fragment invalid, got %s, want %s", got, want)
	}
	if url.User != nil {
		t.Errorf("user invalid, got %s, want %v", url.User, nil)
	}
}

func TestBuildSplunkURLSpecialCharactersInSearch(t *testing.T) {
	client, err := NewDefaultSplunkdClient()
	if err != nil {
		t.Error(err)
	}
	url := client.BuildSplunkURL(nil, "servicesNS", "admin", "search", "saved", "searches", "[some search]")

	if got, want := url.Hostname(), "localhost"; got != want {
		t.Errorf("hostname invalid, got %s, want %s", got, want)
	}
	if got, want := url.Scheme, defaultScheme; got != want {
		t.Errorf("scheme invalid, got %s, want %s", got, want)
	}
	if got, want := url.Port(), "8089"; got != want {
		t.Errorf("port invalid, got %s, want %s", got, want)
	}
	if got, want := url.Path, "servicesNS/admin/search/saved/searches/[some+search]"; got != want {
		t.Errorf("path invalid, got %s, want %s", got, want)
	}
	if got, want := url.Fragment, ""; got != want {
		t.Errorf("fragment invalid, got %s, want %s", got, want)
	}
	if url.User != nil {
		t.Errorf("user invalid, got %s, want %v", url.User, nil)
	}
}

func TestBuildSplunkURLWithEscapedPathPart(t *testing.T) {
	client, err := NewDefaultSplunkdClient()
	if err != nil {
		t.Error(err)
	}
	url := client.BuildSplunkURLWithEscapedPathPart(nil, "A/B+C D", "servicesNS", "nobody", "app", "saved", "searches")

	if got, want := url.Path, "servicesNS/nobody/app/saved/searches/A/B+C D"; got != want {
		t.Errorf("path invalid, got %s, want %s", got, want)
	}
	if got, want := url.EscapedPath(), "servicesNS/nobody/app/saved/searches/A%2FB%2BC%20D"; got != want {
		t.Errorf("escaped path invalid, got %s, want %s", got, want)
	}
}

func TestNamespacedResourceReadersEscapeFinalPathSegment(t *testing.T) {
	os.Setenv(envVarHTTPScheme, "http")
	defer os.Unsetenv(envVarHTTPScheme)

	tests := []struct {
		name string
		want string
		call func(*Client) (*http.Response, error)
	}{
		{
			name: "saved search",
			want: "/servicesNS/nobody/app/saved/searches/A%2FB%2BC%20D",
			call: func(client *Client) (*http.Response, error) {
				return client.ReadSavedSearches("A/B+C D", "nobody", "app")
			},
		},
		{
			name: "dashboard view",
			want: "/servicesNS/nobody/app/data/ui/views/A%2FB%2BC%20D",
			call: func(client *Client) (*http.Response, error) {
				return client.ReadDashboardObject("A/B+C D", "nobody", "app")
			},
		},
		{
			name: "federated provider",
			want: "/services/data/federated/provider/A%2FB%2BC%20D",
			call: func(client *Client) (*http.Response, error) {
				return client.ReadFederatedProvider("A/B+C D")
			},
		},
		{
			name: "configs conf stanza with slashes",
			want: "/servicesNS/nobody/system/configs/conf-inputs/monitor:%2F%2F%2Fdata%2Fsyslog%2F",
			call: func(client *Client) (*http.Response, error) {
				return client.ReadConfigsConfObject("inputs/monitor:///data/syslog/", "nobody", "system")
			},
		},
		{
			name: "configs conf delete stanza with slashes",
			want: "/servicesNS/nobody/system/configs/conf-props/source::%2Fdata%2Fjson%2F...",
			call: func(client *Client) (*http.Response, error) {
				return client.DeleteConfigsConfObject("props/source::/data/json/...", "nobody", "system")
			},
		},
		{
			name: "federated index",
			want: "/services/data/federated/index/federated:remote-main",
			call: func(client *Client) (*http.Response, error) {
				return client.ReadFederatedIndex("federated:remote-main")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got string
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				got = r.URL.EscapedPath()
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`{}`))
			}))
			defer server.Close()

			backend, err := url.Parse(server.URL)
			if err != nil {
				t.Fatal(err)
			}

			client, err := NewSplunkdClient("", defaultAuth, backend.Host, "", server.Client())
			if err != nil {
				t.Fatal(err)
			}

			resp, err := tt.call(client)
			if err != nil {
				t.Fatalf("reader call failed: %s", err)
			}
			defer resp.Body.Close()

			if got != tt.want {
				t.Errorf("escaped path invalid, got %s, want %s", got, tt.want)
			}
		})
	}
}

func TestBuildSplunkURLNoHost(t *testing.T) {
	client, err := NewDefaultSplunkdClient()
	if err != nil {
		t.Error(err)
	}
	url := client.BuildSplunkURL(nil, "services",
		"search", "jobs")

	if got, want := url.Hostname(), "localhost"; got != want {
		t.Errorf("hostname invalid, got %s, want %s", got, want)
	}
	if got, want := url.Scheme, defaultScheme; got != want {
		t.Errorf("scheme invalid, got %s, want %s", got, want)
	}
	if got, want := url.Port(), "8089"; got != want {
		t.Errorf("port invalid, got %s, want %s", got, want)
	}
	if got, want := url.Path, "services/search/jobs"; got != want {
		t.Errorf("path invalid, got %s, want %s", got, want)
	}
	if got, want := url.Fragment, ""; got != want {
		t.Errorf("fragment invalid, got %s, want %s", got, want)
	}
	if url.User != nil {
		t.Errorf("user invalid, got %s, want %v", url.User, nil)
	}
}

func TestBuildSplunkURLHTTPScheme(t *testing.T) {
	os.Setenv("HTTPScheme", "http")
	client, err := NewDefaultSplunkdClient()
	if err != nil {
		t.Error(err)
	}
	url := client.BuildSplunkURL(nil, "")
	os.Unsetenv("HTTPScheme")
	if got, want := url.Hostname(), "localhost"; got != want {
		t.Errorf("hostname invalid, got %s, want %s", got, want)
	}
	if got, want := url.Scheme, "http"; got != want {
		t.Errorf("scheme invalid, got %s, want %s", got, want)
	}
	if got, want := url.Port(), "8089"; got != want {
		t.Errorf("port invalid, got %s, want %s", got, want)
	}
	if got, want := url.Path, ""; got != want {
		t.Errorf("path invalid, got %s, want %s", got, want)
	}
	if got, want := url.Fragment, ""; got != want {
		t.Errorf("fragment invalid, got %s, want %s", got, want)
	}
	if url.User != nil {
		t.Errorf("user invalid, got %s, want %v", url.User, nil)
	}
}

func TestBuildSplunkURLCount(t *testing.T) {
	client, err := NewDefaultSplunkdClient()
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name   string
		values url.Values
		count  string
	}{
		{name: "nil values", values: nil, count: "-1"},
		{name: "missing count", values: url.Values{"datatype": {"all"}}, count: "-1"},
		{name: "explicit count", values: url.Values{"count": {"10"}}, count: "10"},
		{name: "count zero", values: url.Values{"count": {"0"}}, count: "0"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := client.BuildSplunkURL(tt.values, "services", "data", "indexes")
			q, err := url.ParseQuery(got.RawQuery)
			if err != nil {
				t.Fatalf("parse query: %v", err)
			}
			if q.Get("output_mode") != "json" {
				t.Errorf("output_mode = %q, want json", q.Get("output_mode"))
			}
			if q.Get("count") != tt.count {
				t.Errorf("count = %q, want %q", q.Get("count"), tt.count)
			}
			if tt.name == "missing count" && q.Get("datatype") != "all" {
				t.Errorf("datatype = %q, want all", q.Get("datatype"))
			}
		})
	}
}

func TestNewSplunkdHTTPClient(t *testing.T) {
	timeout := time.Second * 10
	skipValidateTLS := true
	testHTTPClient, err := NewSplunkdHTTPClient(timeout, skipValidateTLS)
	if err != nil {
		t.Error(err)
	}
	if got, want := testHTTPClient.Timeout, timeout; got != want {
		t.Errorf("NewDefaultSplunkdClient httpClient is %v, want %v", got, want)
	}
	if got, want := testHTTPClient.Transport.(*http.Transport).TLSClientConfig.InsecureSkipVerify, skipValidateTLS; got != want {
		t.Errorf("NewDefaultSplunkdClient httpClient Transport is %v, want %v", got, want)
	}
}

func TestNewRequest(t *testing.T) {
	client, err := NewDefaultSplunkdClient()
	if err != nil {
		t.Error(err)
	}
	body := []byte(`{"test":"This is a test body"}`)
	expectedBasicAuth := []string{"Basic YWRtaW46Y2hhbmdlbWU="}
	expectedUserAgent := []string{"splunk-simple-go-client"}
	requestBody := bytes.NewBuffer(body)
	tests := []struct {
		method string
		url    string
		body   io.Reader
	}{
		{MethodGet, testURL, nil},
		{MethodPost, testURL, requestBody},
		{MethodPut, testURL, requestBody},
		{MethodPatch, testURL, requestBody},
		{MethodDelete, testURL, nil},
	}
	for _, test := range tests {
		req, err := client.NewRequest(test.method, test.url, test.body)
		if err != nil {
			t.Fatalf("client.NewRequest returns unexpected error: %v", err)
		}
		if got, want := req.Method, test.method; got != want {
			t.Errorf("NewRequest http method is %v, want %v", got, want)
		}
		if got, want := req.URL.String(), test.url; got != want {
			t.Errorf("NewRequest url is %v, want %v", got, want)
		}
		if got, want := req.Header["Authorization"], expectedBasicAuth; !reflect.DeepEqual(got, want) {
			t.Errorf("NewRequest authorization is %v, want %v", got, want)
		}
		if got, want := req.Header["User-Agent"], expectedUserAgent; !reflect.DeepEqual(got, want) {
			t.Errorf("NewRequest user agent is %v, want %v", got, want)
		}
		if test.method == MethodGet || test.method == MethodDelete {
			t.Skipf("Skip NewRequest body test for %v and %v method", MethodGet, MethodDelete)
		} else {
			gotBody, _ := io.ReadAll(req.Body)
			if bytes.Compare(gotBody, body) != -1 {
				t.Errorf("NewRequest url is %v, want %v", gotBody, body)
			}
		}
	}
}

func TestNewRequestBasicAuthHeader(t *testing.T) {
	client, err := NewDefaultSplunkdClient()
	if err != nil {
		t.Error(err)
	}
	req, err := client.NewRequest(MethodGet, testURL, nil)
	if err != nil {
		t.Errorf("NewRequest returns unexpected error %v", err)
	}
	expectedBasicAuth := []string{"Basic YWRtaW46Y2hhbmdlbWU="}
	if got, want := req.Header["Authorization"], expectedBasicAuth; !reflect.DeepEqual(got, want) {
		t.Errorf("NewRequest authorization is %v, want %v", got, want)
	}
}

func TestNewRequestUserAgentHeader(t *testing.T) {
	client, err := NewDefaultSplunkdClient()
	if err != nil {
		t.Error(err)
	}
	req, err := client.NewRequest(MethodGet, testURL, nil)
	if err != nil {
		t.Errorf("NewRequest returns unexpected error %v", err)
	}
	expectedUserAgent := []string{"splunk-simple-go-client"}
	if got, want := req.Header["User-Agent"], expectedUserAgent; !reflect.DeepEqual(got, want) {
		t.Errorf("NewRequest authorization is %v, want %v", got, want)
	}
}

func TestNewRequestWithoutContentTypeHeader(t *testing.T) {
	client, err := NewDefaultSplunkdClient()
	if err != nil {
		t.Error(err)
	}
	req, err := client.NewRequest(MethodGet, testURL, nil)
	if err != nil {
		t.Errorf("NewRequest returns unexpected error %v", err)
	}
	if req.Header["Content-Type"] != nil {
		t.Errorf("NewRequest Content-Type is %v, want nil", req.Header["Content-Type"])
	}
}

func TestNewRequestWithContentTypeHeader(t *testing.T) {
	client, err := NewDefaultSplunkdClient()
	if err != nil {
		t.Error(err)
	}
	client.urlEncoded = true
	req, err := client.NewRequest(MethodGet, testURL, nil)
	if err != nil {
		t.Errorf("NewRequest returns unexpected error %v", err)
	}
	expectedContentType := []string{"application/x-www-form-urlencoded"}
	if got, want := req.Header["Content-Type"], expectedContentType; !reflect.DeepEqual(got, want) {
		t.Errorf("NewRequest Content-Type is %v, want %v", got, want)
	}
}

func TestNewRequestSessionKey(t *testing.T) {
	client, err := NewDefaultSplunkdClient()
	if err != nil {
		t.Error(err)
	}
	client.sessionKey = testSessionKey
	req, err := client.NewRequest(MethodGet, testURL, nil)
	if err != nil {
		t.Errorf("NewRequest returns unexpected error %v", err)
	}
	expectedBasicAuth := []string{"Splunk " + client.sessionKey}
	if got, want := req.Header["Authorization"], expectedBasicAuth; !reflect.DeepEqual(got, want) {
		t.Errorf("NewRequest authorization is %v, want %v", got, want)
	}
}

func TestNewRequestWithAuthToken(t *testing.T) {
	client := &Client{}
	client.authToken = "auth_token"
	req, err := client.NewRequest(MethodGet, testURL, nil)
	if err != nil {
		t.Errorf("NewRequest returns unexpected error %v", err)
	}
	expectedBasicAuth := []string{"Bearer " + client.authToken}
	if got, want := req.Header["Authorization"], expectedBasicAuth; !reflect.DeepEqual(got, want) {
		t.Errorf("NewRequest authorization is %v, want %v", got, want)
	}
}

func TestNewRequestError(t *testing.T) {
	client, err := NewDefaultSplunkdClient()
	if err != nil {
		t.Error(err)
	}
	client.sessionKey = testSessionKey
	_, err = client.NewRequest("#~/", testURL, nil)
	if err == nil {
		t.Errorf("NewRequest expected to return error, got %v", err)
	}
}

func TestEncodeRequestBodyNil(t *testing.T) {
	client, err := NewDefaultSplunkdClient()
	if err != nil {
		t.Error(err)
	}
	response, err := client.EncodeRequestBody(nil)
	if len(response) > 0 {
		t.Errorf("EncodeRequestBody expected to return nil, got %v", response)
	}
	if err != nil {
		t.Errorf("EncodeRequestBody expected to not return error, got %v", err)
	}
}

func TestEncodeRequestBodyString(t *testing.T) {
	client, err := NewDefaultSplunkdClient()
	if err != nil {
		t.Error(err)
	}
	got, err := client.EncodeRequestBody(`{"test":"This is a test body"}`)
	// expect := []byte(`{"test":"This is a test body"}`)
	if value := reflect.ValueOf(got); value.Kind() != reflect.Slice {
		t.Errorf("EncodeRequestBody expected to return []byte, got %v", got)
	}
	if err != nil {
		t.Errorf("EncodeRequestBody expected to not return error, got %v", err)
	}
}

func TestTestEncodeRequestBodyMap(t *testing.T) {
	client, err := NewDefaultSplunkdClient()
	if err != nil {
		t.Error(err)
	}
	testData := map[string]string{
		"testKey": "testValue",
	}
	got, err := client.EncodeRequestBody(testData)
	if value := reflect.ValueOf(got); value.Kind() != reflect.Slice {
		t.Errorf("EncodeRequestBody expected to return []byte, got %v", got)
	}
	if err != nil {
		t.Errorf("EncodeRequestBody expected to not return error, got %v", err)
	}
}

func TestTestEncodeRequestBodyStruct(t *testing.T) {
	client, err := NewDefaultSplunkdClient()
	if err != nil {
		t.Error(err)
	}
	type TestModel struct {
		testID    string
		testValue string
	}
	testData := TestModel{
		testID:    "123",
		testValue: "test",
	}
	got, err := client.EncodeRequestBody(testData)
	if value := reflect.ValueOf(got); value.Kind() != reflect.Slice {
		t.Errorf("EncodeRequestBody expected to return []byte, got %v", got)
	}
	if err != nil {
		t.Errorf("EncodeRequestBody expected to not return error, got %v", err)
	}
}

func TestTestEncodeRequestBodyInvalid(t *testing.T) {
	client, err := NewDefaultSplunkdClient()
	if err != nil {
		t.Error(err)
	}
	_, err = client.EncodeRequestBody(123)
	if err == nil {
		t.Errorf("EncodeRequestBody expected to raise an error, got %v", err)
	}
}

func TestEncodeObjectError(t *testing.T) {
	client, err := NewDefaultSplunkdClient()
	if err != nil {
		t.Error(err)
	}
	_, err = client.EncodeObject(math.Inf(1))
	if err == nil {
		t.Errorf("EncodeObject expected to raise an error, got %v", err)
	}
}

func TestEncodeObjectTypeConversion(t *testing.T) {
	client, err := NewDefaultSplunkdClient()
	if err != nil {
		t.Error(err)
	}
	intVal := 1
	var float32Val float32 = 0.999
	testData := map[string]interface{}{
		"testBool":    true,
		"testInt":     intVal,
		"testFloat32": float32Val,
		"testFloat64": 0.555,
	}
	want := "testBool=true&testFloat32=0.999&testFloat64=0.555&testInt=1"
	got, err := client.EncodeObject(testData)
	gotString := string(got[:])
	if gotString != want {
		t.Errorf("EncodeObject expected to return %v, got %v", want, gotString)
	}
	if err != nil {
		t.Errorf("EncodeObject expected to not return error, got %v", err)
	}
}

func TestSendsCookies(t *testing.T) {
	client, err := NewDefaultSplunkdClient()
	if err != nil {
		t.Error(err)
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if cookie, err := r.Cookie("SESSION"); err != nil {
			http.SetCookie(w, &http.Cookie{Name: "SESSION", Value: "FIRST"})
		} else {
			cookie.Value = "SECOND"
			http.SetCookie(w, cookie)
		}
	}))
	defer ts.Close()

	u, err := url.Parse(ts.URL)
	if err != nil {
		log.Fatal(err)
	}

	_, err = client.Get(*u)
	if err != nil {
		t.Error(err)
	}

	resp, err := client.Get(*u)
	if err != nil {
		t.Error(err)
	}

	if got, want := resp.Cookies()[0].Value, "SECOND"; got != want {
		t.Errorf("Returned Cookie is %v, want %v", got, want)
	}
}

func TestCreateResponsesLeaveBodyOpen(t *testing.T) {
	t.Setenv(envVarHTTPScheme, "http")

	const body = `{"entry":[]}`
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		_, _ = io.WriteString(w, body)
	}))
	defer server.Close()

	backend, err := url.Parse(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	c, err := NewSplunkdClient("", defaultAuth, backend.Host, "", server.Client())
	if err != nil {
		t.Fatal(err)
	}

	dashboard, err := c.CreateDashboardObject("admin", "search", &models.SplunkDashboardsObject{Name: "view"})
	if err != nil {
		t.Fatal(err)
	}
	defer dashboard.Body.Close()
	got, err := io.ReadAll(dashboard.Body)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != body {
		t.Fatalf("dashboard body = %q, want %q", got, body)
	}

	savedSearch, err := c.CreateSavedSearches("alert", "admin", "search", &models.SavedSearchObject{})
	if err != nil {
		t.Fatal(err)
	}
	defer savedSearch.Body.Close()
	got, err = io.ReadAll(savedSearch.Body)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != body {
		t.Fatalf("saved search body = %q, want %q", got, body)
	}
}

func TestEscapeSplunkPathPart(t *testing.T) {
	// Single encoding used by saved searches, dashboards, conf stanzas,
	// federated objects, and SAML/ProxySSO groups. Space stays %20, "/" stays
	// %2F, and a literal "+" is %2B. A bare "+" is not a space.
	tests := []struct {
		in   string
		want string
	}{
		{in: "plain", want: "plain"},
		{in: "Test 30 days", want: "Test%2030%20days"},
		{in: "Test 30+ days", want: "Test%2030%2B%20days"},
		{in: "A/B+C D", want: "A%2FB%2BC%20D"},
		{in: "Example+Saved+Search", want: "Example%2BSaved%2BSearch"},
		{in: "monitor:///var/log/", want: "monitor:%2F%2F%2Fvar%2Flog%2F"},
		{in: "100%+", want: "100%25%2B"},
		{in: "+", want: "%2B"},
		{in: "/var/log/syslog", want: "%2Fvar%2Flog%2Fsyslog"},
		{in: "opt/splunk/var/log/splunk/health.log", want: "opt%2Fsplunk%2Fvar%2Flog%2Fsplunk%2Fhealth.log"},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			if got := EscapeSplunkPathPart(tt.in); got != tt.want {
				t.Errorf("EscapeSplunkPathPart(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

// splunkd decodes "+" in a path segment as a space. Names that already worked
// (spaces as %20, slashes as %2F, no extra characters) must decode back to the
// same name, and a literal "+" must decode back to "+".
func TestSavedSearchNamesRoundTripThroughSplunkdDecoding(t *testing.T) {
	t.Setenv(envVarHTTPScheme, "http")

	names := []string{
		"plain",
		"Test 30 days",
		"Test 30+ days",
		"A/B+C D",
		"Example+Saved+Search",
		"100%+",
	}
	const prefix = "/servicesNS/nobody/search/saved/searches/"

	for _, name := range names {
		t.Run(name, func(t *testing.T) {
			calls := []struct {
				op   string
				call func(*Client) error
			}{
				{
					op: "read",
					call: func(c *Client) error {
						resp, err := c.ReadSavedSearches(name, "nobody", "search")
						if resp != nil {
							defer resp.Body.Close()
						}
						return err
					},
				},
				{
					op: "update",
					call: func(c *Client) error {
						return c.UpdateSavedSearches(name, "nobody", "search", &models.SavedSearchObject{}, nil)
					},
				},
				{
					op: "delete",
					call: func(c *Client) error {
						resp, err := c.DeleteSavedSearches(name, "nobody", "search")
						if resp != nil {
							defer resp.Body.Close()
						}
						return err
					},
				},
			}
			for _, call := range calls {
				t.Run(call.op, func(t *testing.T) {
					var got string
					client := splunkdPathClient(t, func(r *http.Request) {
						got = r.URL.EscapedPath()
					})
					if err := call.call(client); err != nil {
						t.Fatalf("call failed: %s", err)
					}
					if !strings.HasPrefix(got, prefix) {
						t.Fatalf("path %s, want prefix %s", got, prefix)
					}
					segment := got[len(prefix):]
					if segment != EscapeSplunkPathPart(name) {
						t.Errorf("segment %s, want %s", segment, EscapeSplunkPathPart(name))
					}
					decoded, err := url.QueryUnescape(segment)
					if err != nil {
						t.Fatalf("decode %s: %s", segment, err)
					}
					if decoded != name {
						t.Errorf("splunkd would look up %q, want %q", decoded, name)
					}
				})
			}
		})
	}
}

// Monitor, script, and Splunk TCP token inputs encode "/" and space twice,
// because the handler decodes "%2F" and "%20" after the HTTP decode. A literal
// "+" is encoded once: a second encode leaves "%2B" in the object id.
func TestLegacyInputPathsKeepDoubleEncoding(t *testing.T) {
	t.Setenv(envVarHTTPScheme, "http")

	type inputCase struct {
		name   string
		suffix string
	}
	cases := []inputCase{
		{name: "/var/log/syslog", suffix: "%252Fvar%252Flog%252Fsyslog"},
		{name: "opt/splunk/var/log/splunk/health.log", suffix: "opt%252Fsplunk%252Fvar%252Flog%252Fsplunk%252Fhealth.log"},
		{name: "Test 30 days", suffix: "Test%252030%2520days"},
		{name: "Test 30+ days", suffix: "Test%252030%2B%2520days"},
		{name: "plain", suffix: "plain"},
	}

	readers := []struct {
		kind   string
		prefix string
		call   func(*Client, string) (*http.Response, error)
	}{
		{
			kind:   "monitor",
			prefix: "/servicesNS/nobody/search/data/inputs/monitor/",
			call: func(c *Client, name string) (*http.Response, error) {
				return c.ReadMonitorInput(name, "nobody", "search")
			},
		},
		{
			kind:   "script",
			prefix: "/servicesNS/nobody/search/data/inputs/script/",
			call: func(c *Client, name string) (*http.Response, error) {
				return c.ReadScriptedInput(name, "nobody", "search")
			},
		},
		{
			kind:   "splunktcptoken",
			prefix: "/servicesNS/nobody/search/data/inputs/tcp/splunktcptoken/",
			call: func(c *Client, name string) (*http.Response, error) {
				return c.ReadSplunkTCPTokenInput(name, "nobody", "search")
			},
		},
	}

	for _, reader := range readers {
		for _, tc := range cases {
			t.Run(reader.kind+"/"+tc.name, func(t *testing.T) {
				assertLegacyInputPath(t, reader.prefix, tc.suffix, tc.name, func(c *Client) error {
					resp, err := reader.call(c, tc.name)
					if resp != nil {
						defer resp.Body.Close()
					}
					return err
				})
			})
		}
	}

	t.Run("splunktcptoken scheme prefix", func(t *testing.T) {
		const name = "splunktcptoken://a+b"
		assertLegacyInputPath(t, "/servicesNS/nobody/search/data/inputs/tcp/splunktcptoken/", "splunktcptoken:%252F%252Fa%2Bb", name, func(c *Client) error {
			resp, err := c.ReadSplunkTCPTokenInput(name, "nobody", "search")
			if resp != nil {
				defer resp.Body.Close()
			}
			return err
		})
	})
}

func TestLegacyInputACLKeepsDoubleEncoding(t *testing.T) {
	t.Setenv(envVarHTTPScheme, "http")

	acl := &models.ACLObject{Sharing: "app"}
	tests := []struct {
		name     string
		object   string
		wantPath string
		parts    []string
	}{
		{
			name:     "monitor slash path",
			object:   "/var/log/syslog",
			wantPath: "/servicesNS/nobody/search/data/inputs/monitor/%252Fvar%252Flog%252Fsyslog/acl",
			parts:    []string{"data", "inputs", "monitor"},
		},
		{
			name:     "monitor space",
			object:   "Test 30 days",
			wantPath: "/servicesNS/nobody/search/data/inputs/monitor/Test%252030%2520days/acl",
			parts:    []string{"data", "inputs", "monitor"},
		},
		{
			name:     "monitor plus",
			object:   "Test 30+ days",
			wantPath: "/servicesNS/nobody/search/data/inputs/monitor/Test%252030%2B%2520days/acl",
			parts:    []string{"data", "inputs", "monitor"},
		},
		{
			name:     "script plus",
			object:   "a+b",
			wantPath: "/servicesNS/nobody/search/data/inputs/script/a%2Bb/acl",
			parts:    []string{"data", "inputs", "script"},
		},
		{
			name:     "token id plus",
			object:   "a+b",
			wantPath: "/servicesNS/nobody/search/data/inputs/tcp/splunktcptoken/a%2Bb/acl",
			parts:    []string{"data", "inputs", "tcp", "splunktcptoken"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got string
			client := splunkdPathClient(t, func(r *http.Request) {
				got = r.URL.EscapedPath()
			})
			// Resources pass LegacyInputACLName. UpdateAcl escapes that once more.
			err := client.UpdateAcl("nobody", "search", LegacyInputACLName(tt.object), acl, tt.parts...)
			if err != nil {
				t.Fatal(err)
			}
			if got != tt.wantPath {
				t.Errorf("path %s, want %s", got, tt.wantPath)
			}
		})
	}
}

func assertLegacyInputPath(t *testing.T, prefix, suffix, name string, call func(*Client) error) {
	t.Helper()
	var got string
	client := splunkdPathClient(t, func(r *http.Request) {
		got = r.URL.EscapedPath()
	})
	if err := call(client); err != nil {
		t.Fatal(err)
	}
	want := prefix + suffix
	if got != want {
		t.Fatalf("path %s, want %s", got, want)
	}
	// HTTP decodes the segment once. The input handler then decodes %2F and %20
	// and leaves a "+" that was sent as %2B.
	once, err := url.PathUnescape(suffix)
	if err != nil {
		t.Fatal(err)
	}
	gotName := strings.NewReplacer("%2F", "/", "%2f", "/", "%20", " ").Replace(once)
	if gotName != name {
		t.Errorf("object id %q, want %q", gotName, name)
	}
}

func splunkdPathClient(t *testing.T, capture func(*http.Request)) *Client {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capture(r)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{}`))
	}))
	t.Cleanup(server.Close)

	backend, err := url.Parse(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	client, err := NewSplunkdClient("", defaultAuth, backend.Host, "", server.Client())
	if err != nil {
		t.Fatal(err)
	}
	return client
}
