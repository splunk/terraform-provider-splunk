package client

import (
	"github.com/google/go-querystring/query"
	"net/http"
	"terraform-provider-splunk/client/models"
)

func (client *Client) CreateHttpEventCollectorObject(name string, owner string, app string, httpInputConfigObj *models.HttpEventCollectorObject) (*http.Response, error) {
	/*
		Response: {"links":{"create":"/servicesNS/nobody/splunk_httpinput/data/inputs/http/_new","_reload":"/servicesNS/nobody/splunk_httpinput/data/inputs/http/_reload","_acl":"/servicesNS/nobody/splunk_httpinput/data/inputs/http/_acl"},"origin":"https://localhost:8089/servicesNS/nobody/splunk_httpinput/data/inputs/http","updated":"2020-07-27T10:42:08-07:00","generator":{"build":"a6754d8441bf","version":"8.0.3"},"entry":[{"name":"http://new_token","id":"https://localhost:8089/servicesNS/nobody/splunk_httpinput/data/inputs/http/http%3A%252F%252Fnew_token","updated":"2020-07-27T10:42:08-07:00","links":{"alternate":"/servicesNS/nobody/splunk_httpinput/data/inputs/http/http%3A%252F%252Fnew_token","list":"/servicesNS/nobody/splunk_httpinput/data/inputs/http/http%3A%252F%252Fnew_token","_reload":"/servicesNS/nobody/splunk_httpinput/data/inputs/http/http%3A%252F%252Fnew_token/_reload","edit":"/servicesNS/nobody/splunk_httpinput/data/inputs/http/http%3A%252F%252Fnew_token","remove":"/servicesNS/nobody/splunk_httpinput/data/inputs/http/http%3A%252F%252Fnew_token","disable":"/servicesNS/nobody/splunk_httpinput/data/inputs/http/http%3A%252F%252Fnew_token/disable"},"author":"nobody","acl":{"app":"splunk_httpinput","can_change_perms":true,"can_list":true,"can_share_app":true,"can_share_global":true,"can_share_user":false,"can_write":true,"modifiable":true,"owner":"nobody","perms":{"read":["admin"],"write":["admin"]},"removable":true,"sharing":"app"},"content":{"_rcvbuf":1572864,"disabled":false,"eai:acl":null,"eai:appName":"splunk_httpinput","eai:userName":"nobody","host":"ajayaraman-MBP-6E14B","index":"main","indexes":[],"source":"new","sourcetype":"new","token":"04c7c047-f5af-4a88-848c-96df7690d939","useack":"false"}}],"paging":{"total":1,"perPage":30,"offset":0},"messages":[]}
	*/
	values, err := query.Values(httpInputConfigObj)
	values.Add("name", name)
	if err != nil {
		return nil, err
	}

	endpoint := client.BuildSplunkURL(nil, "servicesNS", owner, app, "data", "inputs", "http", name)
	resp, err := client.Post(endpoint, values)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (client *Client) ReadHttpEventCollectorObject(name, owner, app string) (*http.Response, error) {
	endpoint := client.BuildSplunkURL(nil, "servicesNS", owner, app, "data", "inputs", "http", name)
	resp, err := client.Get(endpoint)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (client *Client) UpdateHttpEventCollectorObject(name string, owner string, app string, httpInputConfigObj *models.HttpEventCollectorObject) (*http.Response, error) {
	values, err := query.Values(&httpInputConfigObj)
	endpoint := client.BuildSplunkURL(nil, "servicesNS", owner, app, "data", "inputs", "http", name)
	resp, err := client.Post(endpoint, values)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (client *Client) DeleteHttpEventCollectorObject(name, owner, app string) (*http.Response, error) {
	endpoint := client.BuildSplunkURL(nil, "servicesNS", owner, app, "data", "inputs", "http", name)
	resp, err := client.Delete(endpoint)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
