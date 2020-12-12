package client

import (
	"github.com/jaware-splunk/terraform-provider-splunk/client/models"
	"net/http"
	"strings"

	"github.com/google/go-querystring/query"
)

func (client *Client) CreateConfigsConfObject(name string, owner string, app string, configsConfConfigObj *models.ConfigsConfObject) error {
	values, err := query.Values(configsConfConfigObj)
	if err != nil {
		return err
	}

	for k, v := range configsConfConfigObj.Variables {
		values.Add(string(k), string(v))
	}

	conf, stanza := client.SplitConfStanza(name)
	values.Add("name", stanza)
	values.Del("Variables")

	endpoint := client.BuildSplunkURL(nil, "servicesNS", owner, app, "configs", "conf-"+conf)
	resp, err := client.Post(endpoint, values)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (client *Client) ReadConfigsConfObject(name, owner, app string) (*http.Response, error) {
	conf, stanza := client.SplitConfStanza(name)

	endpoint := client.BuildSplunkURL(nil, "servicesNS", owner, app, "configs", "conf-"+conf, stanza)
	resp, err := client.Get(endpoint)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (client *Client) UpdateConfigsConfObject(name string, owner string, app string, configsConfConfigObj *models.ConfigsConfObject) error {
	values, err := query.Values(&configsConfConfigObj)
	if err != nil {
		return err
	}

	for k, v := range configsConfConfigObj.Variables {
		values.Add(string(k), string(v))
	}

	conf, stanza := client.SplitConfStanza(name)
	values.Del("Variables")

	endpoint := client.BuildSplunkURL(nil, "servicesNS", owner, app, "configs", "conf-"+conf, stanza)
	resp, err := client.Post(endpoint, values)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (client *Client) DeleteConfigsConfObject(name, owner, app string) (*http.Response, error) {
	conf, stanza := client.SplitConfStanza(name)

	endpoint := client.BuildSplunkURL(nil, "servicesNS", owner, app, "configs", "conf-"+conf, stanza)

	resp, err := client.Delete(endpoint)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (client *Client) ReadAllConfigsConfObject(name string) (*http.Response, error) {
	conf, _ := client.SplitConfStanza(name)

	endpoint := client.BuildSplunkURL(nil, "servicesNS", "-", "-", "configs", "conf-"+conf)

	resp, err := client.Get(endpoint)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// Takes a '/' separated string and returns the 0, 1 indexed strings from the split
func (client *Client) SplitConfStanza(name string) (conf string, stanza string) {
	split := strings.Split(name, "/")
	return split[0], split[1]
}
