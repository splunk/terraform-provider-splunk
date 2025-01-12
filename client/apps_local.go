package client

import (
	"net/http"

	"github.com/nealbrown/terraform-provider-splunk/client/models"

	"github.com/google/go-querystring/query"
)

func (client *Client) CreateAppsLocalObject(app string, localAppObject *models.AppsLocalObject) error {
	values, err := query.Values(localAppObject)
	if err != nil {
		return err
	}
	endpoint := client.BuildSplunkURL(nil, "servicesNS", "nobody", "system", "apps", "local")
	resp, err := client.Post(endpoint, values)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (client *Client) ReadAppsLocalObject(app string) (*http.Response, error) {
	endpoint := client.BuildSplunkURL(nil, "servicesNS", "nobody", "system", "apps", "local", app)
	resp, err := client.Get(endpoint)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (client *Client) UpdateAppsLocalObject(app string, localAppObject *models.AppsLocalObject) error {
	values, err := query.Values(&localAppObject)
	if err != nil {
		return err
	}
	values.Del("name")             // Handler does not support "name" argument
	values.Del("update")           // Handler does not support "udpate" argument
	values.Del("filename")         // Handler does not support "filename" argument
	values.Del("explicit_appname") // Handler does not support "explicit_appname" argument
	endpoint := client.BuildSplunkURL(nil, "servicesNS", "nobody", "system", "apps", "local", app)
	resp, err := client.Post(endpoint, values)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (client *Client) DeleteAppsLocalObject(app string) (*http.Response, error) {
	endpoint := client.BuildSplunkURL(nil, "servicesNS", "nobody", "system", "apps", "local", app)
	resp, err := client.Delete(endpoint)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (client *Client) ReadAllAppsLocalObject() (*http.Response, error) {
	endpoint := client.BuildSplunkURL(nil, "servicesNS", "-", "-", "apps", "local")
	resp, err := client.Get(endpoint)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
