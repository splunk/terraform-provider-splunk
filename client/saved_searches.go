package client

import (
	"net/http"

	"github.com/splunk/terraform-provider-splunk/client/models"

	"github.com/google/go-querystring/query"
)

func (client *Client) CreateSavedSearches(name, owner, app string, savedSearchObject *models.SavedSearchObject) error {
	values, err := query.Values(savedSearchObject)
	values.Add("name", name)
	if err != nil {
		return err
	}

	endpoint := client.BuildSplunkURL(nil, "servicesNS", owner, app, "saved", "searches")
	resp, err := client.Post(endpoint, values)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (client *Client) ReadSavedSearches(name, owner, app string) (*http.Response, error) {
	endpoint := client.BuildSplunkURL(nil, "servicesNS", owner, app, "saved", "searches", name)
	resp, err := client.Get(endpoint)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (client *Client) UpdateSavedSearches(name string, owner string, app string, savedSearchObject *models.SavedSearchObject) error {
	// Splunk API doesn't update the saved search when the schedule_priority field is set (BUG)
	// As a workaround we can unset SchedulePriority in update requests, since updating this tf field forces new resource creation
	savedSearchObject.SchedulePriority = ""
	values, err := query.Values(&savedSearchObject)
	if err != nil {
		return err
	}
	endpoint := client.BuildSplunkURL(nil, "servicesNS", owner, app, "saved", "searches", name)
	resp, err := client.Post(endpoint, values)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (client *Client) DeleteSavedSearches(name, owner, app string) (*http.Response, error) {
	endpoint := client.BuildSplunkURL(nil, "servicesNS", owner, app, "saved", "searches", name)
	resp, err := client.Delete(endpoint)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
