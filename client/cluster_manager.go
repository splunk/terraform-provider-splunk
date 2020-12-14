package client

import (
	"net/http"

	"github.com/jaware-splunk/terraform-provider-splunk/client/models"

	"github.com/google/go-querystring/query"
)

func (client *Client) CreateClusterManager(name string, clusterManagerObject *models.ClusterManagerObject) error {
	values, err := query.Values(clusterManagerObject)
	if err != nil {
		return err
	}

	endpoint := client.BuildSplunkURL(nil, "services", "cluster", "config", "config")
	resp, err := client.Post(endpoint, values)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (client *Client) ReadClusterManager(name string) (*http.Response, error) {
	endpoint := client.BuildSplunkURL(nil, "services", "cluster", "config", "config")
	resp, err := client.Get(endpoint)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (client *Client) UpdateClusterManager(name string, clusterManagerObject *models.ClusterManagerObject) error {
	values, err := query.Values(&clusterManagerObject)
	if err != nil {
		return err
	}
	values.Del("cluster_label")
	values.Del("mode")
	endpoint := client.BuildSplunkURL(nil, "services", "cluster", "config", "config")
	resp, err := client.Post(endpoint, values)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (client *Client) DeleteClusterManager(name string) (*http.Response, error) {
	endpoint := client.BuildSplunkURL(nil, "services", "cluster", "config", "config")
	resp, err := client.Delete(endpoint)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

//func (client *Client) ReadAllClusterManagerObject() (*http.Response, error) {
//	endpoint := client.BuildSplunkURL(nil, "servicesNS", "-", "-", "cluster", "config", "config")
//	resp, err := client.Get(endpoint)
//	if err != nil {
//		return nil, err
//	}
//
//	return resp, nil
//}
