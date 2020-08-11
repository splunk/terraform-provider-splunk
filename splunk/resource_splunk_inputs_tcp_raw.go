package splunk

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"net/http"
	"regexp"
	"terraform-provider-splunk/client/models"
)

func inputsTCPRaw() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Required. The input port which receives raw data.",
			},
			"index": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Index to store generated events. Defaults to default.",
			},
			"host": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Host from which the indexer gets data. ",
			},
			"source": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: " 	Sets the source key/field for events from this input. Defaults to the input file path." +
					"Sets the source key initial value. The key is used during parsing/indexing, in particular to set the source field during indexing." +
					"It is also the source field used at search time. As a convenience, the chosen string is prepended with 'source::'." +
					"Note: Overriding the source key is generally not recommended. " +
					"Typically, the input layer provides a more accurate string to aid in problem analysis and investigation, " +
					"accurately recording the file from which the data was retrieved. " +
					"Consider use of source types, tagging, and search wildcards before overriding this value. ",
			},
			"sourcetype": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: "Set the source type for events from this input. \"sourcetype=\" is automatically prepended to <string>." +
					"Defaults to audittrail (if signedaudit=true) or fschange (if signedaudit=false).",
			},
			"disabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Indicates if input is disabled.",
			},
			"queue": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: "Valid values: (parsingQueue | indexQueue) " +
					"Specifies where the input processor should deposit the events it reads. Defaults to parsingQueue." +
					"Set queue to parsingQueue to apply props.conf and other parsing rules to your data. " +
					"Set queue to indexQueue to send your data directly into the index.",
			},
			"restrict_to_host": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "Allows for restricting this input to only accept data from the host specified here.",
			},
			"connection_host": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: "Valid values: (ip | dns | none)" +
					"Set the host for the remote server that is sending data." +
					"ip sets the host to the IP address of the remote server sending data." +
					"dns sets the host to the reverse DNS entry for the IP address of the remote server sending data." +
					"none leaves the host as specified in inputs.conf, which is typically the Splunk system hostname." +
					"Default value is dns. ",
			},
			"raw_tcp_done_timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  10,
				Description: " 	Specifies in seconds the timeout value for adding a Done-key. Default value is 10 seconds. " +
					"If a connection over the port specified by name remains idle after receiving data for specified number of seconds, it adds a Done-key. " +
					"This implies the last event is completely received. ",
			},
			"acl": aclSchema(),
		},
		Read:   inputsTCPRawRead,
		Create: inputsTCPRawCreate,
		Delete: inputsTCPRawDelete,
		Update: inputsTCPRawUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

// Functions
func inputsTCPRawCreate(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	name := d.Get("name").(string)
	inputsTCPRawConfig := getinputsTCPRawConfig(d)
	aclObject := &models.ACLObject{}
	if r, ok := d.GetOk("acl"); ok {
		aclObject = getACLConfig(r.([]interface{}))
	} else {
		aclObject.Owner = "nobody"
		aclObject.App = "search"
	}

	err := (*provider.Client).CreateTCPRawInput(name, aclObject.Owner, aclObject.App, inputsTCPRawConfig)
	if err != nil {
		return err
	}

	if _, ok := d.GetOk("acl"); ok {
		err = (*provider.Client).UpdateAcl(aclObject.Owner, aclObject.App, d.Id(), aclObject, "data", "inputs", "tcp", "raw")
		if err != nil {
			return err
		}
	}

	d.SetId(name)
	return inputsTCPRawRead(d, meta)
}

func inputsTCPRawRead(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	// We first get list of inputs to get owner and app name for the specific input
	resp, err := (*provider.Client).ReadTCPRawInputs()
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	entry, err := getinputsTCPRawConfigByName(d.Id(), resp)
	if err != nil {
		return err
	}

	if entry == nil {
		return errors.New(fmt.Sprintf("Unable to find resource: %v", d.Id()))
	}

	// Now we read the input configuration with proper owner and app
	resp, err = (*provider.Client).ReadTCPRawInput(entry.Name, entry.ACL.Owner, entry.ACL.App)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	entry, err = getinputsTCPRawConfigByName(d.Id(), resp)
	if err != nil {
		return err
	}

	if entry == nil {
		return errors.New(fmt.Sprintf("Unable to find resource: %v", d.Id()))
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

	if err = d.Set("raw_tcp_done_timeout", entry.Content.RawTcpDoneTime); err != nil {
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

func inputsTCPRawUpdate(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	inputsTCPRawConfig := getinputsTCPRawConfig(d)
	aclObject := getACLConfig(d.Get("acl").([]interface{}))
	name := d.Id()
	if r, ok := d.GetOk("restrict_to_host"); ok {
		name = r.(string) + ":" + name
	}
	err := (*provider.Client).UpdateTCPRawInput(name, aclObject.Owner, aclObject.App, inputsTCPRawConfig)
	if err != nil {
		return err
	}

	//ACL update
	err = (*provider.Client).UpdateAcl(aclObject.Owner, aclObject.App, name, aclObject, "data", "inputs", "tcp", "raw")
	if err != nil {
		return err
	}

	return inputsTCPRawRead(d, meta)
}

func inputsTCPRawDelete(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	name := d.Id()
	if r, ok := d.GetOk("restrict_to_host"); ok {
		name = r.(string) + ":" + name
	}
	aclObject := getACLConfig(d.Get("acl").([]interface{}))
	resp, err := (*provider.Client).DeleteTCPRawInput(name, aclObject.Owner, aclObject.App)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 200, 201:
		return nil

	default:
		errorResponse := &models.InputsTCPRawResponse{}
		_ = json.NewDecoder(resp.Body).Decode(errorResponse)
		err := errors.New(errorResponse.Messages[0].Text)
		return err
	}
}

// Helpers
func getinputsTCPRawConfig(d *schema.ResourceData) (inputsTCPRawObj *models.InputsTCPRawObject) {
	inputsTCPRawObj = &models.InputsTCPRawObject{}
	inputsTCPRawObj.Host = d.Get("host").(string)
	inputsTCPRawObj.Index = d.Get("index").(string)
	inputsTCPRawObj.Source = d.Get("source").(string)
	inputsTCPRawObj.SourceType = d.Get("sourcetype").(string)
	inputsTCPRawObj.Disabled = d.Get("disabled").(bool)
	inputsTCPRawObj.ConnectionHost = d.Get("connection_host").(string)
	inputsTCPRawObj.Queue = d.Get("queue").(string)
	inputsTCPRawObj.RestrictToHost = d.Get("restrict_to_host").(string)
	inputsTCPRawObj.RawTcpDoneTime = d.Get("raw_tcp_done_timeout").(int)
	return
}

func getinputsTCPRawConfigByName(name string, httpResponse *http.Response) (entry *models.InputsTCPRawEntry, err error) {
	response := &models.InputsTCPRawResponse{}
	switch httpResponse.StatusCode {
	case 200, 201:
		err = json.NewDecoder(httpResponse.Body).Decode(&response)
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
