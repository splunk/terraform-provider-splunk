package splunk

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"terraform-provider-splunk/client/models"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func outputsTCPDefault() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"tcpout"}, false),
				Description:  "Configuration to be edited. The only valid value is tcpout",
			},
			"default_group": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: "Comma-separated list of one or more target group names, specified later in [tcpout:<target_group>] stanzas of outputs.conf.spec file." +
					"The forwarder sends all data to the specified groups. If you do not want to forward data automatically, do not set this attribute. " +
					"Can be overridden by an inputs.conf _TCP_ROUTING setting, which in turn can be overridden by a props.conf/transforms.conf modifier.",
			},
			"disabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Disables default tcpout settings",
			},
			"drop_events_on_queue_full": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				Description: "If set to a positive number, wait the specified number of seconds before throwing out all new events until the output queue has space. " +
					"Defaults to -1 (do not drop events)." +
					"CAUTION: Do not set this value to a positive integer if you are monitoring files." +
					"Setting this to -1 or 0 causes the output queue to block when it gets full, which causes further blocking up the processing chain. " +
					"If any target group queue is blocked, no more data reaches any other target group." +
					"Using auto load-balancing is the best way to minimize this condition, because, in that case," +
					" multiple receivers must be down (or jammed up) before queue blocking can occur. ",
			},
			"heartbeat_frequency": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				Description: "How often (in seconds) to send a heartbeat packet to the receiving server." +
					"Heartbeats are only sent if sendCookedData=true. Defaults to 30 seconds. ",
			},
			"index_and_forward": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
				Description: "Specifies whether to index all data locally, in addition to forwarding it. " +
					"Defaults to false." +
					"This is known as an \"index-and-forward\" configuration. " +
					"This attribute is only available for heavy forwarders. " +
					"It is available only at the top level [tcpout] stanza in outputs.conf. " +
					"It cannot be overridden in a target group. ",
			},
			"max_queue_size": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringMatch(validMaxQueueSize, "valid values: integer[KB|MB|GB]"),
				Description: "Specify an integer or integer[KB|MB|GB]." +
					"Sets the maximum size of the forwarder output queue. " +
					"It also sets the maximum size of the wait queue to 3x this value, if you have enabled indexer acknowledgment (useACK=true)." +
					"Although the wait queue and the output queues are both configured by this attribute, they are separate queues. " +
					"The setting determines the maximum size of the queue in-memory (RAM) buffer." +
					"For heavy forwarders sending parsed data, maxQueueSize is the maximum number of events. " +
					"Since events are typically much shorter than data blocks, the memory consumed by the queue " +
					"on a parsing forwarder is likely to be much smaller than on a non-parsing forwarder, if you use this version of the setting." +
					"If specified as a lone integer (for example, maxQueueSize=100), " +
					"maxQueueSize indicates the maximum number of queued events (for parsed data) or blocks of data (for unparsed data). " +
					"A block of data is approximately 64KB. For non-parsing forwarders, such as universal forwarders, " +
					"that send unparsed data, maxQueueSize is the maximum number of data blocks." +
					"If specified as an integer followed by KB, MB, or GB (for example, maxQueueSize=100MB), " +
					"maxQueueSize indicates the maximum RAM allocated to the queue buffer. " +
					"Defaults to 500KB (which means a maximum size of 500KB for the output queue and 1500KB for the wait queue, if any). ",
			},
			"send_cooked_data": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
				Description: "If true, events are cooked (processed by Splunk software). If false, events are raw and untouched prior to sending. " +
					"Defaults to true. Set to false if you are sending to a third-party system. ",
			},
			"acl": aclSchema(),
		},
		Read:   outputsTCPDefaultRead,
		Create: outputsTCPDefaultCreate,
		Update: outputsTCPDefaultUpdate,
		Delete: outputsTCPDefaultDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

var validMaxQueueSize = regexp.MustCompile(`^\d+[KB|MB|GB]+`)

// Functions
func outputsTCPDefaultCreate(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	name := d.Get("name").(string)
	outputsTCPDefaultConfig := getOutputsTCPDefaultConfig(d)
	aclObject := &models.ACLObject{}
	if r, ok := d.GetOk("acl"); ok {
		aclObject = getACLConfig(r.([]interface{}))
	} else {
		aclObject.Owner = "nobody"
		aclObject.App = "search"
	}

	err := (*provider.Client).CreateTCPDefaultOutput(name, aclObject.Owner, aclObject.App, outputsTCPDefaultConfig)
	if err != nil {
		return err
	}

	if _, ok := d.GetOk("acl"); ok {
		err = (*provider.Client).UpdateAcl(aclObject.Owner, aclObject.App, d.Id(), aclObject, "data", "outputs", "tcp", "default")
		if err != nil {
			return err
		}
	}

	d.SetId(name)
	return outputsTCPDefaultRead(d, meta)
}

