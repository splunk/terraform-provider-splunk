package splunk

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"github.com/terraform-providers/terraform-provider-splunk/client/models"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func inputsUDP() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`\d+`), "Must be a Integer"),
				Description:  "Required. The UDP port that this input should listen on.",
			},
			"index": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Which index events from this input should be stored in. Defaults to default.",
			},
			"host": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The value to populate in the host field for incoming events. This is used during parsing/indexing, in particular to set the host field. It is also the host field used at search time.",
			},
			"source": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The value to populate in the source field for incoming events. The same source should not be used for multiple data inputs.",
			},
			"sourcetype": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The value to populate in the sourcetype field for incoming events.",
			},
			"disabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Indicates if input is disabled.",
			},
			"queue": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Which queue events from this input should be sent to. Generally this does not need to be changed.",
			},
			"restrict_to_host": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Restrict incoming connections on this port to the host specified here. If this is not set, the value specified in [udp://<remote server>:<port>] in inputs.conf is used.",
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
					"Default value is ip. ",
			},
			"no_appending_timestamp": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "If set to true, prevents Splunk software from prepending a timestamp and hostname to incoming events.",
			},
			"no_priority_stripping": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "If set to true, Splunk software does not remove the priority field from incoming syslog events. ",
			},
			"acl": aclSchema(),
		},
		Read:   inputsUDPRead,
		Create: inputsUDPCreate,
		Delete: inputsUDPDelete,
		Update: inputsUDPUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

// Functions
func inputsUDPCreate(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	name := d.Get("name").(string)
	inputsUDPConfig := getInputsUDPConfig(d)
	aclObject := &models.ACLObject{}
	if r, ok := d.GetOk("acl"); ok {
		aclObject = getACLConfig(r.([]interface{}))
	} else {
		aclObject.Owner = "nobody"
		aclObject.App = "search"
	}

	err := (*provider.Client).CreateUDPInput(name, aclObject.Owner, aclObject.App, inputsUDPConfig)
	if err != nil {
		return err
	}

	if _, ok := d.GetOk("acl"); ok {
		err := (*provider.Client).UpdateAcl(aclObject.Owner, aclObject.App, name, aclObject, "data", "inputs", "udp")
		if err != nil {
			return err
		}
	}

	d.SetId(name)
	return inputsUDPRead(d, meta)
}

func inputsUDPRead(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	name := d.Id()
	// We first get list of inputs to get Scriptowner and app name for the specific input
	resp, err := (*provider.Client).ReadUDPInputs()
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	entry, err := getInputsUDPConfigByName(name, resp)
	if err != nil {
		return err
	}

	if entry == nil {
		return fmt.Errorf("Unable to find resource: %v", name)
	}

	// Now we read the input configuration with proper owner and app
	resp, err = (*provider.Client).ReadUDPInput(name, entry.ACL.Owner, entry.ACL.App)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	entry, err = getInputsUDPConfigByName(name, resp)
	if err != nil {
		return err
	}

	if entry == nil {
		return fmt.Errorf("Unable to find resource: %v", name)
	}

	if err = d.Set("name", d.Id()); err != nil {
		return err
	}

	if err = d.Set("host", entry.Content.Host); err != nil {
		return err
	}

	if err = d.Set("index", entry.Content.Index); err != nil {
		return err
	}

	if err = d.Set("sourcetype", entry.Content.SourceType); err != nil {
		return err
	}

	if err = d.Set("source", entry.Content.Source); err != nil {
		return err
	}

	if err = d.Set("disabled", entry.Content.Disabled); err != nil {
		return err
	}

	if err = d.Set("queue", entry.Content.Queue); err != nil {
		return err
	}

	if err = d.Set("restrict_to_host", entry.Content.RestrictToHost); err != nil {
		return err
	}

	if err = d.Set("no_appending_timestamp", entry.Content.NoAppendingTimestamp); err != nil {
		return err
	}

	if err = d.Set("connection_host", entry.Content.ConnectionHost); err != nil {
		return err
	}

	if err = d.Set("no_priority_stripping", entry.Content.NoPriorityStripping); err != nil {
		return err
	}

	err = d.Set("acl", flattenACL(&entry.ACL))
	if err != nil {
		return err
	}

	return nil
}

func inputsUDPUpdate(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	inputsUDPConfig := getInputsUDPConfig(d)
	aclObject := getACLConfig(d.Get("acl").([]interface{}))
	err := (*provider.Client).UpdateUDPInput(d.Id(), aclObject.Owner, aclObject.App, inputsUDPConfig)
	if err != nil {
		return err
	}

	//ACL update
	err = (*provider.Client).UpdateAcl(aclObject.Owner, aclObject.App, d.Id(), aclObject, "data", "inputs", "udp")
	if err != nil {
		return err
	}

	return inputsUDPRead(d, meta)
}

func inputsUDPDelete(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	aclObject := getACLConfig(d.Get("acl").([]interface{}))
	resp, err := (*provider.Client).DeleteUDPInput(d.Id(), aclObject.Owner, aclObject.App)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 200, 201:
		return nil

	default:
		errorResponse := &models.InputsUDPResponse{}
		_ = json.NewDecoder(resp.Body).Decode(errorResponse)
		err := errors.New(errorResponse.Messages[0].Text)
		return err
	}
}

// Helpers
func getInputsUDPConfig(d *schema.ResourceData) (inputsUDPObj *models.InputsUDPObject) {
	inputsUDPObj = &models.InputsUDPObject{}
	inputsUDPObj.Host = d.Get("host").(string)
	inputsUDPObj.Index = d.Get("index").(string)
	inputsUDPObj.Source = d.Get("source").(string)
	inputsUDPObj.SourceType = d.Get("sourcetype").(string)
	inputsUDPObj.Disabled = d.Get("disabled").(bool)
	inputsUDPObj.ConnectionHost = d.Get("connection_host").(string)
	inputsUDPObj.Queue = d.Get("queue").(string)
	inputsUDPObj.RestrictToHost = d.Get("restrict_to_host").(string)
	inputsUDPObj.NoAppendingTimestamp = d.Get("no_appending_timestamp").(bool)
	inputsUDPObj.NoPriorityStripping = d.Get("no_priority_stripping").(bool)
	return
}

func getInputsUDPConfigByName(name string, httpResponse *http.Response) (entry *models.InputsUDPEntry, err error) {
	response := &models.InputsUDPResponse{}
	switch httpResponse.StatusCode {
	case 200, 201:
		err = json.NewDecoder(httpResponse.Body).Decode(&response)
		if err != nil {
			return nil, err
		}
		re := regexp.MustCompile(`(.*)`)
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
