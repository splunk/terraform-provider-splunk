package client

import (
	"net/http"

	"github.com/nealbrown/terraform-provider-splunk/client/models"

	"github.com/google/go-querystring/query"
)

func (client *Client) CreateAuthenticationUser(name string, authenticationUserObject *models.AuthenticationUserObject) error {
	values, err := query.Values(authenticationUserObject)
	if err != nil {
		return err
	}

	endpoint := client.BuildSplunkURL(nil, "services", "authentication", "users")
	resp, err := client.Post(endpoint, values)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (client *Client) ReadAuthenticationUser(name, owner, app string) (*http.Response, error) {
	endpoint := client.BuildSplunkURL(nil, "services", "authentication", "users", name)
	resp, err := client.Get(endpoint)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (client *Client) UpdateAuthenticationUser(name string, authenticationUserObject *models.AuthenticationUserObject) error {
	values, err := query.Values(&authenticationUserObject)
	if err != nil {
		return err
	}
	// Not required for updating user information
	values.Del("name")
	endpoint := client.BuildSplunkURL(nil, "services", "authentication", "users", name)
	resp, err := client.Post(endpoint, values)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (client *Client) DeleteAuthenticationUser(name string) (*http.Response, error) {
	endpoint := client.BuildSplunkURL(nil, "services", "authentication", "users", name)
	resp, err := client.Delete(endpoint)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// services/authentication/users
func (client *Client) ReadAuthenticationUsers() (*http.Response, error) {
	endpoint := client.BuildSplunkURL(nil, "services", "authentication", "users")
	resp, err := client.Get(endpoint)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
