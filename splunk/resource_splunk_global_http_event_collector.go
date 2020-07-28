package splunk

import (
	"encoding/json"
	"errors"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"net/http"
	"terraform-provider-splunk/client/models"
)

type GlobalHttpInputConfig struct {
	Name      string
	Disabled  bool
	EnableSSL bool
	Port      int
}

func globalHttpEventCollector() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"disabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"port": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  8088,
			},
			"enable_ssl": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"dedicated_io_threads": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  2,
			},
			"max_sockets": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},
			"max_threads": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},
			"use_deployment_server": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
		Read:   globalHttpInputRead,
		Create: globalHttpInputCreate,
		Delete: globalHttpInputDelete,
		Update: globalHttpInputUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

// Functions
func globalHttpInputCreate(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	httpInputConfigObj := createGlobalHttpInputConfigObject(d)
	resp, err := (*provider.Client).CreateGlobalHttpEventCollectorObject(*httpInputConfigObj)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = unmarshalGlobalHttpInputResponse(resp)
	if err != nil {
		return err
	}

	d.SetId("http")
	return nil
}

func globalHttpInputRead(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	resp, err := (*provider.Client).ReadGlobalHttpEventCollectorObject()
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = unmarshalGlobalHttpInputResponse(resp)
	if err != nil {
		return err
	}

	d.SetId("http")
	return nil
}

func globalHttpInputDelete(d *schema.ResourceData, meta interface{}) error {
	// Global Http input resource object cannot be deleted
	return nil
}

func globalHttpInputUpdate(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	httpInputConfigObj := createGlobalHttpInputConfigObject(d)
	resp, err := (*provider.Client).CreateGlobalHttpEventCollectorObject(*httpInputConfigObj)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = unmarshalGlobalHttpInputResponse(resp)
	if err != nil {
		return err
	}

	d.SetId("http")
	return nil
}

// Helpers
func createGlobalHttpInputConfigObject(d *schema.ResourceData) (globalHttpInputConfigObject *models.GlobalHttpEventCollectorObject) {
	globalHttpInputConfigObject = &models.GlobalHttpEventCollectorObject{}
	globalHttpInputConfigObject.Disabled = d.Get("disabled").(bool)
	globalHttpInputConfigObject.Port = d.Get("port").(int)
	globalHttpInputConfigObject.EnableSSL = d.Get("enable_ssl").(bool)
	globalHttpInputConfigObject.DedicatedIoThreads = d.Get("dedicated_io_threads").(int)
	globalHttpInputConfigObject.MaxSockets = d.Get("max_sockets").(int)
	globalHttpInputConfigObject.MaxThreads = d.Get("max_threads").(int)
	globalHttpInputConfigObject.UseDeploymentServer = d.Get("use_deployment_server").(bool)
	return globalHttpInputConfigObject
}

func unmarshalGlobalHttpInputResponse(httpResponse *http.Response) (globalHttpEventCollectorObject *models.GlobalHttpEventCollectorObject, err error) {
	response := &models.GlobalHECResponse{}
	switch httpResponse.StatusCode {
	case 200, 201:
		_ = json.NewDecoder(httpResponse.Body).Decode(&response)
		return &response.Entry[0].Content, nil

	default:
		_ = json.NewDecoder(httpResponse.Body).Decode(response)
		err := errors.New(response.Messages[0].Text)
		return globalHttpEventCollectorObject, err
	}
}
