package splunk

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"terraform-provider-splunk/client/models"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func inputsMonitor() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The file or directory path to monitor on the system.",
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
				Description: "The value to populate in the host field for events from this data input.",
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
				Description: "Indicates if input monitoring is disabled.",
			},
			"rename_source": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The value to populate in the source field for events from this data input. The same source should not be used for multiple data inputs.",
			},
			"blacklist": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specify a regular expression for a file path. The file path that matches this regular expression is not indexed.",
			},
			"whitelist": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specify a regular expression for a file path. Only file paths that match this regular expression are indexed.",
			},
			"crc_salt": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "A string that modifies the file tracking identity for files in this input. The magic value <SOURCE> invokes special behavior (see admin documentation).",
			},
			"follow_tail": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "If set to true, files that are seen for the first time is read from the end.",
			},
			"recursive": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Setting this to false prevents monitoring of any subdirectories encountered within this data input.",
			},
			"host_regex": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specify a regular expression for a file path. If the path for a file matches this regular expression, the captured value is used to populate the host field for events from this data input. The regular expression must have one capture group.",
			},
			"host_segment": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Use the specified slash-separate segment of the filepath as the host field value.",
			},
			"ignore_older_than": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specify a time value. If the modification time of a file being monitored falls outside of this rolling time window, the file is no longer being monitored.",
			},
			"time_before_close": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "When Splunk software reaches the end of a file that is being read, the file is kept open for a minimum of the number of seconds specified in this value. After this period has elapsed, the file is checked again for more data.",
			},
			"acl": aclSchema(),
		},
		Read:   inputsMonitorRead,
		Create: inputsMonitorCreate,
		Delete: inputsMonitorDelete,
		Update: inputsMonitorUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

// Functions
func inputsMonitorCreate(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	name := d.Get("name").(string)
	inputsMonitorConfig := getInputsMonitorConfig(d)
	aclObject := &models.ACLObject{}
	if r, ok := d.GetOk("acl"); ok {
		aclObject = getACLConfig(r.([]interface{}))
	} else {
		aclObject.Owner = "nobody"
		aclObject.App = "search"
	}

	err := (*provider.Client).CreateMonitorInput(name, aclObject.Owner, aclObject.App, inputsMonitorConfig)
	if err != nil {
		return err
	}

	if _, ok := d.GetOk("acl"); ok {
		err := (*provider.Client).UpdateAcl(aclObject.Owner, aclObject.App, url.PathEscape(name), aclObject, "data", "inputs", "monitor")
		if err != nil {
			return err
		}
	}

	d.SetId(name)
	return inputsMonitorRead(d, meta)
}

func inputsMonitorRead(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	name := d.Id()
	// We first get list of inputs to get Scriptowner and app name for the specific input
	resp, err := (*provider.Client).ReadMonitorInputs()
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	entry, err := getInputsMonitorConfigByName(name, resp)
	if err != nil {
		return err
	}

	if entry == nil {
		return errors.New(fmt.Sprintf("Unable to find resource: %v", name))
	}

	// Now we read the input configuration with proper owner and app
	resp, err = (*provider.Client).ReadMonitorInput(name, entry.ACL.Owner, entry.ACL.App)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	entry, err = getInputsMonitorConfigByName(name, resp)
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

	if err = d.Set("sourcetype", entry.Content.SourceType); err != nil {
		return err
	}

	if err = d.Set("disabled", entry.Content.Disabled); err != nil {
		return err
	}

	if err = d.Set("crc_salt", entry.Content.CrcSalt); err != nil {
		return err
	}

	if err = d.Set("follow_tail", entry.Content.FollowTail); err != nil {
		return err
	}

	if err = d.Set("recursive", entry.Content.Recursive); err != nil {
		return err
	}

	if err = d.Set("host_regex", entry.Content.HostRegex); err != nil {
		return err
	}

	if err = d.Set("host_segment", entry.Content.HostSegment); err != nil {
		return err
	}

	if err = d.Set("time_before_close", entry.Content.TimeBeforeClose); err != nil {
		return err
	}

	if err = d.Set("ignore_older_than", entry.Content.IgnoreOlderThan); err != nil {
		return err
	}

	if err = d.Set("blacklist", entry.Content.Blacklist); err != nil {
		return err
	}

	if err = d.Set("whitelist", entry.Content.Whitelist); err != nil {
		return err
	}

	err = d.Set("acl", flattenACL(&entry.ACL))
	if err != nil {
		return err
	}

	return nil
}

func inputsMonitorUpdate(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	inputsMonitorConfig := getInputsMonitorConfig(d)
	aclObject := getACLConfig(d.Get("acl").([]interface{}))
	err := (*provider.Client).UpdateMonitorInput(d.Id(), aclObject.Owner, aclObject.App, inputsMonitorConfig)
	if err != nil {
		return err
	}

	//ACL update
	err = (*provider.Client).UpdateAcl(aclObject.Owner, aclObject.App, url.PathEscape(d.Id()), aclObject, "data", "inputs", "monitor")
	if err != nil {
		return err
	}

	return inputsMonitorRead(d, meta)
}

func inputsMonitorDelete(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	aclObject := getACLConfig(d.Get("acl").([]interface{}))
	resp, err := (*provider.Client).DeleteMonitorInput(d.Id(), aclObject.Owner, aclObject.App)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 200, 201:
		return nil

	default:
		errorResponse := &models.InputsMonitorResponse{}
		_ = json.NewDecoder(resp.Body).Decode(errorResponse)
		err := errors.New(errorResponse.Messages[0].Text)
		return err
	}
}

// Helpers
func getInputsMonitorConfig(d *schema.ResourceData) (inputsMonitorObj *models.InputsMonitorObject) {
	inputsMonitorObj = &models.InputsMonitorObject{}
	inputsMonitorObj.Host = d.Get("host").(string)
	inputsMonitorObj.Index = d.Get("index").(string)
	inputsMonitorObj.SourceType = d.Get("sourcetype").(string)
	inputsMonitorObj.RenameSource = d.Get("rename_source").(string)
	inputsMonitorObj.Disabled = d.Get("disabled").(bool)
	inputsMonitorObj.CrcSalt = d.Get("crc_salt").(string)
	inputsMonitorObj.FollowTail = d.Get("follow_tail").(bool)
	inputsMonitorObj.Recursive = d.Get("recursive").(bool)
	inputsMonitorObj.HostRegex = d.Get("host_regex").(string)
	inputsMonitorObj.HostSegment = d.Get("host_segment").(int)
	inputsMonitorObj.TimeBeforeClose = d.Get("time_before_close").(int)
	inputsMonitorObj.IgnoreOlderThan = d.Get("ignore_older_than").(string)
	inputsMonitorObj.Blacklist = d.Get("blacklist").(string)
	inputsMonitorObj.Whitelist = d.Get("whitelist").(string)
	return
}

func getInputsMonitorConfigByName(name string, httpResponse *http.Response) (entry *models.InputsMonitorEntry, err error) {
	response := &models.InputsMonitorResponse{}
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
