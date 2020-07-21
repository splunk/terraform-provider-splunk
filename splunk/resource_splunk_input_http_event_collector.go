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

type HttpInputConfig struct {
	Name        string
	AppContext  string
	Index       string
	Source      string
	SourceType  string
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
	name := d.Get("name").(string)
	appContext := d.Get("app_context").(string)
	values := url.Values{}
	values.Add("name", name)
	endpoint := (*provider.Client).BuildSplunkdURL(nil, "servicesNS", "nobody", appContext, "data", "inputs", "http", name)
	resp, err := (*provider.Client).Get(endpoint)
	defer resp.Body.Close()
	if err != nil {
		return err
	}

	httpInputConfig, err := unmarshalHttpInputResponse(resp)
	if err != nil {
		return err
	}

	d.SetId(httpInputConfig.Name)
	err = d.Set("token", httpInputConfig.Token)
	if err != nil {
		return err
	}

	return nil
}

func httpEventCollectorInputCreate(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	name := d.Get("name").(string)
	appContext := d.Get("app_context").(string)
	values := url.Values{}
	values.Add("name", name)
	values.Add("index", d.Get("index").(string))
	values.Add("source", d.Get("source").(string))
	values.Add("sourcetype", d.Get("sourcetype").(string))
	values.Add("useACK", strconv.FormatBool(d.Get("use_ack").(bool)))
	values.Add("disabled", strconv.FormatBool(d.Get("disabled").(bool)))
	endpoint := (*provider.Client).BuildSplunkdURL(nil, "servicesNS", "nobody", appContext, "data", "inputs", "http", name)
	resp, err := (*provider.Client).Post(endpoint, values)
	defer resp.Body.Close()
	if err != nil {
		return err
	}

	httpInputConfig, err := unmarshalHttpInputResponse(resp)
	if err != nil {
		return err
	}

	d.SetId(httpInputConfig.Name)
	err = d.Set("token", httpInputConfig.Token)
	if err != nil {
		return err
	}

	return nil
}

func httpEventCollectorInputDelete(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	name := d.Get("name").(string)
	appContext := d.Get("app_context").(string)
	endpoint := (*provider.Client).BuildSplunkdURL(nil, "servicesNS", "nobody", appContext, "data", "inputs", "http", name)
	resp, err := (*provider.Client).Delete(endpoint)
	defer resp.Body.Close()
	if err != nil {
		return err
	}

	switch resp.StatusCode {
	case 200, 201:
		return nil

	default:
		errorResponse := &Response{}
		_ = json.NewDecoder(resp.Body).Decode(errorResponse)
		err := errors.New(errorResponse.Messages[0].Text)
		return err
	}
}

func httpEventCollectorInputUpdate(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	name := d.Get("name").(string)
	appContext := d.Get("app_context").(string)
	values := url.Values{}
	values.Add("name", name)
	values.Add("index", d.Get("index").(string))
	values.Add("source", d.Get("source").(string))
	values.Add("sourcetype", d.Get("sourcetype").(string))
	values.Add("useACK", strconv.FormatBool(d.Get("use_ack").(bool)))
	values.Add("disabled", strconv.FormatBool(d.Get("disabled").(bool)))
	endpoint := (*provider.Client).BuildSplunkdURL(nil, "servicesNS", "nobody", appContext, "data", "inputs", "http", name)
	resp, err := (*provider.Client).Post(endpoint, values)
	defer resp.Body.Close()
	if err != nil {
		return err
	}

	httpInputConfig, err := unmarshalHttpInputResponse(resp)
	if err != nil {
		return err
	}

	d.SetId(httpInputConfig.Name)
	err = d.Set("token", httpInputConfig.Token)
	if err != nil {
		return err
	}

	return nil
}

func unmarshalHttpInputResponse(response *http.Response) (*HttpInputConfig, error) {
	switch response.StatusCode {
	case 200, 201:
		httpInputConfig := &HttpInputConfig{}
		successResponse := &Response{}
		_ = json.NewDecoder(response.Body).Decode(&successResponse)
		re := regexp.MustCompile(`http://(.*)`)
		httpInputConfig.Name = re.FindStringSubmatch(successResponse.Entry[0].Name)[1]
		httpInputConfig.Token = successResponse.Entry[0].Content.Token
		return httpInputConfig, nil

	default:
		errorResponse := &Response{}
		_ = json.NewDecoder(response.Body).Decode(errorResponse)
		err := errors.New(errorResponse.Messages[0].Text)
		return nil, err
	}
}

