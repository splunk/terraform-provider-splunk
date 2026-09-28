package client

import (
	"net/http"

	"github.com/google/go-querystring/query"
	"github.com/splunk/terraform-provider-splunk/client/models"
)

// CreateFederatedProvider creates a Splunk-to-Splunk federated provider.
func (client *Client) CreateFederatedProvider(name string, provider *models.FederatedProviderObject) error {
	values, err := query.Values(provider)
	if err != nil {
		return err
	}
	values.Set("name", name)

	resp, err := client.Post(client.BuildSplunkURL(nil, "services", "data", "federated", "provider"), values)
	if resp != nil {
		defer func() { _ = resp.Body.Close() }()
	}
	return err
}

// ReadFederatedProvider gets a provider by name. The response is returned on
// HTTP errors so callers can distinguish a missing object from other failures.
func (client *Client) ReadFederatedProvider(name string) (*http.Response, error) {
	endpoint := client.BuildSplunkURLWithEscapedPathPart(nil, name, "services", "data", "federated", "provider")
	return client.Get(endpoint)
}

// UpdateFederatedProvider updates a federated provider definition.
func (client *Client) UpdateFederatedProvider(name string, provider *models.FederatedProviderObject) error {
	values, err := query.Values(provider)
	if err != nil {
		return err
	}

	endpoint := client.BuildSplunkURLWithEscapedPathPart(nil, name, "services", "data", "federated", "provider")
	resp, err := client.Post(endpoint, values)
	if resp != nil {
		defer func() { _ = resp.Body.Close() }()
	}
	return err
}

// DeleteFederatedProvider removes a federated provider definition.
func (client *Client) DeleteFederatedProvider(name string) (*http.Response, error) {
	endpoint := client.BuildSplunkURLWithEscapedPathPart(nil, name, "services", "data", "federated", "provider")
	return client.Delete(endpoint)
}
