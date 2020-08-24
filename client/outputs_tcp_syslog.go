package client

import (
	"github.com/google/go-querystring/query"
	"net/http"
	"terraform-provider-splunk/client/models"
)

func (client *Client) CreateTCPSyslogOutput(name string, owner string, app string, outputsTCPSyslogObject *models.OutputsTCPSyslogObject) error {
	values, err := query.Values(outputsTCPSyslogObject)
	values.Add("name", name)
	if err != nil {
		return err
	}

	endpoint := client.BuildSplunkURL(nil, "servicesNS", owner, app, "data", "outputs", "tcp", "syslog")
	resp, err := client.Post(endpoint, values)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (client *Client) ReadTCPSyslogOutput(name, owner, app string) (*http.Response, error) {
	endpoint := client.BuildSplunkURL(nil, "servicesNS", owner, app, "data", "outputs", "tcp", "syslog", name)
	resp, err := client.Get(endpoint)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (client *Client) UpdateTCPSyslogOutput(name string, owner string, app string, outputsTCPSyslogObject *models.OutputsTCPSyslogObject) error {
	values, err := query.Values(&outputsTCPSyslogObject)
	endpoint := client.BuildSplunkURL(nil, "servicesNS", owner, app, "data", "outputs", "tcp", "syslog", name)
	resp, err := client.Post(endpoint, values)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (client *Client) DeleteTCPSyslogOutput(name, owner, app string) (*http.Response, error) {
	endpoint := client.BuildSplunkURL(nil, "servicesNS", owner, app, "data", "outputs", "tcp", "syslog", name)
	resp, err := client.Delete(endpoint)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// services/data/outputs/tcp/syslog
func (client *Client) ReadTCPSyslogOutputs() (*http.Response, error) {
	endpoint := client.BuildSplunkURL(nil, "servicesNS", "-", "-", "data", "outputs", "tcp", "syslog")
	resp, err := client.Get(endpoint)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
