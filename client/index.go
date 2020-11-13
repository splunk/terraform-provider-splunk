package client

import (
	"net/http"
	"github.com/terraform-providers/terraform-provider-splunk/client/models"

	"github.com/google/go-querystring/query"
)

func (client *Client) CreateIndexObject(name string, owner string, app string, indexConfigObj *models.IndexObject) error {
	values, err := query.Values(indexConfigObj)
	if err != nil {
		return err
	}
	values.Add("name", name)
	endpoint := client.BuildSplunkURL(nil, "servicesNS", owner, app, "data", "indexes", name)
	resp, err := client.Post(endpoint, values)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (client *Client) ReadIndexObject(name, owner, app string) (*http.Response, error) {
	endpoint := client.BuildSplunkURL(nil, "servicesNS", owner, app, "data", "indexes", name)
	resp, err := client.Get(endpoint)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (client *Client) UpdateIndexObject(name string, owner string, app string, indexConfigObj *models.IndexObject) error {
	values, err := query.Values(&indexConfigObj)
	if err != nil {
		return err
	}
	values.Del("coldPath")
	values.Del("datatype")
	values.Del("homePath")
	values.Del("thawedPath")
	values.Del("tstatsHomePath")
	endpoint := client.BuildSplunkURL(nil, "servicesNS", owner, app, "data", "indexes", name)
	resp, err := client.Post(endpoint, values)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (client *Client) DeleteIndexObject(name, owner, app string) (*http.Response, error) {
	endpoint := client.BuildSplunkURL(nil, "servicesNS", owner, app, "data", "indexes", name)
	resp, err := client.Delete(endpoint)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (client *Client) ReadAllIndexObject() (*http.Response, error) {
	endpoint := client.BuildSplunkURL(nil, "servicesNS", "-", "-", "data", "indexes")
	resp, err := client.Get(endpoint)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
