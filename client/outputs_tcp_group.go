package client

import (
	"github.com/jaware-splunk/terraform-provider-splunk/client/models"
	"net/http"
	"strings"

	"github.com/google/go-querystring/query"
)

func (client *Client) CreateTCPGroupOutput(name string, owner string, app string, outputsTCPGroupObject *models.OutputsTCPGroupObject) error {
	values, err := query.Values(outputsTCPGroupObject)
	values.Add("name", name)
	if err != nil {
		return err
	}

	values.Set("servers", strings.Join(outputsTCPGroupObject.Servers, ","))
	endpoint := client.BuildSplunkURL(nil, "servicesNS", owner, app, "data", "outputs", "tcp", "group")
	resp, err := client.Post(endpoint, values)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (client *Client) ReadTCPGroupOutput(name, owner, app string) (*http.Response, error) {
	endpoint := client.BuildSplunkURL(nil, "servicesNS", owner, app, "data", "outputs", "tcp", "group", name)
	resp, err := client.Get(endpoint)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (client *Client) UpdateTCPGroupOutput(name string, owner string, app string, outputsTCPGroupObject *models.OutputsTCPGroupObject) error {
	values, err := query.Values(&outputsTCPGroupObject)
	if err != nil {
		return err
	}
	values.Set("servers", strings.Join(outputsTCPGroupObject.Servers, ","))
	endpoint := client.BuildSplunkURL(nil, "servicesNS", owner, app, "data", "outputs", "tcp", "group", name)
	resp, err := client.Post(endpoint, values)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (client *Client) DeleteTCPGroupOutput(name, owner, app string) (*http.Response, error) {
	endpoint := client.BuildSplunkURL(nil, "servicesNS", owner, app, "data", "outputs", "tcp", "group", name)
	resp, err := client.Delete(endpoint)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// services/data/outputs/tcp/group
func (client *Client) ReadTCPGroupOutputs() (*http.Response, error) {
	endpoint := client.BuildSplunkURL(nil, "servicesNS", "-", "-", "data", "outputs", "tcp", "group")
	resp, err := client.Get(endpoint)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
