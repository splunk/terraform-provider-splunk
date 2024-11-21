package client

import (
	"github.com/splunk/terraform-provider-splunk/client/models"
	"net/http"

	"github.com/google/go-querystring/query"
)

func (client *Client) CreateSavedEventTypes(name, owner, app string, savedEventTypeObject *models.SavedEventTypeObject) error {
	values, err := query.Values(savedEventTypeObject)
	values.Add("name", name)
	if err != nil {
		return err
	}

	endpoint := client.BuildSplunkURL(nil, "servicesNS", owner, app, "saved", "eventtypes")
	resp, err := client.Post(endpoint, values)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (client *Client) ReadSavedEventTypes(name, owner, app string) (*http.Response, error) {
	endpoint := client.BuildSplunkURL(nil, "servicesNS", owner, app, "saved", "eventtypes", name)
	resp, err := client.Get(endpoint)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (client *Client) UpdateSavedEventTypes(name string, owner string, app string, savedEventTypeObject *models.SavedEventTypeObject) error {
	values, err := query.Values(&savedEventTypeObject)
	if err != nil {
		return err
	}
	endpoint := client.BuildSplunkURL(nil, "servicesNS", owner, app, "saved", "eventtypes", name)
	resp, err := client.Post(endpoint, values)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (client *Client) DeleteSavedEventTypes(name, owner, app string) (*http.Response, error) {
	endpoint := client.BuildSplunkURL(nil, "servicesNS", owner, app, "saved", "eventtypes", name)
	resp, err := client.Delete(endpoint)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// services/saved/eventtypes
func (client *Client) ReadAllSavedEventTypes() (*http.Response, error) {
	endpoint := client.BuildSplunkURL(nil, "servicesNS", "-", "-", "saved", "eventtypes")
	resp, err := client.Get(endpoint)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
