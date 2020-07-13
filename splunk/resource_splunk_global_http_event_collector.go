package splunk

import (
	"encoding/json"
	"errors"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"net/http"
	"net/url"
	"strconv"
)

type SplunkGlobalHttpInputConfig struct {
	Name        string
	Disabled    bool
	EnableSSL   bool
	Port        int
}


func globalHttpEventCollector() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:     true,
				Default:     "http",
			},
			"disabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"port": {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     8088,
			},
			"enable_ssl": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
		},
		Read:globalHttpInputRead,
		Create:globalHttpInputCreate,
		Delete:globalHttpInputDelete,
		Update:globalHttpInputUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func globalHttpInputRead(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	input := SplunkGlobalHttpInputConfig{
		Name: d.Get("name").(string),
	}

	httpConfigObj, err := (*provider.Client).doGetGlobalHttpInput(&input)
	if err != nil {
		return err
	}

	d.SetId(httpConfigObj.Name)
	if err != nil {
		return err
	}
	return nil
}

func globalHttpInputCreate(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	input := SplunkGlobalHttpInputConfig{
		Name: d.Get("name").(string),
		Disabled:d.Get("disabled").(bool),
		EnableSSL:d.Get("enable_ssl").(bool),
		Port:d.Get("port").(int),
	}

	httpConfigObj, err := (*provider.Client).doCreateGlobalHttpInput(&input)
	if err != nil {
		return err
	}

	d.SetId(httpConfigObj.Name)
	if err != nil {
		return err
	}
	return nil
}

func globalHttpInputDelete(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	input := SplunkGlobalHttpInputConfig{
		Name: d.Get("name").(string),
	}

	err := (*provider.Client).doGlobalDeleteHttpInput(&input)
	if err != nil {
		return err
	}
	return nil
}

func globalHttpInputUpdate(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	input := SplunkGlobalHttpInputConfig{
		Name: d.Get("name").(string),
		Disabled:d.Get("disabled").(bool),
		EnableSSL:d.Get("enable_ssl").(bool),
		Port:d.Get("port").(int),
	}

	httpConfigObj, err := (*provider.Client).doGlobalUpdateHttpInput(&input)
	if err != nil {
		return err
	}

	d.SetId(httpConfigObj.Name)
	if err != nil {
		return err
	}
	return nil
}

// Calls

func (c *SplunkClient) doGetGlobalHttpInput(m *SplunkGlobalHttpInputConfig) (*SplunkGlobalHttpInputConfig, error) {
	httpInputEndpoint := splunkGlobalHttpInputEndpoint + m.Name
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
		splunkGlobalHttpInputConfig, err := unmarshalSplunkGlobalHttpInputResponse(resp)
		if err != nil {
			return nil, err
		}
		return splunkGlobalHttpInputConfig, nil

	default:
		errorResponse := &Response{}
		_ = json.NewDecoder(resp.Body).Decode(errorResponse)
		err = errors.New(errorResponse.Messages[0].Text)
		return nil, err
	}
}

func (c *SplunkClient) doCreateGlobalHttpInput(m *SplunkGlobalHttpInputConfig) (*SplunkGlobalHttpInputConfig, error) {
	httpInputEndpoint := splunkGlobalHttpInputEndpoint + m.Name
	values := url.Values{}
	values.Add("disabled", strconv.FormatBool(m.Disabled))
	values.Add("enableSSL", strconv.FormatBool(m.EnableSSL))
	values.Add("port", strconv.Itoa(m.Port))
	resp, err := c.doRequest("POST", httpInputEndpoint, values)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	switch resp.StatusCode {
	case 200, 201:
		splunkGlobalHttpInputConfig, err := unmarshalSplunkGlobalHttpInputResponse(resp)
		if err != nil {
			return nil, err
		}

		return splunkGlobalHttpInputConfig, nil

	default:
		errorResponse := &Response{}
		_ = json.NewDecoder(resp.Body).Decode(errorResponse)
		err = errors.New(errorResponse.Messages[0].Text)
		return nil, err
	}
}

func (c *SplunkClient) doGlobalUpdateHttpInput(m *SplunkGlobalHttpInputConfig) (*SplunkGlobalHttpInputConfig, error) {
	httpInputEndpoint := splunkGlobalHttpInputEndpoint + m.Name
	values := url.Values{}
	values.Add("disabled", strconv.FormatBool(m.Disabled))
	values.Add("enableSSL", strconv.FormatBool(m.EnableSSL))
	values.Add("port", strconv.Itoa(m.Port))
	resp, err := c.doRequest("POST", httpInputEndpoint, values)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	switch resp.StatusCode {
	case 200, 201:
		splunkGlobalHttpInputConfig, err := unmarshalSplunkGlobalHttpInputResponse(resp)
		if err != nil {
			return nil, err
		}

		return splunkGlobalHttpInputConfig, nil

	default:
		errorResponse := &Response{}
		_ = json.NewDecoder(resp.Body).Decode(errorResponse)
		err = errors.New(errorResponse.Messages[0].Text)
		return nil, err
	}

}

func (c *SplunkClient) doGlobalDeleteHttpInput(m *SplunkGlobalHttpInputConfig) error {
	httpInputEndpoint := splunkGlobalHttpInputEndpoint + m.Name
	values := url.Values{}
	_, err := c.doRequest("DELETE", httpInputEndpoint, values)
	if err != nil {
		return err
	}
	return nil
}

func unmarshalSplunkGlobalHttpInputResponse(response *http.Response) (*SplunkGlobalHttpInputConfig, error) {
	splunkGlobalHttpInputConfig := &SplunkGlobalHttpInputConfig{}
	successResponse := &Response{}
	_ = json.NewDecoder(response.Body).Decode(&successResponse)
	splunkGlobalHttpInputConfig.Name = successResponse.Entry[0].Name
	return splunkGlobalHttpInputConfig, nil
}



