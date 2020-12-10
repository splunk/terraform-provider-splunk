package client

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/splunk/terraform-provider-splunk/client/utils"
)

// Declare constants for service package
const (
	defaultTimeOut   = time.Second * 30
	defaultHost      = "localhost:8089"
	defaultScheme    = "https"
	MethodGet        = "GET"
	MethodPost       = "POST"
	MethodPut        = "PUT"
	MethodPatch      = "PATCH"
	MethodDelete     = "DELETE"
	envVarHTTPScheme = "HTTPScheme"
)

var defaultAuth = [2]string{"admin", "changeme"}

// A Client is used to communicate with Splunkd endpoints
type Client struct {
	authToken  string
	sessionKey string
	auth       [2]string
	host       string
	httpClient *http.Client
	userAgent  string
}

// NewRequest creates a new HTTP Request and set proper header
func (c *Client) NewRequest(httpMethod, url string, body io.Reader) (*http.Request, error) {
	request, err := http.NewRequest(httpMethod, url, body)
	if err != nil {
		return nil, err
	}
	if c.authToken != "" {
		request.Header.Add("Authorization", "Bearer "+c.authToken)
	} else if c.sessionKey != "" {
		request.Header.Add("Authorization", "Splunk "+c.sessionKey)
	} else {
		request.SetBasicAuth(c.auth[0], c.auth[1])
	}
	request.Header.Set("Accept", "application/json")
	request.Header.Set("User-Agent", c.userAgent)
	return request, nil
}

func getEnv(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		if value == "http" {
			return value
		}
	}
	return defaultValue
}

func (c *Client) BuildSplunkURL(queryValues url.Values, urlPathParts ...string) url.URL {
	buildPath := ""
	for _, pathPart := range urlPathParts {
		pathPart = strings.ReplaceAll(pathPart, " ", "+") // url parameters cannot have spaces
		buildPath = path.Join(buildPath, pathPart)
	}
	if queryValues == nil {
		queryValues = url.Values{}
	}

	queryValues.Set("output_mode", "json")
	queryValues.Set("count", "-1") // To avoid http response truncation
	httpScheme := getEnv(envVarHTTPScheme, defaultScheme)

	return url.URL{
		Scheme:   httpScheme,
		Host:     c.host,
		Path:     buildPath,
		RawQuery: queryValues.Encode(),
	}
}

// Do sends out request and returns HTTP response
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	response, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	return response, nil
}

// Get implements HTTP Get call
func (c *Client) Get(getURL url.URL) (*http.Response, error) {
	return c.DoRequest(MethodGet, getURL, nil)
}

// Post implements HTTP POST call
func (c *Client) Post(postURL url.URL, body interface{}) (*http.Response, error) {
	return c.DoRequest(MethodPost, postURL, body)
}

// Put implements HTTP PUT call
func (c *Client) Put(putURL url.URL, body interface{}) (*http.Response, error) {
	return c.DoRequest(MethodPut, putURL, body)
}

// Delete implements HTTP DELETE call
func (c *Client) Delete(deleteURL url.URL) (*http.Response, error) {
	return c.DoRequest(MethodDelete, deleteURL, nil)
}

// Patch implements HTTP Patch call
func (c *Client) Patch(patchURL url.URL, body interface{}) (*http.Response, error) {
	return c.DoRequest(MethodPatch, patchURL, body)
}

// DoRequest creates and execute a new request
func (c *Client) DoRequest(method string, requestURL url.URL, body interface{}) (*http.Response, error) {
	var buffer *bytes.Buffer
	if contentBytes, ok := body.([]byte); ok {
		buffer = bytes.NewBuffer(contentBytes)
	} else {
		if content, err := c.EncodeRequestBody(body); err == nil {
			buffer = bytes.NewBuffer(content)
		} else {
			return nil, err
		}
	}
	request, err := c.NewRequest(method, requestURL.String(), buffer)
	if err != nil {
		return nil, err
	}
	response, err := c.Do(request)
	if err != nil {
		return nil, err
	}
	return utils.ParseHTTPStatusCodeInResponse(response)
}

