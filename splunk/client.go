package splunk

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"
)

type SplunkClient struct {
	HttpClient         *http.Client
	Url                string
	Username           string
	Password           string
	InsecureSkipVerify bool
	SessionKey         string
	OutputMode         string
}

type Response struct {
	SessionKey    string `json:"sessionKey"`
}

type ClientParams func(*SplunkClient) error

func NewClient(url string, username string, password string, opts... ClientParams) (*SplunkClient, error) {
	client := &SplunkClient{
		Url:url,
		Username: username,
		Password: password,
		HttpClient: &http.Client{
			Timeout: time.Second * 30,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		},
	}

	for _, option := range opts {
		option(client)
	}

	return client, nil
}

func httpClient(httpClient *http.Client) ClientParams {
	return func(client *SplunkClient) error {
		client.HttpClient = httpClient
		return nil
	}
}

func sessionKey(sessionKey string) ClientParams {
	return func(client *SplunkClient) error {
		client.SessionKey = sessionKey
		return nil
	}
}

func (s *SplunkClient) setSessionKey() (e error) {
	values := url.Values{}
	values.Add("username", s.Username)
	values.Add("password", s.Password)
	values.Add("output_mode", "json")
	u, err := url.Parse(s.Url)
	u.Path = path.Join(u.Path, "services/auth/login")
	req, err := http.NewRequest("POST", u.String(), strings.NewReader(values.Encode()))
	req.SetBasicAuth(s.Username, s.Password)
	res, err := s.HttpClient.Do(req)
	defer res.Body.Close()
	if err != nil {
		return err
	}

	switch res.StatusCode {
	case 200:
		decoded := Response{}
		_ = json.NewDecoder(res.Body).Decode(&decoded)
		s.SessionKey = decoded.SessionKey
	default:
		buf := new(bytes.Buffer)
		_, _ = buf.ReadFrom(res.Body)
		responseBody := buf.String()
		err = errors.New(responseBody)

	}
	return err
}

func (c *SplunkClient) doRequest(method string, endpoint string, params url.Values) (*http.Response, error) {
	if c.SessionKey == "" {
		err := c.setSessionKey()
		if err != nil {
			return nil, err
		}
	}
	return c.doRequestWithSessionKey(method, endpoint, c.SessionKey, params)
}

func (c *SplunkClient) doRequestWithSessionKey(method string, endpoint string, sessionKey string, parameters url.Values) (*http.Response, error) {
	u, err := url.Parse(c.Url)
	if err != nil {
		return nil, err
	}
	u.Path = path.Join(u.Path, endpoint)
	parameters.Add("output_mode", "json")
	req, err := http.NewRequest(method, u.String(), strings.NewReader(parameters.Encode()))
	req.Header.Add("Authorization", "Splunk " + sessionKey)
	return c.HttpClient.Do(req)
}