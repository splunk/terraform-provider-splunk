package splunk

import (
	"bytes"
	"github.com/likexian/gokit/assert"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
)

type RoundTripFunc func(req *http.Request) *http.Response

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

func NewTestClient(fn RoundTripFunc) *SplunkClient {
	client := &SplunkClient{}
	client.HttpClient = &http.Client{
		Transport: RoundTripFunc(fn),
	}
	return client
}

func TestSplunkClient(t *testing.T) {
	client := NewTestClient(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewBufferString(`OK`)),
			Header:     make(http.Header),
		}
	})

	body, err := client.doRequest("POST", "/", url.Values{})
	assert.Nil(t, err)
	bodyBytes, _ := ioutil.ReadAll(body.Body)
	assert.Equal(t, bodyBytes, []byte("OK"))
}
