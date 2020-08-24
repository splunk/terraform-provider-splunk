package client

import (
	"github.com/google/go-querystring/query"
	"net/http"
	"terraform-provider-splunk/client/models"
	"strings"
)

func (client *Client) CreateConfStanzaObject(name string, owner string, app string, confStanzaConfigObj *models.ConfStanzaObject) error {
	values, err := query.Values(confStanzaConfigObj)

	for k, v := range confStanzaConfigObj.Variables {
		values.Add(string(k), string(v))
	}
	split_name := strings.Split(name, "/")
	conf_name := split_name[0]
	stanza_name := split_name[1]
	values.Add("name", stanza_name)
	values.Del("variables")

	endpoint := client.BuildSplunkURL(nil, "servicesNS", owner, app, "configs", "conf-" + conf_name)
	resp, err := client.Post(endpoint, values)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (client *Client) ReadConfStanzaObject(name, owner, app string) (*http.Response, error) {
	split_name := strings.Split(name, "/")
	conf_name := split_name[0]
	stanza_name := split_name[1]

	endpoint := client.BuildSplunkURL(nil, "servicesNS", owner, app, "configs", "conf-" + conf_name, stanza_name)
	resp, err := client.Get(endpoint)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (client *Client) UpdateConfStanzaObject(name string, owner string, app string, confStanzaConfigObj *models.ConfStanzaObject) error {
	values, err := query.Values(&confStanzaConfigObj)

	for k, v := range confStanzaConfigObj.Variables {
		values.Add(string(k), string(v))
	}

	split_name := strings.Split(name, "/")
	conf_name := split_name[0]
	stanza_name := split_name[1]
	values.Del("variables")

	endpoint := client.BuildSplunkURL(nil, "servicesNS", owner, app, "configs", "conf-" + conf_name, stanza_name)
	resp, err := client.Post(endpoint, values)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (client *Client) DeleteConfStanzaObject(name, owner, app string) (*http.Response, error) {
	split_name := strings.Split(name, "/")
	conf_name := split_name[0]
	stanza_name := split_name[1]

	endpoint := client.BuildSplunkURL(nil, "servicesNS", owner, app, "configs", "conf-" + conf_name, stanza_name)

	resp, err := client.Delete(endpoint)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (client *Client) ReadAllConfStanzaObject(name string) (*http.Response, error) {
	split_name := strings.Split(name, "/")
	conf_name := split_name[0]

	endpoint := client.BuildSplunkURL(nil, "services", "configs", "conf-" + conf_name)

	resp, err := client.Get(endpoint)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
