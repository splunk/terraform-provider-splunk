package client

import (
	"log"
	"net/http"
	"net/http/httputil"

	"github.com/google/go-querystring/query"
	"github.com/nealbrown/terraform-provider-splunk/client/models"
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

	respBody, error := httputil.DumpResponse(resp, true)
	if error != nil {
		log.Printf("[ERROR] Error occured during CreateSavedSearches %s", error)
	}

	log.Printf("[DEBUG] Response object returned from CreateSavedSearches is: %s", string(respBody))

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

	respBody, error := httputil.DumpResponse(resp, true)
	if error != nil {
		log.Printf("[ERROR] Error occured during DeleteSavedSearches %s", error)
	}

	log.Printf("[DEBUG] Response object returned from DeleteSavedSearches is: %s", string(respBody))

	return resp, nil
}