func (c *Client) Login() (e error) {
	loginURL := c.BuildSplunkURL(nil, "services", "auth", "login")
	bodyValues := map[string]string{
		"username": c.auth[0],
		"password": c.auth[1],
	}
	response, err := c.Post(loginURL, bodyValues)
	if response != nil {
		defer response.Body.Close()
	}
	if err != nil {
		return err
	}

	switch response.StatusCode {
	case 200:
		decoded := struct {
			SessionKey string `json:"sessionKey"`
		}{}
		_ = json.NewDecoder(response.Body).Decode(&decoded)
		c.sessionKey = decoded.SessionKey
	default:
		buf := new(bytes.Buffer)
		_, _ = buf.ReadFrom(response.Body)
		responseBody := buf.String()
		err = errors.New(responseBody)

	}
	return err
}

// EncodeRequestBody takes a json string or object and serializes it to be used in request body
func (c *Client) EncodeRequestBody(content interface{}) ([]byte, error) {
	if content != nil {
		switch value := reflect.ValueOf(content); value.Kind() {
		case reflect.String:
			return []byte(value.String()), nil
		case reflect.Map:
			return c.EncodeObject(value.Interface())
		case reflect.Struct:
			return c.EncodeObject(value.Interface())
		default:
			return nil, &utils.HTTPError{Status: 400, Message: "Bad Request"}
		}
	}
	return nil, nil
}

// EncodeObject encodes an object into url-encoded string
func (c *Client) EncodeObject(content interface{}) ([]byte, error) {
	URLValues := url.Values{}
	marshalContent, err := json.Marshal(content)
	if err != nil {
		return nil, err
	}
	var valueMap map[string]interface{}
	if err := json.Unmarshal(marshalContent, &valueMap); err != nil {
		return nil, err
	}
	for k, v := range valueMap {
		//k = strings.ToLower(k)
		switch val := v.(type) {
		case []interface{}:
			for _, ele := range val {
				if encoded, err := encodeValue(ele); err == nil && len(encoded) > 0 {
					URLValues.Add(k, encoded)
				}
			}
		case map[string]interface{}:
			for innerK, innerV := range val {
				if encoded, err := encodeValue(innerV); err == nil && len(encoded) > 0 {
					URLValues.Set(innerK, encoded)
				}
			}
		default:
			if encoded, err := encodeValue(val); err == nil && len(encoded) > 0 {
				URLValues.Set(k, encoded)
			}
		}
	}
	return []byte(URLValues.Encode()), nil
}

func encodeValue(v interface{}) (string, error) {
	switch val := v.(type) {
	case string:
		return val, nil
	case bool:
		return strconv.FormatBool(val), nil
	case int:
		return strconv.FormatInt(int64(val), 10), nil
	case float32:
		return strconv.FormatFloat(float64(val), 'f', -1, 32), nil
	case float64:
		return strconv.FormatFloat(float64(val), 'f', -1, 64), nil
	default:
		return "", fmt.Errorf("could not encode type %T", v)
	}
}

// NewDefaultSplunkdClient creates a Client with default values
func NewDefaultSplunkdClient() *Client {
	httpClient := NewSplunkdHTTPClient(defaultTimeOut, true)
	c := &Client{auth: defaultAuth, host: defaultHost, httpClient: httpClient, userAgent: "splunk-simple-go-client"}
	return c
}

// NewSplunkdClient creates a Client with custom values passed in
func NewSplunkdClient(sessionKey string, auth [2]string, host string, httpClient *http.Client) *Client {
	c := NewDefaultSplunkdClient()
	c.auth = auth
	c.host = host
	c.sessionKey = sessionKey
	if httpClient != nil {
		c.httpClient = httpClient
	}
	return c
}

// NewSplunkdClient creates a Client with custom values passed in
func NewSplunkdClientWithAuthToken(authToken string, auth [2]string, host string, httpClient *http.Client) *Client {
	c := NewDefaultSplunkdClient()
	c.auth = auth
	c.host = host
	c.authToken = authToken
	if httpClient != nil {
		c.httpClient = httpClient
	}
	return c
}

// NewSplunkdHTTPClient returns a HTTP Client with timeout and tls validation setup
func NewSplunkdHTTPClient(timeout time.Duration, skipValidateTLS bool) *http.Client {
	httpScheme := getEnv(envVarHTTPScheme, defaultScheme)
	if httpScheme == "http" {
		return &http.Client{
			Timeout: timeout,
		}
	}

	return &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: skipValidateTLS},
		},
	}

}
