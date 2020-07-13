package splunk

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
)

type SplunkHttpInputConfig struct {
	Name        string
	AppContext  string
	Index       string
	Source      string
	Sourcetype  string
	Disabled    bool
	UseACK      bool
	Token       string
}

func inputHttpEventCollector() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "HTTP Event Collector Token name",
			},
			"token": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "HTTP Event Collector Token name",
			},
			"app_context": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "splunk_httpinput",
				Description: "App context for the input (App must exist before creating an input)",
				ForceNew:    true,
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
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func httpEventCollectorInputRead(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	input := SplunkHttpInputConfig{
		Name: d.Get("name").(string),
		AppContext: d.Get("app_context").(string),
	}

	httpConfigObj, err := (*provider.Client).doGetHttpInput(&input)
	if err != nil {
		return err
	}

	d.SetId(httpConfigObj.Name)
	err = d.Set("token", httpConfigObj.Token)
	if err != nil {
		return err
	}
	return nil
}

func httpEventCollectorInputCreate(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	input := SplunkHttpInputConfig{
		Name: d.Get("name").(string),
		AppContext: d.Get("app_context").(string),
		Index:d.Get("index").(string),
		Source:d.Get("source").(string),
		Sourcetype:d.Get("sourcetype").(string),
		Disabled:d.Get("disabled").(bool),
		UseACK:d.Get("use_ack").(bool),
	}

	httpConfigObj, err := (*provider.Client).doCreateHttpInput(&input)
	if err != nil {
		return err
	}

	d.SetId(httpConfigObj.Name)
	err = d.Set("token", httpConfigObj.Token)
	if err != nil {
		return err
	}
	return nil
}

func httpEventCollectorInputDelete(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	input := SplunkHttpInputConfig{
		Name: d.Get("name").(string),
		AppContext: d.Get("app_context").(string),
	}

	err := (*provider.Client).doDeleteHttpInput(&input)
	if err != nil {
		return err
	}
	return nil
}

func httpEventCollectorInputUpdate(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	input := SplunkHttpInputConfig{
		Name: d.Get("name").(string),
		AppContext: d.Get("app_context").(string),
		Index:d.Get("index").(string),
		Source:d.Get("source").(string),
		Sourcetype:d.Get("sourcetype").(string),
		Disabled:d.Get("disabled").(bool),
		UseACK:d.Get("use_ack").(bool),
	}

	httpConfigObj, err := (*provider.Client).doUpdateHttpInput(&input)
	if err != nil {
		return err
	}

	d.SetId(httpConfigObj.Name)
	err = d.Set("token", httpConfigObj.Token)
	if err != nil {
		return err
	}
	return nil
}

// Calls

func (c *SplunkClient) doGetHttpInput(m *SplunkHttpInputConfig) (*SplunkHttpInputConfig, error) {
	httpInputEndpoint := splunkHttpInputEndpoint + m.Name
	httpInputEndpoint = fmt.Sprintf(httpInputEndpoint, m.AppContext)
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
		splunkHttpInputConfig, err := unmarshalSplunkHttpInputResponse(resp)
		if err != nil {
			return nil, err
		}
		return splunkHttpInputConfig, nil

	default:
		errorResponse := &Response{}
		_ = json.NewDecoder(resp.Body).Decode(errorResponse)
		err = errors.New(errorResponse.Messages[0].Text)
		return nil, err
	}
}

func (c *SplunkClient) doCreateHttpInput(m *SplunkHttpInputConfig) (*SplunkHttpInputConfig, error) {
	httpInputEndpoint := splunkHttpInputEndpoint
	httpInputEndpoint = fmt.Sprintf(httpInputEndpoint, m.AppContext)
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
		splunkHttpInputConfig, err := unmarshalSplunkHttpInputResponse(resp)
		if err != nil {
			return nil, err
		}

		return splunkHttpInputConfig, nil

	default:
		errorResponse := &Response{}
		_ = json.NewDecoder(resp.Body).Decode(errorResponse)
		err = errors.New(errorResponse.Messages[0].Text)
		return nil, err
	}
}

func (c *SplunkClient) doUpdateHttpInput(m *SplunkHttpInputConfig) (*SplunkHttpInputConfig, error) {
	httpInputEndpoint := splunkHttpInputEndpoint + m.Name
	httpInputEndpoint = fmt.Sprintf(httpInputEndpoint, m.AppContext)
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
		splunkHttpInputConfig, err := unmarshalSplunkHttpInputResponse(resp)
		if err != nil {
			return nil, err
		}

		return splunkHttpInputConfig, nil

	default:
		errorResponse := &Response{}
		_ = json.NewDecoder(resp.Body).Decode(errorResponse)
		err = errors.New(errorResponse.Messages[0].Text)
		return nil, err
	}

}

func (c *SplunkClient) doDeleteHttpInput(m *SplunkHttpInputConfig) error {
	httpInputEndpoint := splunkHttpInputEndpoint + m.Name
	httpInputEndpoint = fmt.Sprintf(httpInputEndpoint, m.AppContext)
	values := url.Values{}
	_, err := c.doRequest("DELETE", httpInputEndpoint, values)
	if err != nil {
		return err
	}
	return nil
}

func unmarshalSplunkHttpInputResponse(response *http.Response) (*SplunkHttpInputConfig, error) {
	splunkHttpInputConfig := &SplunkHttpInputConfig{}
	successResponse := &Response{}
	_ = json.NewDecoder(response.Body).Decode(&successResponse)
	re := regexp.MustCompile(`http://(.*)`)
	splunkHttpInputConfig.Name = re.FindStringSubmatch(successResponse.Entry[0].Name)[1]
	splunkHttpInputConfig.Token = successResponse.Entry[0].Content.Token
	return splunkHttpInputConfig, nil
}

