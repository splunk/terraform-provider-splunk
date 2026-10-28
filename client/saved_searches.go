package client

import (
	"github.com/google/go-querystring/query"
	"github.com/splunk/terraform-provider-splunk/client/models"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func (client *Client) CreateSavedSearches(name, owner, app string, savedSearchObject *models.SavedSearchObject) (*http.Response, error) {
	values, err := query.Values(savedSearchObject)
	values.Add("name", name)
	if err != nil {
		return nil, err
	}

	endpoint := client.BuildSplunkURL(nil, "servicesNS", owner, app, "saved", "searches")
	resp, err := client.Post(endpoint, values)
	if err != nil {
		return resp, err
	}

	respBody, error := httputil.DumpResponse(resp, true)
	if error != nil {
		log.Printf("[ERROR] Error occured during CreateSavedSearches %s", error)
	}

	log.Printf("[DEBUG] Response object returned from CreateSavedSearches is: %s", string(respBody))

	// Leave the body open. DumpResponse replaced it with a readable copy, and the caller closes it.
	return resp, nil
}

// ReadSavedSearches fetches one saved search. The response is returned together
// with any status error so the caller can tell a 404 from other failures.
func (client *Client) ReadSavedSearches(name, owner, app string) (*http.Response, error) {
	endpoint := client.BuildSplunkURLWithEscapedPathPart(nil, name, "servicesNS", owner, app, "saved", "searches")
	return client.Get(endpoint)
}

// UpdateSavedSearches posts the saved search. Splunk keeps any key an update
// leaves out, so each key in clearKeys is posted with an empty value to clear it.
func (client *Client) UpdateSavedSearches(name string, owner string, app string, savedSearchObject *models.SavedSearchObject, clearKeys []string) error {
	values, err := query.Values(&savedSearchObject)
	if err != nil {
		return err
	}
	// EncodeRequestBody drops empty values, so the cleared keys are appended after it.
	body, err := client.EncodeRequestBody(values)
	if err != nil {
		return err
	}
	cleared := url.Values{}
	for _, key := range clearKeys {
		cleared.Set(key, "")
	}
	if len(cleared) > 0 {
		if len(body) > 0 {
			body = append(body, '&')
		}
		body = append(body, cleared.Encode()...)
	}
	endpoint := client.BuildSplunkURLWithEscapedPathPart(nil, name, "servicesNS", owner, app, "saved", "searches")
	resp, err := client.Post(endpoint, body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (client *Client) DeleteSavedSearches(name, owner, app string) (*http.Response, error) {
	endpoint := client.BuildSplunkURLWithEscapedPathPart(nil, name, "servicesNS", owner, app, "saved", "searches")
	resp, err := client.Delete(endpoint)
	if err != nil {
		return resp, err
	}

	respBody, error := httputil.DumpResponse(resp, true)
	if error != nil {
		log.Printf("[ERROR] Error occured during DeleteSavedSearches %s", error)
	}

	log.Printf("[DEBUG] Response object returned from DeleteSavedSearches is: %s", string(respBody))

	return resp, nil
}
