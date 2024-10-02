package client

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
)

func (client *Client) CreateLookupTableFile(name string, owner string, app string, contents string) error {
	values := []byte(fmt.Sprintf("namespace=%s&lookup_file=%s&owner=%s&contents=%s", app, name, owner, contents))
	endpoint := client.BuildSplunkURL(nil, "services", "data", "lookup_edit", "lookup_contents")
	client.urlEncoded = true
	resp, err := client.Post(endpoint, values)
	if err != nil {
		return err
	}

	respBody, error := httputil.DumpResponse(resp, true)
	if error != nil {
		log.Printf("[ERROR] Error occured during CreateLookup %s", error)
	}

	log.Printf("[DEBUG] Response object returned from CreateLookup is: %s", string(respBody))

	defer resp.Body.Close()
	return nil
}

func (client *Client) ReadLookupTableFile(name, owner, app string) (*http.Response, error) {
	values := []byte(fmt.Sprintf("namespace=%s&lookup_file=%s&owner=%s", app, name, owner))
	client.urlEncoded = true
	endpoint := client.BuildSplunkURL(nil, "services", "data", "lookup_edit", "lookup_data")
	resp, err := client.Post(endpoint, values)
	return resp, err
}

func (client *Client) UpdateLookupTableFile(name string, owner string, app string, contents string) error {
	values := []byte(fmt.Sprintf("namespace=%s&lookup_file=%s&owner=%s&contents=%s", app, name, owner, contents))
	endpoint := client.BuildSplunkURL(nil, "services", "data", "lookup_edit", "lookup_contents")
	client.urlEncoded = true
	resp, err := client.Post(endpoint, values)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (client *Client) DeleteLookupTableFile(name string, owner string, app string) (*http.Response, error) {
	endpoint := client.BuildSplunkURL(nil, "servicesNS", owner, app, "data", "lookup-table-files", name)
	resp, err := client.Delete(endpoint)
	if err != nil {
		return nil, err
	}

	respBody, error := httputil.DumpResponse(resp, true)
	if error != nil {
		log.Printf("[ERROR] Error occured during DeleteLookup %s", error)
	}

	log.Printf("[DEBUG] Response object returned from DeleteLookup is: %s", string(respBody))

	return resp, nil
}
