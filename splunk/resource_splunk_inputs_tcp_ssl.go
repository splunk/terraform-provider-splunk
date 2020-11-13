package splunk

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"github.com/terraform-providers/terraform-provider-splunk/client/models"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func inputsTCPSSL() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Sensitive:   true,
				Description: "Server certificate password, if any.",
			},
			"root_ca": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Certificate authority list (root file)",
			},
			"server_cert": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Full path to the server certificate.",
			},
			"disabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Indicates if input is disabled.",
			},
			"require_client_cert": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Determines whether a client must authenticate.",
			},
		},
		Read:   inputsTCPSSLRead,
		Create: inputsTCPSSLCreate,
		Delete: inputsTCPSSLDelete,
		Update: inputsTCPSSLUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

// Functions
func inputsTCPSSLCreate(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	inputsTCPSSLConfig := getInputsTCPSSLConfig(d)
	err := (*provider.Client).CreateTCPSSLInput(inputsTCPSSLConfig)
	if err != nil {
		return err
	}

	d.SetId("ssl")
	return inputsTCPSSLRead(d, meta)
}

func inputsTCPSSLRead(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	name := d.Id()
	resp, err := (*provider.Client).ReadTCPSSLInput()
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	content, err := getInputsTCPSSLConfigByName(name, resp)
	if err != nil {
		return err
	}

	if content == nil {
		return fmt.Errorf("Unable to find resource: %v", d.Id())
	}

	if err = d.Set("disabled", content.Disabled); err != nil {
		return err
	}

	if err = d.Set("require_client_cert", content.RequireClientCert); err != nil {
		return err
	}

	if err = d.Set("root_ca", content.RootCA); err != nil {
		return err
	}

	if err = d.Set("password", content.Password); err != nil {
		return err
	}

	if err = d.Set("server_cert", content.ServerCert); err != nil {
		return err
	}

	return nil
}

func inputsTCPSSLUpdate(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	inputsTCPSSLConfig := getInputsTCPSSLConfig(d)
	name := d.Id()
	err := (*provider.Client).UpdateTCPSSLInput(name, inputsTCPSSLConfig)
	if err != nil {
		return err
	}

	return inputsTCPSSLRead(d, meta)
}

func inputsTCPSSLDelete(d *schema.ResourceData, meta interface{}) error {
	// SSL object cannot be removed
	return nil
}

// Helpers
func getInputsTCPSSLConfig(d *schema.ResourceData) (inputsTCPSSLObj *models.InputsTCPSSLObject) {
	inputsTCPSSLObj = &models.InputsTCPSSLObject{}
	inputsTCPSSLObj.Disabled = d.Get("disabled").(bool)
	inputsTCPSSLObj.Password = d.Get("password").(string)
	inputsTCPSSLObj.RootCA = d.Get("root_ca").(string)
	inputsTCPSSLObj.ServerCert = d.Get("server_cert").(string)
	inputsTCPSSLObj.RequireClientCert = d.Get("require_client_cert").(bool)
	return
}

func getInputsTCPSSLConfigByName(name string, httpResponse *http.Response) (inputsTCPSSLObject *models.InputsTCPSSLObject, err error) {
	response := &models.InputsTCPSSLResponse{}
	switch httpResponse.StatusCode {
	case 200, 201:
		err = json.NewDecoder(httpResponse.Body).Decode(&response)
		if err != nil {
			return nil, err
		}
		return &response.Entry[0].Content, nil

	default:
		_ = json.NewDecoder(httpResponse.Body).Decode(response)
		err := errors.New(response.Messages[0].Text)
		return inputsTCPSSLObject, err
	}
}
