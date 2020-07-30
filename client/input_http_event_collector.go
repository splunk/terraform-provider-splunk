package client

import (
	"github.com/google/go-querystring/query"
	"net/http"
	"terraform-provider-splunk/client/models"
)

func (client *Client) CreateHttpEventCollectorObject(name string, owner string, app string, httpInputConfigObj *models.HttpEventCollectorObject) (*http.Response, error) {
	values, err := query.Values(httpInputConfigObj)
	values.Add("name", name)
	if err != nil {
		return nil, err
	}

	endpoint := client.BuildSplunkURL(nil, "servicesNS", owner, app, "data", "inputs", "http", name)
	resp, err := client.Post(endpoint, values)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (client *Client) ReadHttpEventCollectorObject(name, owner, app string) (*http.Response, error) {
	endpoint := client.BuildSplunkURL(nil, "servicesNS", owner, app, "data", "inputs", "http", name)
	resp, err := client.Get(endpoint)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (client *Client) UpdateHttpEventCollectorObject(name string, owner string, app string, httpInputConfigObj *models.HttpEventCollectorObject) (*http.Response, error) {
	values, err := query.Values(&httpInputConfigObj)
	endpoint := client.BuildSplunkURL(nil, "servicesNS", owner, app, "data", "inputs", "http", name)
	resp, err := client.Post(endpoint, values)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (client *Client) DeleteHttpEventCollectorObject(name, owner, app string) (*http.Response, error) {
	endpoint := client.BuildSplunkURL(nil, "servicesNS", owner, app, "data", "inputs", "http", name)
	resp, err := client.Delete(endpoint)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// services/data/inputs/http
func (client *Client) ReadAllHttpEventCollectorObject() (*http.Response, error) {
	endpoint := client.BuildSplunkURL(nil, "services", "data", "inputs", "http")
	resp, err := client.Get(endpoint)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
