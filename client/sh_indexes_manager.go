package client

import (
	"net/http"
	"net/url"

	"github.com/rsrdesarrollo/terraform-provider-splunk/client/models"

	"github.com/google/go-querystring/query"
)

func (client *Client) CreateShIndexesManagerObject(name string, owner string, app string, indexConfigObj *models.ShIndexesManagerObject) error {
	values, err := query.Values(indexConfigObj)
	if err != nil {
		return err
	}
	values.Add("name", name)
	endpoint := client.BuildSplunkURL(nil, "servicesNS", owner, app, "cluster_blaster_indexes", "sh_indexes_manager", name)
	resp, err := client.Post(endpoint, values)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (client *Client) ReadShIndexesManagerObject(name, owner, app string) (*http.Response, error) {
	queryValues := url.Values{}
	queryValues.Add("datatype", "all")

	endpoint := client.BuildSplunkURL(queryValues, "servicesNS", owner, app, "cluster_blaster_indexes", "sh_indexes_manager", name)
	resp, err := client.Get(endpoint)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (client *Client) UpdateShIndexesManagerObject(name string, owner string, app string, indexConfigObj *models.ShIndexesManagerObject) error {
	values, err := query.Values(indexConfigObj)
	values.Del("datatype")
	if err != nil {
		return err
	}

	endpoint := client.BuildSplunkURL(nil, "servicesNS", owner, app, "cluster_blaster_indexes", "sh_indexes_manager", name)
	resp, err := client.Post(endpoint, values)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (client *Client) DeleteShIndexesManagerObject(name, owner, app string) (*http.Response, error) {
	endpoint := client.BuildSplunkURL(nil, "servicesNS", owner, app, "cluster_blaster_indexes", "sh_indexes_manager", name)
	resp, err := client.Delete(endpoint)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (client *Client) ReadAllShIndexesManagerObject() (*http.Response, error) {
	queryValues := url.Values{}
	queryValues.Add("datatype", "all")

	endpoint := client.BuildSplunkURL(queryValues, "servicesNS", "-", "-", "cluster_blaster_indexes", "sh_indexes_manager")
	resp, err := client.Get(endpoint)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
