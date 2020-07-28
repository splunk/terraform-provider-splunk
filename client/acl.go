package client

import (
	"github.com/google/go-querystring/query"
	"net/http"
	"strings"
	"terraform-provider-splunk/client/models"
)

//https://docs.splunk.com/Documentation/Splunk/8.0.4/RESTUM/RESTusing#Access_Control_List
func (client *Client) GetAcl(owner, app, resource, name string) (*http.Response, error) {
	/*
		Response: {"links":{"create":"/services/data/inputs/http/_new","_reload":"/services/data/inputs/http/_reload","_acl":"/services/data/inputs/http/_acl"},"origin":"https://localhost:8089/services/data/inputs/http","updated":"2020-07-28T12:27:41-07:00","generator":{"build":"a6754d8441bf","version":"8.0.3"},"entry":[{"name":"http://new-token","id":"https://localhost:8089/servicesNS/nobody/splunk_httpinput/data/inputs/http/http%3A%252F%252Fnew-token","updated":"2020-07-28T12:26:53-07:00","links":{"alternate":"/servicesNS/nobody/splunk_httpinput/data/inputs/http/http%3A%252F%252Fnew-token","list":"/servicesNS/nobody/splunk_httpinput/data/inputs/http/http%3A%252F%252Fnew-token","_reload":"/servicesNS/nobody/splunk_httpinput/data/inputs/http/http%3A%252F%252Fnew-token/_reload","edit":"/servicesNS/nobody/splunk_httpinput/data/inputs/http/http%3A%252F%252Fnew-token","remove":"/servicesNS/nobody/splunk_httpinput/data/inputs/http/http%3A%252F%252Fnew-token","disable":"/servicesNS/nobody/splunk_httpinput/data/inputs/http/http%3A%252F%252Fnew-token/disable"},"author":"splunker","acl":{"app":"splunk_httpinput","can_change_perms":true,"can_list":true,"can_share_app":true,"can_share_global":true,"can_share_user":true,"can_write":true,"modifiable":true,"owner":"splunker","perms":{"read":["*"],"write":["admin","user"]},"removable":true,"sharing":"app"},"content":{"eai:acl":null}}],"paging":{"total":1,"perPage":30,"offset":0},"messages":[]}
	*/
	endpoint := client.BuildSplunkURL(nil, "servicesNS", owner, app, "data", "inputs", resource, name, "acl")
	resp, err := client.Get(endpoint)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (client *Client) UpdateAcl(owner, app, resource, name string, acl *models.ACLObject) (*http.Response, error) {
	values, err := query.Values(&acl)
	// remove app from url values during POST
	values.Del("app")
	values.Del("perms[read]")
	values.Del("perms[write]")
	// Flatten []string
	values.Set("perms.read", strings.Join(acl.Perms.Read, ","))
	values.Set("perms.write", strings.Join(acl.Perms.Write, ","))
	endpoint := client.BuildSplunkURL(nil, "servicesNS", owner, app, "data", "inputs", resource, name, "acl")
	resp, err := client.Post(endpoint, values)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (client *Client) Move(owner, app, resource, name string, acl *models.ACLObject) (*http.Response, error) {
	values, err := query.Values(&acl)
	endpoint := client.BuildSplunkURL(nil, "servicesNS", owner, app, "data", "inputs", "http", resource, name, "move")
	resp, err := client.Post(endpoint, values)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
