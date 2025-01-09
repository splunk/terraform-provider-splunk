package client

import (
	"net/http"

	"github.com/splunk/terraform-provider-splunk/client/models"

	"github.com/google/go-querystring/query"
)

// CreateServerClassObject creates a new server class
func (client *Client) CreateServerClassObject(serverClass string, serverClassObject *models.ServerClassObject) error {
	values, err := query.Values(serverClassObject)
	if err != nil {
		return err
	}
	endpoint := client.BuildSplunkURL(nil, "services", "deployment", "server", "serverclasses", serverClass)
	resp, err := client.Post(endpoint, values)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

// ReadServerClassObject reads an existing server class
func (client *Client) ReadServerClassObject(serverClass string) (*http.Response, error) {
	endpoint := client.BuildSplunkURL(nil, "services", "deployment", "server", "serverclasses", serverClass)
	resp, err := client.Get(endpoint)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// UpdateServerClassObject updates an existing server class
func (client *Client) UpdateServerClassObject(serverClass string, serverClassObject *models.ServerClassObject) error {
	values, err := query.Values(serverClassObject)
	if err != nil {
		return err
	}
	endpoint := client.BuildSplunkURL(nil, "services", "deployment", "server", "serverclasses", serverClass)
	resp, err := client.Post(endpoint, values)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

// DeleteServerClassObject deletes an existing server class
func (client *Client) DeleteServerClassObject(serverClass string) error {
	endpoint := client.BuildSplunkURL(nil, "services", "deployment", "server", "serverclasses", serverClass)
	resp, err := client.Delete(endpoint)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}
