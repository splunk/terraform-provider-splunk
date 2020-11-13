package client

import (
	"net/http"
	"github.com/terraform-providers/terraform-provider-splunk/client/models"

	"github.com/google/go-querystring/query"
)

func (client *Client) CreateAdminSAMLGroups(name string, adminSAMLGroupsObject *models.AdminSAMLGroupsObject) error {
	values, err := query.Values(adminSAMLGroupsObject)
	if err != nil {
		return err
	}

	endpoint := client.BuildSplunkURL(nil, "services", "admin", "SAML-groups")
	resp, err := client.Post(endpoint, values)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (client *Client) ReadAdminSAMLGroups(name string) (*http.Response, error) {
	endpoint := client.BuildSplunkURL(nil, "services", "admin", "SAML-groups", name)
	resp, err := client.Get(endpoint)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (client *Client) UpdateAdminSAMLGroups(name string, adminSAMLGroupsObject *models.AdminSAMLGroupsObject) error {
	values, err := query.Values(&adminSAMLGroupsObject)
	if err != nil {
		return err
	}
	// Not required for updates
	values.Del("name")
	endpoint := client.BuildSplunkURL(nil, "services", "admin", "SAML-groups", name)
	resp, err := client.Post(endpoint, values)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (client *Client) DeleteAdminSAMLGroups(name string) (*http.Response, error) {
	endpoint := client.BuildSplunkURL(nil, "services", "admin", "SAML-groups", name)
	resp, err := client.Delete(endpoint)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
