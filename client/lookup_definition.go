package client

import (
	"log"
	"net/http"
	"net/http/httputil"

	"github.com/google/go-querystring/query"
	"github.com/splunk/terraform-provider-splunk/client/models"
)

func (client *Client) CreateLookupDefinitionObject(owner string, app string, splunkLookupDefObj *models.SplunkLookupDefinitionObject) error {
	values, err := query.Values(&splunkLookupDefObj)
	if err != nil {
		return err
	}
	endpoint := client.BuildSplunkURL(nil, "servicesNS", owner, app, "data", "transforms", "lookups")
	resp, err := client.Post(endpoint, values)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, error := httputil.DumpResponse(resp, true)
	if error != nil {
		log.Printf("[ERROR] Error occured during CreateLookupDefinition %s", error)
	}

	log.Printf("[DEBUG] Response object returned from CreateLookupDefinition is: %s", string(respBody))
	return nil
}

func (client *Client) ReadLookupDefinitionObject(name, owner, app string) (*http.Response, error) {
	endpoint := client.BuildSplunkURL(nil, "servicesNS", owner, app, "data", "transforms", "lookups", name)
	resp, err := client.Get(endpoint)
	requestBody, _ := httputil.DumpRequest(resp.Request, true)
	if err != nil {
		log.Printf("[ERROR] Error occured during ReadLookupDefinitionObject %s", string(requestBody))
		return nil, err
	}
	respBody, error := httputil.DumpResponse(resp, true)
	if error != nil {
		log.Printf("[ERROR] Error occured during ReadLookupDefinitionObject %s", error)
	}

	log.Printf("[DEBUG] Request object returned from ReadLookupDefinitionObject is: %s", string(requestBody))
	log.Printf("[DEBUG] Response object returned from ReadLookupDefinitionObject is: %s", string(respBody))

	return resp, nil
}

func (client *Client) UpdateLookupDefinitionObject(owner string, app string, name string, splunkLookupDefObj *models.SplunkLookupDefinitionObject) error {
	values, err := query.Values(&splunkLookupDefObj)
	values.Del("name")
	if err != nil {
		return err
	}
	endpoint := client.BuildSplunkURL(nil, "servicesNS", owner, app, "data", "transforms", "lookups", name)
	resp, err := client.Post(endpoint, values)
	if err != nil {
		return err
	}
	respBody, error := httputil.DumpResponse(resp, true)
	if error != nil {
		log.Printf("[ERROR] Error occured during UpdateLookupDefinitionObject %s", error)
	}

	log.Printf("[DEBUG] Response object returned from UpdateLookupDefinitionObject is: %s", string(respBody))

	defer resp.Body.Close()

	return nil
}

func (client *Client) DeleteLookupDefinitionObject(owner string, app string, name string) (*http.Response, error) {
	endpoint := client.BuildSplunkURL(nil, "servicesNS", owner, app, "data", "transforms", "lookups", name)
	resp, err := client.Delete(endpoint)
	if err != nil {
		return nil, err
	}
	respBody, error := httputil.DumpResponse(resp, true)
	if error != nil {
		log.Printf("[ERROR] Error occured during DeleteLookupDefinitionObject %s", error)
	}

	log.Printf("[DEBUG] Response object returned from DeleteLookupDefinitionObject is: %s", string(respBody))

	defer resp.Body.Close()

	return resp, nil
}
