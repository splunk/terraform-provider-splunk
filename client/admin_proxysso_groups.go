package client

import (
	"net/http"

	"github.com/google/go-querystring/query"
	"github.com/splunk/terraform-provider-splunk/client/models"
)

// CreateAdminProxyssoGroups posts a new mapping to the admin/ProxySSO-groups collection.
// The group name belongs in the form body. Posting it in the path is the update endpoint
// and fails when the group does not exist yet.
func (client *Client) CreateAdminProxyssoGroups(name string, adminProxyssoGroupsObject *models.AdminProxyssoGroupsObject) error {
	if adminProxyssoGroupsObject.Name == "" {
		adminProxyssoGroupsObject.Name = name
	}
	values, err := query.Values(adminProxyssoGroupsObject)
	if err != nil {
		return err
	}

	endpoint := client.BuildSplunkURL(nil, "services", "admin", "ProxySSO-groups")
	resp, err := client.Post(endpoint, values)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

// ReadAdminProxyssoGroups fetches one Proxy SSO group. The name is escaped once.
func (client *Client) ReadAdminProxyssoGroups(name string) (*http.Response, error) {
	endpoint := client.BuildSplunkURLWithEscapedPathPart(nil, name, "services", "admin", "ProxySSO-groups")
	resp, err := client.Get(endpoint)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// UpdateAdminProxyssoGroups posts role changes to the named Proxy SSO group.
func (client *Client) UpdateAdminProxyssoGroups(name string, adminProxyssoGroupsObject *models.AdminProxyssoGroupsObject) error {
	values, err := query.Values(adminProxyssoGroupsObject)
	if err != nil {
		return err
	}
	// The group name is the URL path, not a form field, on update.
	values.Del("name")
	endpoint := client.BuildSplunkURLWithEscapedPathPart(nil, name, "services", "admin", "ProxySSO-groups")
	resp, err := client.Post(endpoint, values)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

// DeleteAdminProxyssoGroups removes one Proxy SSO group. The name is escaped once.
func (client *Client) DeleteAdminProxyssoGroups(name string) (*http.Response, error) {
	endpoint := client.BuildSplunkURLWithEscapedPathPart(nil, name, "services", "admin", "ProxySSO-groups")
	resp, err := client.Delete(endpoint)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
