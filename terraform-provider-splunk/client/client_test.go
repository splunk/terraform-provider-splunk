package client

import (
	"bytes"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"reflect"
	"testing"
	"time"
)

const (
	testSessionKey = "123"
	testURL        = "https://test:8089/test"
)

func TestBuildSplunkURLNoURLPath(t *testing.T) {
	client := NewDefaultSplunkdClient()
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
	client := NewDefaultSplunkdClient()
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

func TestBuildSplunkURLNoHost(t *testing.T) {
	client := NewDefaultSplunkdClient()
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
	client := NewDefaultSplunkdClient()
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

func TestNewSplunkdHTTPClient(t *testing.T) {
	timeout := time.Second * 10
	skipValidateTLS := true
	testHTTPClient := NewSplunkdHTTPClient(timeout, skipValidateTLS)
	if got, want := testHTTPClient.Timeout, timeout; got != want {
		t.Errorf("NewDefaultSplunkdClient httpClient is %v, want %v", got, want)
	}
	if got, want := testHTTPClient.Transport.(*http.Transport).TLSClientConfig.InsecureSkipVerify, skipValidateTLS; got != want {
		t.Errorf("NewDefaultSplunkdClient httpClient Transport is %v, want %v", got, want)
	}
}

func TestNewRequest(t *testing.T) {
	client := NewDefaultSplunkdClient()
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
			gotBody, _ := ioutil.ReadAll(req.Body)
			if bytes.Compare(gotBody, body) != -1 {
				t.Errorf("NewRequest url is %v, want %v", gotBody, body)
			}
		}
	}
}

func TestNewRequestBasicAuthHeader(t *testing.T) {
	client := NewDefaultSplunkdClient()
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
	client := NewDefaultSplunkdClient()
	req, err := client.NewRequest(MethodGet, testURL, nil)
	if err != nil {
		t.Errorf("NewRequest returns unexpected error %v", err)
	}
	expectedUserAgent := []string{"splunk-simple-go-client"}
	if got, want := req.Header["User-Agent"], expectedUserAgent; !reflect.DeepEqual(got, want) {
		t.Errorf("NewRequest authorization is %v, want %v", got, want)
	}
}

func TestNewRequestSessionKey(t *testing.T) {
	client := NewDefaultSplunkdClient()
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
	client := NewDefaultSplunkdClient()
	client.sessionKey = testSessionKey
	_, err := client.NewRequest("#~/", testURL, nil)
	if err == nil {
		t.Errorf("NewRequest expected to return error, got %v", err)
	}
}

func TestEncodeRequestBodyNil(t *testing.T) {
	client := NewDefaultSplunkdClient()
	response, err := client.EncodeRequestBody(nil)
	if len(response) > 0 {
		t.Errorf("EncodeRequestBody expected to return nil, got %v", response)
	}
	if err != nil {
		t.Errorf("EncodeRequestBody expected to not return error, got %v", err)
	}
}

func TestEncodeRequestBodyString(t *testing.T) {
	client := NewDefaultSplunkdClient()
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
	client := NewDefaultSplunkdClient()
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
	client := NewDefaultSplunkdClient()
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
	client := NewDefaultSplunkdClient()
	_, err := client.EncodeRequestBody(123)
	if err == nil {
		t.Errorf("EncodeRequestBody expected to raise an error, got %v", err)
	}
}

func TestEncodeObjectError(t *testing.T) {
	client := NewDefaultSplunkdClient()
	_, err := client.EncodeObject(math.Inf(1))
	if err == nil {
		t.Errorf("EncodeObject expected to raise an error, got %v", err)
	}
}

func TestEncodeObjectTypeConversion(t *testing.T) {
	client := NewDefaultSplunkdClient()
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
