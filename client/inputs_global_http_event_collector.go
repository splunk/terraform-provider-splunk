package client

import (
	"net/http"
	"github.com/terraform-providers/terraform-provider-splunk/client/models"

	"github.com/google/go-querystring/query"
)

func (client *Client) CreateGlobalHttpEventCollectorObject(httpInputConfigObj models.GlobalHttpEventCollectorObject) error {
	values, err := query.Values(&httpInputConfigObj)
	if err != nil {
		return err
	}
	endpoint := client.BuildSplunkURL(nil, "services", "data", "inputs", "http", "http")
	resp, err := client.Post(endpoint, values)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (client *Client) ReadGlobalHttpEventCollectorObject() (*http.Response, error) {
	endpoint := client.BuildSplunkURL(nil, "services", "data", "inputs", "http", "http")
	resp, err := client.Get(endpoint)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (client *Client) UpdateGlobalHttpEventCollectorObject(httpInputConfigObj models.GlobalHttpEventCollectorObject) error {
	values, err := query.Values(&httpInputConfigObj)
	if err != nil {
		return err
	}
	endpoint := client.BuildSplunkURL(nil, "services", "data", "inputs", "http", "http")
	resp, err := client.Post(endpoint, values)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
