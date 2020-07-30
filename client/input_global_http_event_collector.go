package client

import (
	"github.com/google/go-querystring/query"
	"net/http"
	"terraform-provider-splunk/client/models"
)

func (client *Client) CreateGlobalHttpEventCollectorObject(httpInputConfigObj models.GlobalHttpEventCollectorObject) (*http.Response, error) {
	values, err := query.Values(&httpInputConfigObj)
	endpoint := client.BuildSplunkURL(nil, "services", "data", "inputs", "http", "http")
	resp, err := client.Post(endpoint, values)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (client *Client) ReadGlobalHttpEventCollectorObject() (*http.Response, error) {
	endpoint := client.BuildSplunkURL(nil, "services", "data", "inputs", "http", "http")
	resp, err := client.Get(endpoint)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (client *Client) UpdateGlobalHttpEventCollectorObject(httpInputConfigObj models.GlobalHttpEventCollectorObject) (*http.Response, error) {
	values, err := query.Values(&httpInputConfigObj)
	endpoint := client.BuildSplunkURL(nil, "services", "data", "inputs", "http", "http")
	resp, err := client.Post(endpoint, values)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
