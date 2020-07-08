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

func outputsTCPSyslog() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"disabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Disables default tcpout settings",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the syslog output group. This is name used when creating syslog configuration in outputs.conf.",
			},
			"priority": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				Description: "Sets syslog priority value.The priority value should specified as an integer. " +
					"See $SPLUNK_HOME/etc/system/README/outputs.conf.spec for details. ",
			},
			"server": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringMatch(validTCPSyslogServer, "<host>:<port> of the Splunk receiver"),
				Description:  "host:port of the server where syslog data should be sent ",
			},
			"syslog_sourcetype": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: "Specifies a rule for handling data in addition to that provided by the \"syslog\" sourcetype. " +
					"By default, there is no value for syslogSourceType." +
					"This string is used as a substring match against the sourcetype key. ",
			},
			"timestamp_format": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: "Format of timestamp to add at start of the events to be forwarded." +
					"The format is a strftime-style timestamp formatting string. See $SPLUNK_HOME/etc/system/README/outputs.conf.spec for details. ",
			},
			"type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"tcp", "udp"}, false),
				Description:  "Protocol to use to send syslog data. Valid values: (tcp | udp ). ",
			},
			"acl": aclSchema(),
		},
		Read:   outputsTCPSyslogRead,
		Create: outputsTCPSyslogCreate,
		Update: outputsTCPSyslogUpdate,
		Delete: outputsTCPSyslogDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

var validTCPSyslogServer = regexp.MustCompile("^[a-zA-Z0-9\\-.]*:\\d+")

// Functions
func outputsTCPSyslogCreate(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	name := d.Get("name").(string)
	outputsTCPSyslogConfig := getOutputsTCPSyslogConfig(d)
	aclObject := &models.ACLObject{}
	if r, ok := d.GetOk("acl"); ok {
		aclObject = getACLConfig(r.([]interface{}))
	} else {
		aclObject.Owner = "nobody"
		aclObject.App = "search"
	}

	err := (*provider.Client).CreateTCPSyslogOutput(name, aclObject.Owner, aclObject.App, outputsTCPSyslogConfig)
	if err != nil {
		return err
	}

	if _, ok := d.GetOk("acl"); ok {
		err = (*provider.Client).UpdateAcl(aclObject.Owner, aclObject.App, d.Id(), aclObject, "data", "outputs", "tcp", "syslog")
		if err != nil {
			return err
		}
	}

	d.SetId(name)
	return outputsTCPSyslogRead(d, meta)
}

func outputsTCPSyslogRead(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	// We first get list of outputs to get owner and app name for the specific input
	resp, err := (*provider.Client).ReadTCPSyslogOutputs()
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	entry, err := getOutputsTCPSyslogConfigByName(d.Id(), resp)
	if err != nil {
		return err
	}

	if entry == nil {
		return errors.New(fmt.Sprintf("Unable to find resource: %v", d.Id()))
	}

	// Now we read the input configuration with proper owner and app
	resp, err = (*provider.Client).ReadTCPSyslogOutput(entry.Name, entry.ACL.Owner, entry.ACL.App)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	entry, err = getOutputsTCPSyslogConfigByName(d.Id(), resp)
	if err != nil {
		return err
	}

	if entry == nil {
		return errors.New(fmt.Sprintf("Unable to find resource: %v", d.Id()))
	}

	if err = d.Set("name", d.Id()); err != nil {
		return err
	}

	if err = d.Set("priority", entry.Content.Priority); err != nil {
		return err
	}

	if err = d.Set("server", entry.Content.Server); err != nil {
		return err
	}

	if err = d.Set("syslog_sourcetype", entry.Content.SyslogSourceType); err != nil {
		return err
	}

	if err = d.Set("timestamp_format", entry.Content.TimestampFormat); err != nil {
		return err
	}

	if err = d.Set("disabled", entry.Content.Disabled); err != nil {
		return err
	}

	if err = d.Set("type", entry.Content.Type); err != nil {
		return err
	}

	err = d.Set("acl", flattenACL(&entry.ACL))
	if err != nil {
		return err
	}

	return nil
}

func outputsTCPSyslogUpdate(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	outputsTCPSyslogConfig := getOutputsTCPSyslogConfig(d)
	aclObject := getACLConfig(d.Get("acl").([]interface{}))
	err := (*provider.Client).UpdateTCPSyslogOutput(d.Id(), aclObject.Owner, aclObject.App, outputsTCPSyslogConfig)
	if err != nil {
		return err
	}

	//ACL update
	err = (*provider.Client).UpdateAcl(aclObject.Owner, aclObject.App, d.Id(), aclObject, "data", "outputs", "tcp", "syslog")
	if err != nil {
		return err
	}

	return outputsTCPSyslogRead(d, meta)
}

func outputsTCPSyslogDelete(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	aclObject := getACLConfig(d.Get("acl").([]interface{}))
	resp, err := (*provider.Client).DeleteTCPSyslogOutput(d.Id(), aclObject.Owner, aclObject.App)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 200, 201:
		return nil

	default:
		errorResponse := &models.OutputsTCPSyslogResponse{}
		_ = json.NewDecoder(resp.Body).Decode(errorResponse)
		err := errors.New(errorResponse.Messages[0].Text)
		return err
	}
}

// Helpers
func getOutputsTCPSyslogConfig(d *schema.ResourceData) (outputsTCPSyslogObj *models.OutputsTCPSyslogObject) {
	outputsTCPSyslogObj = &models.OutputsTCPSyslogObject{}
	outputsTCPSyslogObj.Disabled = d.Get("disabled").(bool)
	outputsTCPSyslogObj.Priority = d.Get("priority").(int)
	outputsTCPSyslogObj.Server = d.Get("server").(string)
	outputsTCPSyslogObj.SyslogSourceType = d.Get("syslog_sourcetype").(string)
	outputsTCPSyslogObj.TimestampFormat = d.Get("timestamp_format").(string)
	outputsTCPSyslogObj.Type = d.Get("type").(string)
	return
}

func getOutputsTCPSyslogConfigByName(name string, httpResponse *http.Response) (entry *models.OutputsTCPSyslogEntry, err error) {
	response := &models.OutputsTCPSyslogResponse{}
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
