package client

import (
	"github.com/splunk/terraform-provider-splunk/client/models"
	"net/http"
	"net/url"
	"github.com/google/go-querystring/query"
)

func (client *Client) CreateAdminProxyssoGroups(name string, AdminProxyssoGroupsObject *models.AdminProxyssoGroupsObject) error {
	values, err := query.Values(AdminProxyssoGroupsObject)
	if err != nil {
		return err
	}

	endpoint := client.BuildSplunkURL(nil, "servicesNS", "nobody", "system", "admin", "ProxySSO-groups", url.PathEscape(name))
	resp, err := client.Post(endpoint, values)

	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (client *Client) ReadAdminProxyssoGroups(name, owner, app string) (*http.Response, error) {
	endpoint := client.BuildSplunkURL(nil, "servicesNS", "nobody", "system", "admin", "ProxySSO-groups", url.PathEscape(name))
	resp, err := client.Get(endpoint)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (client *Client) UpdateAdminProxyssoGroups(name string, AdminProxyssoGroupsObject *models.AdminProxyssoGroupsObject) error {
	values, err := query.Values(&AdminProxyssoGroupsObject)
	if err != nil {
		return err
	}
	// Not required for updating user information
	values.Del("name")
	endpoint := client.BuildSplunkURL(nil, "servicesNS", "nobody", "system", "admin", "ProxySSO-groups", url.PathEscape(name))
	resp, err := client.Post(endpoint, values)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (client *Client) DeleteAdminProxyssoGroups(name string) (*http.Response, error) {
	endpoint := client.BuildSplunkURL(nil, "servicesNS", "nobody", "system", "admin", "ProxySSO-groups", url.PathEscape(name))
	resp, err := client.Delete(endpoint)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// servicesNS/nobody/system/admin/ProxySSO-groups
func (client *Client) ReadAllAdminProxyssoGroups() (*http.Response, error) {
	endpoint := client.BuildSplunkURL(nil, "servicesNS", "nobody", "system", "admin", "ProxySSO-groups")
	resp, err := client.Get(endpoint)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
