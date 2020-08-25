package splunk

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"net/http"
	"regexp"
	"terraform-provider-splunk/client/models"
)

func inputsScript() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Specify the name of the scripted input.",
			},
			"index": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Sets the index for events from this input. Defaults to the main index.",
			},
			"host": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Sets the host for events from this input. Defaults to whatever host sent the event.",
			},
			"source": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Sets the source key/field for events from this input. Defaults to the input file path.",
			},
			"sourcetype": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Sets the sourcetype key/field for events from this input.",
			},
			"rename_source": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specify a new name for the source field for the script.",
			},
			"passauth": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "User to run the script as. If you provide a username, Splunk software generates an auth token for that user and passes it to the script.",
			},
			"disabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Specifies whether the input script is disabled.",
			},
			"interval": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntAtLeast(60),
				Description:  "Specify an integer or cron schedule. This parameter specifies how often to execute the specified script, in seconds or a valid cron schedule. If you specify a cron schedule, the script is not executed on start-up.",
			},
			"acl": aclSchema(),
		},
		Read:   inputsScriptRead,
		Create: inputsScriptCreate,
		Delete: inputsScriptDelete,
		Update: inputsScriptUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

// Functions
func inputsScriptCreate(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	name := d.Get("name").(string)
	inputsScriptConfig := getInputsScriptConfig(d)
	aclObject := &models.ACLObject{}
	if r, ok := d.GetOk("acl"); ok {
		aclObject = getACLConfig(r.([]interface{}))
	} else {
		aclObject.Owner = "nobody"
		aclObject.App = "search"
	}

	err := (*provider.Client).CreateScriptedInput(name, aclObject.Owner, aclObject.App, inputsScriptConfig)
	if err != nil {
		return err
	}

	if _, ok := d.GetOk("acl"); ok {
		err := (*provider.Client).UpdateAcl(aclObject.Owner, aclObject.App, name, aclObject, "data", "inputs", "script")
		if err != nil {
			return err
		}
	}

	d.SetId(name)
	return inputsScriptRead(d, meta)
}

func inputsScriptRead(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	name := d.Id()
	// We first get list of inputs to get owner and app name for the specific input
	resp, err := (*provider.Client).ReadScriptedInputs()
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	entry, err := getScriptedInputsConfigByName(name, resp)
	if err != nil {
		return err
	}

	if entry == nil {
		return errors.New(fmt.Sprintf("Unable to find resource: %v", name))
	}

	// Now we read the input configuration with proper owner and app
	resp, err = (*provider.Client).ReadScriptedInput(name, entry.ACL.Owner, entry.ACL.App)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	entry, err = getScriptedInputsConfigByName(name, resp)
	if err != nil {
		return err
	}

	if entry == nil {
		return errors.New(fmt.Sprintf("Unable to find resource: %v", name))
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

	if err = d.Set("source", entry.Content.Source); err != nil {
		return err
	}

	if err = d.Set("sourcetype", entry.Content.SourceType); err != nil {
		return err
	}

	if err = d.Set("disabled", entry.Content.Disabled); err != nil {
		return err
	}

	if err = d.Set("interval", entry.Content.Interval); err != nil {
		return err
	}

	err = d.Set("acl", flattenACL(&entry.ACL))
	if err != nil {
		return err
	}

	return nil
}

func inputsScriptUpdate(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	inputsScriptConfig := getInputsScriptConfig(d)
	aclObject := getACLConfig(d.Get("acl").([]interface{}))
	err := (*provider.Client).UpdateScriptedInput(d.Id(), aclObject.Owner, aclObject.App, inputsScriptConfig)
	if err != nil {
		return err
	}

	//ACL update
	err = (*provider.Client).UpdateAcl(aclObject.Owner, aclObject.App, d.Id(), aclObject, "data", "inputs", "script")
	if err != nil {
		return err
	}

	return inputsScriptRead(d, meta)
}

func inputsScriptDelete(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	aclObject := getACLConfig(d.Get("acl").([]interface{}))
	resp, err := (*provider.Client).DeleteScriptedInput(d.Id(), aclObject.Owner, aclObject.App)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 200, 201:
		return nil

	default:
		errorResponse := &models.HECResponse{}
		_ = json.NewDecoder(resp.Body).Decode(errorResponse)
		err := errors.New(errorResponse.Messages[0].Text)
		return err
	}
}

// Helpers
func getInputsScriptConfig(d *schema.ResourceData) (inputsScriptObj *models.InputsScriptObject) {
	inputsScriptObj = &models.InputsScriptObject{}
	inputsScriptObj.Host = d.Get("host").(string)
	inputsScriptObj.Index = d.Get("index").(string)
	inputsScriptObj.Source = d.Get("source").(string)
	inputsScriptObj.SourceType = d.Get("sourcetype").(string)
	inputsScriptObj.RenameSource = d.Get("rename_source").(string)
	inputsScriptObj.PassAuth = d.Get("passauth").(string)
	inputsScriptObj.Disabled = d.Get("disabled").(bool)
	inputsScriptObj.Interval = d.Get("interval").(int)
	return
}

func getScriptedInputsConfigByName(name string, httpResponse *http.Response) (entry *models.InputsScriptEntry, err error) {
	response := &models.InputsScriptResponse{}
	switch httpResponse.StatusCode {
	case 200, 201:
		err = json.NewDecoder(httpResponse.Body).Decode(&response)
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
