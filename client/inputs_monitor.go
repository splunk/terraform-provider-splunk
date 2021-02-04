package client

import (
	"github.com/splunk/terraform-provider-splunk/client/models"
	"net/http"
	"net/url"

	"github.com/google/go-querystring/query"
)

func (client *Client) CreateMonitorInput(name string, owner string, app string, inputsMonitorObject *models.InputsMonitorObject) error {
	values, err := query.Values(inputsMonitorObject)
	values.Add("name", name)
	if err != nil {
		return err
	}

	endpoint := client.BuildSplunkURL(nil, "servicesNS", owner, app, "data", "inputs", "monitor")
	resp, err := client.Post(endpoint, values)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (client *Client) ReadMonitorInput(name, owner, app string) (*http.Response, error) {
	endpoint := client.BuildSplunkURL(nil, "servicesNS", owner, app, "data", "inputs", "monitor", url.PathEscape(name))
	resp, err := client.Get(endpoint)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (client *Client) UpdateMonitorInput(name string, owner string, app string, inputsMonitorObject *models.InputsMonitorObject) error {
	values, err := query.Values(&inputsMonitorObject)
	if err != nil {
		return err
	}
	endpoint := client.BuildSplunkURL(nil, "servicesNS", owner, app, "data", "inputs", "monitor", url.PathEscape(name))
	resp, err := client.Post(endpoint, values)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (client *Client) DeleteMonitorInput(name, owner, app string) (*http.Response, error) {
	endpoint := client.BuildSplunkURL(nil, "servicesNS", owner, app, "data", "inputs", "monitor", url.PathEscape(name))
	resp, err := client.Delete(endpoint)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// services/data/inputs/monitor
func (client *Client) ReadMonitorInputs() (*http.Response, error) {
	endpoint := client.BuildSplunkURL(nil, "servicesNS", "-", "-", "data", "inputs", "monitor")
	resp, err := client.Get(endpoint)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
