package splunk

import (
	"encoding/json"
	"errors"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"net/http"
	"net/url"
	"strconv"
)

type GlobalHttpInputConfig struct {
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
	name := d.Get("name").(string)
	values := url.Values{}
	values.Add("name", name)
	endpoint := (*provider.Client).BuildSplunkdURL(nil, "services", "data", "inputs", "http", name)
	resp, err := (*provider.Client).Get(endpoint)
	defer resp.Body.Close()
	if err != nil {
		return err
	}

	globalHttpInputConfig, err := unmarshalGlobalHttpInputResponse(resp)
	if err != nil {
		return err
	}

	d.SetId(globalHttpInputConfig.Name)
	return nil
}

func globalHttpInputCreate(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	name := d.Get("name").(string)
	values := url.Values{}
	values.Add("disabled", strconv.FormatBool(d.Get("disabled").(bool)))
	values.Add("enableSSL", strconv.FormatBool(d.Get("enable_ssl").(bool)))
	values.Add("port", strconv.Itoa(d.Get("port").(int)))
	endpoint := (*provider.Client).BuildSplunkdURL(nil, "services", "data", "inputs", "http", name)
	resp, err := (*provider.Client).Post(endpoint, values)
	defer resp.Body.Close()
	if err != nil {
		return err
	}

	globalHttpInputConfig, err := unmarshalGlobalHttpInputResponse(resp)
	if err != nil {
		return err
	}

	d.SetId(globalHttpInputConfig.Name)
	return nil
}

func globalHttpInputDelete(d *schema.ResourceData, meta interface{}) error {
	// Global Http input resource object cannot be deleted
	return nil
}

func globalHttpInputUpdate(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	name := d.Get("name").(string)
	values := url.Values{}
	values.Add("disabled", strconv.FormatBool(d.Get("disabled").(bool)))
	values.Add("enableSSL", strconv.FormatBool(d.Get("enable_ssl").(bool)))
	values.Add("port", strconv.Itoa(d.Get("port").(int)))
	endpoint := (*provider.Client).BuildSplunkdURL(nil, "services", "data", "inputs", "http", name)
	resp, err := (*provider.Client).Post(endpoint, values)
	defer resp.Body.Close()
	if err != nil {
		return err
	}

	globalHttpInputConfig, err := unmarshalGlobalHttpInputResponse(resp)
	if err != nil {
		return err
	}

	d.SetId(globalHttpInputConfig.Name)
	return nil
}

func unmarshalGlobalHttpInputResponse(response *http.Response) (*GlobalHttpInputConfig, error) {
	switch response.StatusCode {
	case 200, 201:
		globalHttpInputConfig := &GlobalHttpInputConfig{}
		successResponse := &Response{}
		_ = json.NewDecoder(response.Body).Decode(&successResponse)
		globalHttpInputConfig.Name = successResponse.Entry[0].Name
		return globalHttpInputConfig, nil

	default:
		errorResponse := &Response{}
		_ = json.NewDecoder(response.Body).Decode(errorResponse)
		err := errors.New(errorResponse.Messages[0].Text)
		return nil, err
	}
}



