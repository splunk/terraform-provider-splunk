package client

import (
	"net/http"

	"github.com/google/go-querystring/query"
	"github.com/splunk/terraform-provider-splunk/client/models"
)

// CreateFederatedIndex submits a short name; Splunk adds the federated: prefix.
func (client *Client) CreateFederatedIndex(name string, index *models.FederatedIndexObject) error {
	values, err := query.Values(index)
	if err != nil {
		return err
	}
	values.Set("name", name)

	resp, err := client.Post(client.BuildSplunkURL(nil, "services", "data", "federated", "index"), values)
	if resp != nil {
		defer func() { _ = resp.Body.Close() }()
	}
	return err
}

// ReadFederatedIndex gets a federated index by name.
func (client *Client) ReadFederatedIndex(name string) (*http.Response, error) {
	endpoint := client.BuildSplunkURLWithEscapedPathPart(nil, name, "services", "data", "federated", "index")
	return client.Get(endpoint)
}

// UpdateFederatedIndex updates a federated index definition.
func (client *Client) UpdateFederatedIndex(name string, index *models.FederatedIndexObject) error {
	values, err := query.Values(index)
	if err != nil {
		return err
	}

	endpoint := client.BuildSplunkURLWithEscapedPathPart(nil, name, "services", "data", "federated", "index")
	resp, err := client.Post(endpoint, values)
	if resp != nil {
		defer func() { _ = resp.Body.Close() }()
	}
	return err
}

// DeleteFederatedIndex removes a federated index definition.
func (client *Client) DeleteFederatedIndex(name string) (*http.Response, error) {
	endpoint := client.BuildSplunkURLWithEscapedPathPart(nil, name, "services", "data", "federated", "index")
	return client.Delete(endpoint)
}
