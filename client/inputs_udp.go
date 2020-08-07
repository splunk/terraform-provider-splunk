package client

import (
	"github.com/google/go-querystring/query"
	"net/http"
	"terraform-provider-splunk/client/models"
)

func (client *Client) CreateUDPInput(name string, owner string, app string, inputsUDPObject *models.InputsUDPObject) error {
	values, err := query.Values(inputsUDPObject)
	values.Add("name", name)
	if err != nil {
		return err
	}

	endpoint := client.BuildSplunkURL(nil, "servicesNS", owner, app, "data", "inputs", "udp")
	resp, err := client.Post(endpoint, values)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (client *Client) ReadUDPInput(name, owner, app string) (*http.Response, error) {
	endpoint := client.BuildSplunkURL(nil, "servicesNS", owner, app, "data", "inputs", "udp", name)
	resp, err := client.Get(endpoint)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (client *Client) UpdateUDPInput(name string, owner string, app string, inputsUDPObject *models.InputsUDPObject) error {
	values, err := query.Values(&inputsUDPObject)
	endpoint := client.BuildSplunkURL(nil, "servicesNS", owner, app, "data", "inputs", "udp", name)
	resp, err := client.Post(endpoint, values)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (client *Client) DeleteUDPInput(name, owner, app string) (*http.Response, error) {
	endpoint := client.BuildSplunkURL(nil, "servicesNS", owner, app, "data", "inputs", "udp", name)
	resp, err := client.Delete(endpoint)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// services/data/inputs/udp
func (client *Client) ReadUDPInputs() (*http.Response, error) {
	endpoint := client.BuildSplunkURL(nil, "services", "data", "inputs", "udp")
	resp, err := client.Get(endpoint)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
