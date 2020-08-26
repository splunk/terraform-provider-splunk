package client

import (
	"github.com/google/go-querystring/query"
	"net/http"
	"terraform-provider-splunk/client/models"
)

func (client *Client) CreateConfStanzaObject(name string, owner string, app string, confStanzaConfigObj *models.ConfStanzaObject) error {
	values, err := query.Values(confStanzaConfigObj)

	for k, v := range confStanzaConfigObj.Variables {
		values.Add(string(k), string(v))
	}

	conf, stanza := client.SplitConfStanza(name)
	values.Add("name", stanza)
	values.Del("Variables")

	endpoint := client.BuildSplunkURL(nil, "servicesNS", owner, app, "configs", "conf-" + conf)
	resp, err := client.Post(endpoint, values)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (client *Client) ReadConfStanzaObject(name, owner, app string) (*http.Response, error) {
	conf, stanza := client.SplitConfStanza(name)

	endpoint := client.BuildSplunkURL(nil, "servicesNS", owner, app, "configs", "conf-" + conf, stanza)
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

	conf, stanza := client.SplitConfStanza(name)
	values.Del("Variables")

	endpoint := client.BuildSplunkURL(nil, "servicesNS", owner, app, "configs", "conf-" + conf, stanza)
	resp, err := client.Post(endpoint, values)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (client *Client) DeleteConfStanzaObject(name, owner, app string) (*http.Response, error) {
	conf, stanza := client.SplitConfStanza(name)

	endpoint := client.BuildSplunkURL(nil, "servicesNS", owner, app, "configs", "conf-" + conf, stanza)

	resp, err := client.Delete(endpoint)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (client *Client) ReadAllConfStanzaObject(name string) (*http.Response, error) {
	conf, _ := client.SplitConfStanza(name)

	endpoint := client.BuildSplunkURL(nil, "services", "configs", "conf-" + conf)

	resp, err := client.Get(endpoint)
	if err != nil {
		return nil, err
	}

	return resp, nil
}