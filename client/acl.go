package client

import (
	"github.com/google/go-querystring/query"
	"net/http"
	"strings"
	"terraform-provider-splunk/client/models"
)

//https://docs.splunk.com/Documentation/Splunk/8.0.4/RESTUM/RESTusing#Access_Control_List
func (client *Client) GetAcl(owner, app, resource, name string) (*http.Response, error) {
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
