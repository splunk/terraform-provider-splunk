package client

import (
	"net/http"
	"github.com/splunk/terraform-provider-splunk/client/models"

	"github.com/google/go-querystring/query"
)

func (client *Client) CreateTCPServerOutput(name string, owner string, app string, outputsTCPServerObject *models.OutputsTCPServerObject) error {
	values, err := query.Values(outputsTCPServerObject)
	values.Add("name", name)
	if err != nil {
		return err
	}

	endpoint := client.BuildSplunkURL(nil, "servicesNS", owner, app, "data", "outputs", "tcp", "server")
	resp, err := client.Post(endpoint, values)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (client *Client) ReadTCPServerOutput(name, owner, app string) (*http.Response, error) {
	endpoint := client.BuildSplunkURL(nil, "servicesNS", owner, app, "data", "outputs", "tcp", "server", name)
	resp, err := client.Get(endpoint)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (client *Client) UpdateTCPServerOutput(name string, owner string, app string, outputsTCPServerObject *models.OutputsTCPServerObject) error {
	values, err := query.Values(&outputsTCPServerObject)
	if err != nil {
		return err
	}
	endpoint := client.BuildSplunkURL(nil, "servicesNS", owner, app, "data", "outputs", "tcp", "server", name)
	resp, err := client.Post(endpoint, values)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (client *Client) DeleteTCPServerOutput(name, owner, app string) (*http.Response, error) {
	endpoint := client.BuildSplunkURL(nil, "servicesNS", owner, app, "data", "outputs", "tcp", "server", name)
	resp, err := client.Delete(endpoint)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// services/data/outputs/tcp/server
func (client *Client) ReadTCPServerOutputs() (*http.Response, error) {
	endpoint := client.BuildSplunkURL(nil, "servicesNS", "-", "-", "data", "outputs", "tcp", "server")
	resp, err := client.Get(endpoint)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
