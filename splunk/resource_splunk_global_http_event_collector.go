package splunk

import (
	"encoding/json"
	"errors"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"net/http"
	"terraform-provider-splunk/client/models"
)

func globalHttpEventCollector() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"disabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Input disabled indicator: 0 = Input Not disabled, 1 = Input disabled.",
			},
			"port": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(1024, 65536),
				Description:  "HTTP data input IP port.",
			},
			"enable_ssl": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Enable SSL protocol for HTTP data input. 1 = SSL enabled, 0 = SSL disabled.",
			},
			"dedicated_io_threads": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Number of threads used by HTTP Input server.",
			},
			"max_sockets": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				Description: "Maximum number of simultaneous HTTP connections accepted. " +
					"Adjusting this value may cause server performance issues and is not generally recommended. " +
					"Possible values for this setting vary by OS.",
			},
			"max_threads": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				Description: "Maximum number of threads that can be used by active HTTP transactions. " +
					"Adjusting this value may cause server performance issues and is not generally recommended. " +
					"Possible values for this setting vary by OS.",
			},
			"use_deployment_server": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
				Description: "Indicates whether the event collector input writes its configuration to a deployment server repository." +
					"When this setting is set to 1 (enabled), the input writes its configuration to the directory specified as repositoryLocation in serverclass.conf." +
					"Copy the full contents of the splunk_httpinput app directory to this directory for the configuration to work." +
					"When enabled, only the tokens defined in the splunk_httpinput app in this repository are viewable and editable on the API and the Data Inputs page in Splunk Web." +
					"When disabled, the input writes its configuration to $SPLUNK_HOME/etc/apps by default." +
					"Defaults to 0 (disabled). ",
			},
		},
		Read:   globalHttpInputRead,
		Create: globalHttpInputCreate,
		Update: globalHttpInputUpdate,
		Delete: globalHttpInputDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

// Functions
func globalHttpInputCreate(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	httpInputConfigObj := createGlobalHttpInputConfigObject(d)
	err := (*provider.Client).CreateGlobalHttpEventCollectorObject(*httpInputConfigObj)
	if err != nil {
		return err
	}

	d.SetId("http")
	return globalHttpInputRead(d, meta)
}

func globalHttpInputRead(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	resp, err := (*provider.Client).ReadGlobalHttpEventCollectorObject()
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	content, err := unmarshalGlobalHttpInputResponse(resp)
	if err != nil {
		return err
	}

	if err = d.Set("port", content.Port); err != nil {
		return err
	}

	if err = d.Set("dedicated_io_threads", content.DedicatedIoThreads); err != nil {
		return err
	}

	if err = d.Set("max_sockets", content.MaxSockets); err != nil {
		return err
	}

	if err = d.Set("max_threads", content.MaxThreads); err != nil {
		return err
	}

	if err = d.Set("disabled", content.Disabled); err != nil {
		return err
	}

	if err = d.Set("enable_ssl", content.EnableSSL); err != nil {
		return err
	}

	if err = d.Set("use_deployment_server", content.UseDeploymentServer); err != nil {
		return err
	}

	return nil
}

func globalHttpInputUpdate(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	httpInputConfigObj := createGlobalHttpInputConfigObject(d)
	err := (*provider.Client).UpdateGlobalHttpEventCollectorObject(*httpInputConfigObj)
	if err != nil {
		return err
	}

	return globalHttpInputRead(d, meta)
}

func globalHttpInputDelete(d *schema.ResourceData, meta interface{}) error {
	// Global Http input resource object cannot be deleted
	return nil
}

// Helpers
func createGlobalHttpInputConfigObject(d *schema.ResourceData) (globalHttpInputConfigObject *models.GlobalHttpEventCollectorObject) {
	globalHttpInputConfigObject = &models.GlobalHttpEventCollectorObject{}
	globalHttpInputConfigObject.Disabled = d.Get("disabled").(bool)
	globalHttpInputConfigObject.EnableSSL = d.Get("enable_ssl").(bool)
	globalHttpInputConfigObject.Port = d.Get("port").(int)
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
