package client

import (
	"net/http"
	"github.com/terraform-providers/terraform-provider-splunk/client/models"

	"github.com/google/go-querystring/query"
)

func (client *Client) CreateAuthorizationRoles(name string, authorizationRolesObject *models.AuthorizationRolesObject) error {
	values, err := query.Values(authorizationRolesObject)
	if err != nil {
		return err
	}

	endpoint := client.BuildSplunkURL(nil, "services", "authorization", "roles")
	resp, err := client.Post(endpoint, values)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (client *Client) ReadAuthorizationRoles(name, owner, app string) (*http.Response, error) {
	endpoint := client.BuildSplunkURL(nil, "services", "authorization", "roles", name)
	resp, err := client.Get(endpoint)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (client *Client) UpdateAuthorizationRoles(name string, authorizationRolesObject *models.AuthorizationRolesObject) error {
	values, err := query.Values(&authorizationRolesObject)
	if err != nil {
		return err
	}
	// Not required for updating user information
	values.Del("name")
	endpoint := client.BuildSplunkURL(nil, "services", "authorization", "roles", name)
	resp, err := client.Post(endpoint, values)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (client *Client) DeleteAuthorizationRoles(name string) (*http.Response, error) {
	endpoint := client.BuildSplunkURL(nil, "services", "authorization", "roles", name)
	resp, err := client.Delete(endpoint)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// services/authorization/roles
func (client *Client) ReadAllAuthorizationRoles() (*http.Response, error) {
	endpoint := client.BuildSplunkURL(nil, "services", "authorization", "roles")
	resp, err := client.Get(endpoint)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