func outputsTCPDefaultRead(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	// We first get list of outputs to get owner and app name for the specific input
	resp, err := (*provider.Client).ReadTCPDefaultOutputs()
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	entry, err := getOutputsTCPDefaultConfigByName(d.Id(), resp)
	if err != nil {
		return err
	}

	if entry == nil {
		return fmt.Errorf("Unable to find resource: %v", d.Id())
	}

	// Now we read the input configuration with proper owner and app
	resp, err = (*provider.Client).ReadTCPDefaultOutput(entry.Name, entry.ACL.Owner, entry.ACL.App)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	entry, err = getOutputsTCPDefaultConfigByName(d.Id(), resp)
	if err != nil {
		return err
	}

	if entry == nil {
		return fmt.Errorf("Unable to find resource: %v", d.Id())
	}

	if err = d.Set("name", d.Id()); err != nil {
		return err
	}

	if err = d.Set("default_group", entry.Content.DefaultGroup); err != nil {
		return err
	}

	if err = d.Set("drop_events_on_queue_full", entry.Content.DropEventsOnQueueFull); err != nil {
		return err
	}

	if err = d.Set("heartbeat_frequency", entry.Content.HeartbeatFrequency); err != nil {
		return err
	}

	if err = d.Set("max_queue_size", entry.Content.MaxQueueSize); err != nil {
		return err
	}

	if err = d.Set("disabled", entry.Content.Disabled); err != nil {
		return err
	}

	if err = d.Set("index_and_forward", entry.Content.IndexAndForward); err != nil {
		return err
	}

	if err = d.Set("send_cooked_data", entry.Content.SendCookedData); err != nil {
		return err
	}

	err = d.Set("acl", flattenACL(&entry.ACL))
	if err != nil {
		return err
	}

	return nil
}

func outputsTCPDefaultUpdate(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	outputsTCPDefaultConfig := getOutputsTCPDefaultConfig(d)
	aclObject := getACLConfig(d.Get("acl").([]interface{}))
	err := (*provider.Client).UpdateTCPDefaultOutput(d.Id(), aclObject.Owner, aclObject.App, outputsTCPDefaultConfig)
	if err != nil {
		return err
	}

	//ACL update
	err = (*provider.Client).UpdateAcl(aclObject.Owner, aclObject.App, d.Id(), aclObject, "data", "outputs", "tcp", "default")
	if err != nil {
		return err
	}

	return outputsTCPDefaultRead(d, meta)
}

func outputsTCPDefaultDelete(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	aclObject := getACLConfig(d.Get("acl").([]interface{}))
	resp, err := (*provider.Client).DeleteTCPDefaultOutput(d.Id(), aclObject.Owner, aclObject.App)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 200, 201:
		return nil

	default:
		errorResponse := &models.OutputsTCPDefaultResponse{}
		_ = json.NewDecoder(resp.Body).Decode(errorResponse)
		err := errors.New(errorResponse.Messages[0].Text)
		return err
	}
}

// Helpers
func getOutputsTCPDefaultConfig(d *schema.ResourceData) (outputsTCPDefaultObj *models.OutputsTCPDefaultObject) {
	outputsTCPDefaultObj = &models.OutputsTCPDefaultObject{}
	outputsTCPDefaultObj.DefaultGroup = d.Get("default_group").(string)
	outputsTCPDefaultObj.DropEventsOnQueueFull = d.Get("drop_events_on_queue_full").(int)
	outputsTCPDefaultObj.HeartbeatFrequency = d.Get("heartbeat_frequency").(int)
	outputsTCPDefaultObj.MaxQueueSize = d.Get("max_queue_size").(string)
	outputsTCPDefaultObj.Disabled = d.Get("disabled").(bool)
	outputsTCPDefaultObj.IndexAndForward = d.Get("index_and_forward").(bool)
	outputsTCPDefaultObj.SendCookedData = d.Get("send_cooked_data").(bool)
	return
}

func getOutputsTCPDefaultConfigByName(name string, httpResponse *http.Response) (entry *models.OutputsTCPDefaultEntry, err error) {
	response := &models.OutputsTCPDefaultResponse{}
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
