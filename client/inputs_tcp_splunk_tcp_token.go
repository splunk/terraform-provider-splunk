package client

import (
	"net/http"
	"net/url"

	"github.com/rsrdesarrollo/terraform-provider-splunk/client/models"

	"github.com/google/go-querystring/query"
)

func (client *Client) CreateSplunkTCPTokenInput(owner string, app string, inputsSplunkTCPTokenObject *models.InputsSplunkTCPTokenObject) error {
	values, err := query.Values(inputsSplunkTCPTokenObject)
	if err != nil {
		return err
	}

	endpoint := client.BuildSplunkURL(nil, "servicesNS", owner, app, "data", "inputs", "tcp", "splunktcptoken")
	resp, err := client.Post(endpoint, values)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (client *Client) ReadSplunkTCPTokenInput(name, owner, app string) (*http.Response, error) {
	endpoint := client.BuildSplunkURL(nil, "servicesNS", owner, app, "data", "inputs", "tcp", "splunktcptoken", url.PathEscape(name))
	resp, err := client.Get(endpoint)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (client *Client) UpdateSplunkTCPTokenInput(name string, owner string, app string, inputsSplunkTCPTokenObject *models.InputsSplunkTCPTokenObject) error {
	values, err := query.Values(&inputsSplunkTCPTokenObject)
	if err != nil {
		return err
	}
	endpoint := client.BuildSplunkURL(nil, "servicesNS", owner, app, "data", "inputs", "tcp", "splunktcptoken", url.PathEscape(name))
	resp, err := client.Post(endpoint, values)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (client *Client) DeleteSplunkTCPTokenInput(name, owner, app string) (*http.Response, error) {
	endpoint := client.BuildSplunkURL(nil, "servicesNS", owner, app, "data", "inputs", "tcp", "splunktcptoken", url.PathEscape(name))
	resp, err := client.Delete(endpoint)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// services/data/inputs/tcp/splunktcptoken
func (client *Client) ReadSplunkTCPTokenInputs() (*http.Response, error) {
	endpoint := client.BuildSplunkURL(nil, "servicesNS", "-", "-", "data", "inputs", "tcp", "splunktcptoken")
	resp, err := client.Get(endpoint)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
