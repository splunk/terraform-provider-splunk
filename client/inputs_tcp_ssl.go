package client

import (
	"github.com/jaware-splunk/terraform-provider-splunk/client/models"
	"net/http"

	"github.com/google/go-querystring/query"
)

func (client *Client) CreateTCPSSLInput(inputsTCPSSLObject *models.InputsTCPSSLObject) error {
	values, err := query.Values(inputsTCPSSLObject)
	if err != nil {
		return err
	}

	endpoint := client.BuildSplunkURL(nil, "services", "data", "inputs", "tcp", "ssl", "ssl")
	resp, err := client.Post(endpoint, values)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (client *Client) ReadTCPSSLInput() (*http.Response, error) {
	endpoint := client.BuildSplunkURL(nil, "servicesNS", "-", "-", "data", "inputs", "tcp", "ssl", "ssl")
	resp, err := client.Get(endpoint)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (client *Client) UpdateTCPSSLInput(name string, inputsTCPSSLObject *models.InputsTCPSSLObject) error {
	values, err := query.Values(&inputsTCPSSLObject)
	if err != nil {
		return err
	}
	endpoint := client.BuildSplunkURL(nil, "services", "data", "inputs", "tcp", "ssl", name)
	resp, err := client.Post(endpoint, values)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (client *Client) ReadTCPSSLInputs() (*http.Response, error) {
	endpoint := client.BuildSplunkURL(nil, "services", "data", "inputs", "tcp", "ssl")
	resp, err := client.Get(endpoint)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
