package client

import (
	"net/http"

	"github.com/rsrdesarrollo/terraform-provider-splunk/client/models"

	"github.com/google/go-querystring/query"
)

func (client *Client) CreateTCPDefaultOutput(name string, owner string, app string, outputsTCPDefaultObject *models.OutputsTCPDefaultObject) error {
	values, err := query.Values(outputsTCPDefaultObject)
	//values.Add("name", name)
	if err != nil {
		return err
	}

	endpoint := client.BuildSplunkURL(nil, "servicesNS", owner, app, "data", "outputs", "tcp", "default", name)
	resp, err := client.Post(endpoint, values)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (client *Client) ReadTCPDefaultOutput(name, owner, app string) (*http.Response, error) {
	endpoint := client.BuildSplunkURL(nil, "servicesNS", owner, app, "data", "outputs", "tcp", "default", name)
	resp, err := client.Get(endpoint)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (client *Client) UpdateTCPDefaultOutput(name string, owner string, app string, outputsTCPDefaultObject *models.OutputsTCPDefaultObject) error {
	values, err := query.Values(&outputsTCPDefaultObject)
	if err != nil {
		return err
	}
	endpoint := client.BuildSplunkURL(nil, "servicesNS", owner, app, "data", "outputs", "tcp", "default", name)
	resp, err := client.Post(endpoint, values)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (client *Client) DeleteTCPDefaultOutput(name, owner, app string) (*http.Response, error) {
	endpoint := client.BuildSplunkURL(nil, "servicesNS", owner, app, "data", "outputs", "tcp", "default", name)
	resp, err := client.Delete(endpoint)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// services/data/outputs/tcp/default
func (client *Client) ReadTCPDefaultOutputs() (*http.Response, error) {
	endpoint := client.BuildSplunkURL(nil, "servicesNS", "-", "-", "data", "outputs", "tcp", "default")
	resp, err := client.Get(endpoint)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
