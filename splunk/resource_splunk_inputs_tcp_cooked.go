package splunk

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jaware-splunk/terraform-provider-splunk/client/models"
	"net/http"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func inputsTCPCooked() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`\d+`), "Must be a Integer"),
				Description:  "Required. The port number of this input.",
			},
			"host": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Host from which the indexer gets data. ",
			},
			"disabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Indicates if input is disabled.",
			},
			"restrict_to_host": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "Allows for restricting this input to only accept data from the host specified here.",
			},
			"connection_host": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"ip", "dns", "none"}, false),
				Description: "Valid values: (ip | dns | none)" +
					"Set the host for the remote server that is sending data." +
					"ip sets the host to the IP address of the remote server sending data." +
					"dns sets the host to the reverse DNS entry for the IP address of the remote server sending data." +
					"none leaves the host as specified in inputs.conf, which is typically the Splunk system hostname." +
					"Default value is dns. ",
			},
			"acl": aclSchema(),
		},
		Read:   inputsTCPCookedRead,
		Create: inputsTCPCookedCreate,
		Delete: inputsTCPCookedDelete,
		Update: inputsTCPCookedUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

// Functions
func inputsTCPCookedCreate(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	name := d.Get("name").(string)
	inputsTCPCookedConfig := getinputsTCPCookedConfig(d)
	aclObject := &models.ACLObject{}
	if r, ok := d.GetOk("acl"); ok {
		aclObject = getACLConfig(r.([]interface{}))
	} else {
		aclObject.Owner = "nobody"
		aclObject.App = "search"
	}

	err := (*provider.Client).CreateTCPCookedInput(name, aclObject.Owner, aclObject.App, inputsTCPCookedConfig)
	if err != nil {
		return err
	}

	if _, ok := d.GetOk("acl"); ok {
		err = (*provider.Client).UpdateAcl(aclObject.Owner, aclObject.App, d.Id(), aclObject, "data", "inputs", "tcp", "cooked")
		if err != nil {
			return err
		}
	}

	d.SetId(name)
	return inputsTCPCookedRead(d, meta)
}

func inputsTCPCookedRead(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	// We first get list of inputs to get owner and app name for the specific input
	resp, err := (*provider.Client).ReadTCPCookedInputs()
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	entry, err := getinputsTCPCookedConfigByName(d.Id(), resp)
	if err != nil {
		return err
	}

	if entry == nil {
		return fmt.Errorf("Unable to find resource: %v", d.Id())
	}

	// Now we read the input configuration with proper owner and app
	resp, err = (*provider.Client).ReadTCPCookedInput(entry.Name, entry.ACL.Owner, entry.ACL.App)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	entry, err = getinputsTCPCookedConfigByName(d.Id(), resp)
	if err != nil {
		return err
	}

	if entry == nil {
		return fmt.Errorf("Unable to find resource: %v", d.Id())
	}

	if err = d.Set("name", d.Id()); err != nil {
		return err
	}

	if err = d.Set("host", entry.Content.Host); err != nil {
		return err
	}

	if err = d.Set("disabled", entry.Content.Disabled); err != nil {
		return err
	}

	if err = d.Set("restrict_to_host", entry.Content.RestrictToHost); err != nil {
		return err
	}

	if err = d.Set("connection_host", entry.Content.ConnectionHost); err != nil {
		return err
	}

	err = d.Set("acl", flattenACL(&entry.ACL))
	if err != nil {
		return err
	}

	return nil
}

func inputsTCPCookedUpdate(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	inputsTCPCookedConfig := getinputsTCPCookedConfig(d)
	aclObject := getACLConfig(d.Get("acl").([]interface{}))
	name := d.Id()
	if r, ok := d.GetOk("restrict_to_host"); ok {
		name = r.(string) + ":" + name
	}
	err := (*provider.Client).UpdateTCPCookedInput(name, aclObject.Owner, aclObject.App, inputsTCPCookedConfig)
	if err != nil {
		return err
	}

	//ACL update
	err = (*provider.Client).UpdateAcl(aclObject.Owner, aclObject.App, name, aclObject, "data", "inputs", "tcp", "cooked")
	if err != nil {
		return err
	}

	return inputsTCPCookedRead(d, meta)
}

func inputsTCPCookedDelete(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	name := d.Id()
	if r, ok := d.GetOk("restrict_to_host"); ok {
		name = r.(string) + ":" + name
	}
	aclObject := getACLConfig(d.Get("acl").([]interface{}))
	resp, err := (*provider.Client).DeleteTCPCookedInput(name, aclObject.Owner, aclObject.App)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 200, 201:
		return nil

	default:
		errorResponse := &models.InputsTCPCookedResponse{}
		_ = json.NewDecoder(resp.Body).Decode(errorResponse)
		err := errors.New(errorResponse.Messages[0].Text)
		return err
	}
}

// Helpers
func getinputsTCPCookedConfig(d *schema.ResourceData) (InputsTCPCookedObject *models.InputsTCPCookedObject) {
	InputsTCPCookedObject = &models.InputsTCPCookedObject{}
	InputsTCPCookedObject.Host = d.Get("host").(string)
	InputsTCPCookedObject.Disabled = d.Get("disabled").(bool)
	InputsTCPCookedObject.ConnectionHost = d.Get("connection_host").(string)
	InputsTCPCookedObject.RestrictToHost = d.Get("restrict_to_host").(string)
	return
}

func getinputsTCPCookedConfigByName(name string, httpResponse *http.Response) (entry *models.InputsTCPCookedEntry, err error) {
	response := &models.InputsTCPCookedResponse{}
	switch httpResponse.StatusCode {
	case 200, 201:
		err = json.NewDecoder(httpResponse.Body).Decode(&response)
		if err != nil {
			return nil, err
		}
		re := regexp.MustCompile(`(\d+)$`)
		for _, e := range response.Entry {
			if name == re.FindStringSubmatch(e.Name)[1] {
				return &e, nil
			}
		}

	default:
		_ = json.NewDecoder(httpResponse.Body).Decode(response)
		err := errors.New(response.Messages[0].Text)
		return entry, err
	}

	return entry, nil
}
