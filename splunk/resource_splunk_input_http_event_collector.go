package splunk

import (
	"encoding/json"
	"errors"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
)

const (
  splunkHttpInputEndpoint = "/services/data/inputs/http/"
)

type SplunkHttpInputConfig struct {
	Name        string
	Index       string
	Source      string
	Sourcetype  string
	Disabled    bool
	UseACK      bool
	Token       string
}

type SuccessResponse struct {
	Entry []Entry `json:"entry"`
	Messages []ErrorMessage `json:"messages"`
}

type Entry struct {
	Name    string  `json:"name"`
	Content Content `json:"content"`
}

type Content struct {
	Token string `json:"token"`
	Index string `json:"index"`
}

type ErrorMessage struct {
	Type   string `json:"type"`
	Text   string `json:"text"`
}

type ErrorResponse struct {
	Messages []ErrorMessage `json:"messages"`
}

func inputHttpEventCollector() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "HTTP Event Collector Token name",
			},
			"index": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "main",
				Description: "Index to store generated events",
			},
			"source": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Default source for events with this token",
			},
			"sourcetype": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Default sourcetype for events with this token",
			},
			"disabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Input disabled indicator: 0 = Input Not disabled, 1 = Input disabled",
			},
			"use_ack": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Indexer acknowledgement for this token",
			},
		},
		Read:httpEventCollectorInputRead,
		Create:httpEventCollectorInputCreate,
		Delete:httpEventCollectorInputDelete,
		Update:httpEventCollectorInputUpdate,
		Exists:httpEventCollectorInputExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func httpEventCollectorInputRead(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	input := SplunkHttpInputConfig{
		Name: d.Get("name").(string),
	}

	createdSplunkHttpInput, err := (*provider.Client).doGetHttpInput(&input)
	if err != nil {
		return err
	}

	d.SetId(createdSplunkHttpInput.Name)
	_ = d.Set("token", createdSplunkHttpInput.Token)
	return nil
}

func httpEventCollectorInputCreate(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	input := SplunkHttpInputConfig{
		Name: d.Get("name").(string),
		Index:d.Get("index").(string),
		Source:d.Get("source").(string),
		Sourcetype:d.Get("sourcetype").(string),
		Disabled:d.Get("disabled").(bool),
		UseACK:d.Get("use_ack").(bool),
	}

	createdSplunkHttpInput, err := (*provider.Client).doCreateHttpInput(&input)
	if err != nil {
		return err
	}

	d.SetId(createdSplunkHttpInput.Name)
	_ = d.Set("token", createdSplunkHttpInput.Token)
	return nil
}

func httpEventCollectorInputDelete(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	input := SplunkHttpInputConfig{
		Name: d.Get("name").(string),
	}

	err := (*provider.Client).doDeleteHttpInput(&input)
	if err != nil {
		return err
	}
	return nil
}

func httpEventCollectorInputExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	return true, nil
}

func httpEventCollectorInputUpdate(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	input := SplunkHttpInputConfig{
		Name: d.Get("name").(string),
		Index:d.Get("index").(string),
		Source:d.Get("source").(string),
		Sourcetype:d.Get("sourcetype").(string),
		Disabled:d.Get("disabled").(bool),
		UseACK:d.Get("use_ack").(bool),
	}

	createdSplunkHttpInput, err := (*provider.Client).doUpdateHttpInput(&input)
	if err != nil {
		return err
	}

	d.SetId(createdSplunkHttpInput.Name)
	_ = d.Set("token", createdSplunkHttpInput.Token)
	return nil
}

// Make requests

func (c *SplunkClient) doGetHttpInput(m *SplunkHttpInputConfig) (*SplunkHttpInputConfig, error) {
	httpInputEndpoint := splunkHttpInputEndpoint + m.Name
	values := url.Values{}
	values.Add("name", m.Name)
	resp, err := c.doRequest("GET", httpInputEndpoint, values)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	switch resp.StatusCode {
	case 200, 201:
		splunkHttpInputConfig, err := unmarshalSuccessReponse(resp)
		if err != nil {
			return nil, err
		}
		return splunkHttpInputConfig, nil

	default:
		errorResponse := &ErrorResponse{}
		_ = json.NewDecoder(resp.Body).Decode(errorResponse)
		err = errors.New(errorResponse.Messages[0].Text)
		return nil, err
	}
}

func (c *SplunkClient) doCreateHttpInput(m *SplunkHttpInputConfig) (*SplunkHttpInputConfig, error) {
	httpInputEndpoint := splunkHttpInputEndpoint
	values := url.Values{}
	values.Add("name", m.Name)
	values.Add("index", m.Index)
	values.Add("source", m.Source)
	values.Add("sourcetype", m.Sourcetype)
	values.Add("useACK", strconv.FormatBool(m.UseACK))
	values.Add("disabled", strconv.FormatBool(m.Disabled))
	resp, err := c.doRequest("POST", httpInputEndpoint, values)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	switch resp.StatusCode {
	case 200, 201:
		splunkHttpInputConfig, err := unmarshalSuccessReponse(resp)
		if err != nil {
			return nil, err
		}

		return splunkHttpInputConfig, nil

	default:
		errorResponse := &ErrorResponse{}
		_ = json.NewDecoder(resp.Body).Decode(errorResponse)
		err = errors.New(errorResponse.Messages[0].Text)
		return nil, err
	}
}

func (c *SplunkClient) doUpdateHttpInput(m *SplunkHttpInputConfig) (*SplunkHttpInputConfig, error) {
	httpInputEndpoint := splunkHttpInputEndpoint + m.Name
	values := url.Values{}
	values.Add("index", m.Index)
	values.Add("source", m.Source)
	values.Add("sourcetype", m.Sourcetype)
	values.Add("disabled", strconv.FormatBool(m.Disabled))
	values.Add("useACK", strconv.FormatBool(m.UseACK))
	resp, err := c.doRequest("POST", httpInputEndpoint, values)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	switch resp.StatusCode {
	case 200, 201:
		splunkHttpInputConfig, err := unmarshalSuccessReponse(resp)
		if err != nil {
			return nil, err
		}

		return splunkHttpInputConfig, nil

	default:
		errorResponse := &ErrorResponse{}
		_ = json.NewDecoder(resp.Body).Decode(errorResponse)
		err = errors.New(errorResponse.Messages[0].Text)
		return nil, err
	}

}

func (c *SplunkClient) doDeleteHttpInput(m *SplunkHttpInputConfig) error {
	httpInputEndpoint := splunkHttpInputEndpoint + m.Name
	values := url.Values{}
	_, err := c.doRequest("DELETE", httpInputEndpoint, values)
	if err != nil {
		return err
	}
	return nil
}

func unmarshalSuccessReponse(response *http.Response) (*SplunkHttpInputConfig, error) {
	splunkHttpInputConfig := &SplunkHttpInputConfig{}
	successResponse := &SuccessResponse{}
	_ = json.NewDecoder(response.Body).Decode(&successResponse)
	re := regexp.MustCompile(`http://(.*)`)
	splunkHttpInputConfig.Name = re.FindStringSubmatch(successResponse.Entry[0].Name)[1]
	splunkHttpInputConfig.Token = successResponse.Entry[0].Content.Token
	return splunkHttpInputConfig, nil
}

