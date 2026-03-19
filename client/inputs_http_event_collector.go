package client

import (
	"github.com/splunk/terraform-provider-splunk/client/models"
	"net/http"

	"github.com/google/go-querystring/query"
)

func (client *Client) CreateHttpEventCollectorObject(name string, owner string, app string, httpInputConfigObj *models.HttpEventCollectorObject) error {
	values, err := query.Values(httpInputConfigObj)
	if err != nil {
		return err
	}
	values.Add("name", name)
	endpoint := client.BuildSplunkURL(nil, "servicesNS", owner, app, "data", "inputs", "http", name)
	resp, err := client.Post(endpoint, values)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (client *Client) ReadHttpEventCollectorObject(name, owner, app string) (*http.Response, error) {
	endpoint := client.BuildSplunkURL(nil, "servicesNS", owner, app, "data", "inputs", "http", name)
	resp, err := client.Get(endpoint)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (client *Client) UpdateHttpEventCollectorObject(name string, owner string, app string, httpInputConfigObj *models.HttpEventCollectorObject) error {
	values, err := query.Values(&httpInputConfigObj)
	if err != nil {
		return err
	}
	endpoint := client.BuildSplunkURL(nil, "servicesNS", owner, app, "data", "inputs", "http", name)
	resp, err := client.Post(endpoint, values)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (client *Client) DeleteHttpEventCollectorObject(name, owner, app string) (*http.Response, error) {
	endpoint := client.BuildSplunkURL(nil, "servicesNS", owner, app, "data", "inputs", "http", name)
	resp, err := client.Delete(endpoint)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// ReadAllHttpEventCollectorObject returns the list of HTTP Event Collector tokens.
// Per Splunk RESTREF (https://docs.splunk.com/Documentation/Splunk/latest/RESTREF/RESTinput#data.2Finputs.2Fhttp),
// the documented list endpoint is GET /services/data/inputs/http (global, no servicesNS).
// GET /servicesNS/-/-/data/inputs/http returns empty in some deployments (see issue #56).
func (client *Client) ReadAllHttpEventCollectorObject() (*http.Response, error) {
	endpoint := client.BuildSplunkURL(nil, "services", "data", "inputs", "http")
	resp, err := client.Get(endpoint)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
