package client

import (
	"net/http"

	"github.com/splunk/terraform-provider-splunk/client/models"

	"github.com/google/go-querystring/query"
)

func (client *Client) CreateDashboardObject(owner string, app string, splunkDashboardsObj *models.SplunkDashboardsObject) error {
	values, err := query.Values(&splunkDashboardsObj)
	if err != nil {
		return err
	}
	endpoint := client.BuildSplunkURL(nil, "servicesNS", owner, app, "data", "ui", "views")
	resp, err := client.Post(endpoint, values)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (client *Client) ReadDashboardObject(name, owner, app string) (*http.Response, error) {
	endpoint := client.BuildSplunkURL(nil, "servicesNS", owner, app, "data", "ui", "views", name)
	resp, err := client.Get(endpoint)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (client *Client) UpdateDashboardObject(owner string, app string, name string, splunkDashboardsObj *models.SplunkDashboardsObject) error {
	values, err := query.Values(&splunkDashboardsObj)
	values.Del("name")
	if err != nil {
		return err
	}
	endpoint := client.BuildSplunkURL(nil, "servicesNS", owner, app, "data", "ui", "views", name)
	resp, err := client.Post(endpoint, values)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (client *Client) DeleteDashboardObject(owner string, app string, name string) (*http.Response, error) {
	endpoint := client.BuildSplunkURL(nil, "servicesNS", owner, app, "data", "ui", "views", name)
	resp, err := client.Delete(endpoint)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return resp, nil
}

func (client *Client) ReadAllDashboardObject() (*http.Response, error) {
	endpoint := client.BuildSplunkURL(nil, "servicesNS", "-", "-", "data", "ui", "views")
	resp, err := client.Get(endpoint)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
