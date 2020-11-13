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

func outputsTCPServer() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringMatch(validTCPServer, "<host>:<port> of the Splunk receiver"),
				Description: "<host>:<port> of the Splunk receiver. " +
					"<host> can be either an ip address or server name. " +
					"<port> is the that port that the Splunk receiver is listening on.",
			},
			"disabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Disables default tcpout settings",
			},
			"method": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"clone", "balance", "autobalance"}, false),
				Description: "Valid values: (clone | balance | autobalance)" +
					"The data distribution method used when two or more servers exist in the same forwarder group. ",
			},
			"ssl_alt_name_to_check": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The alternate name to match in the remote server's SSL certificate. ",
			},
			"ssl_cert_path": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The alternate name to match in the remote server's SSL certificate. ",
			},
			"ssl_cipher": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "SSL Cipher in the form ALL:!aNULL:!eNULL:!LOW:!EXP:RC4+RSA:+HIGH:+MEDIUM ",
			},
			"ssl_common_name_to_check": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: "Check the common name of the server's certificate against this name." +
					"If there is no match, assume that Splunk Enterprise is not authenticated against this server. " +
					"You must specify this setting if sslVerifyServerCert is true.",
			},
			"ssl_root_ca_path": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The path to the root certificate authority file (optional). ",
			},
			"ssl_password": {
				Type:      schema.TypeString,
				Optional:  true,
				Computed:  true,
				Sensitive: true,
				Description: "The password associated with the CAcert." +
					"The default Splunk Enterprise CAcert uses the password password.",
			},
			"ssl_verify_server_cert": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
				Description: " 	If true, make sure that the server you are connecting to is a valid one (authenticated). " +
					"Both the common name and the alternate name of the server are then checked for a match.",
			},
			"acl": aclSchema(),
		},
		Read:   outputsTCPServerRead,
		Create: outputsTCPServerCreate,
		Update: outputsTCPServerUpdate,
		Delete: outputsTCPServerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

var validTCPServer = regexp.MustCompile(`^[a-zA-Z0-9\-.]*:\d+`)

// Functions
func outputsTCPServerCreate(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	name := d.Get("name").(string)
	outputsTCPServerConfig := getOutputsTCPServerConfig(d)
	aclObject := &models.ACLObject{}
	if r, ok := d.GetOk("acl"); ok {
		aclObject = getACLConfig(r.([]interface{}))
	} else {
		aclObject.Owner = "nobody"
		aclObject.App = "search"
	}

	err := (*provider.Client).CreateTCPServerOutput(name, aclObject.Owner, aclObject.App, outputsTCPServerConfig)
	if err != nil {
		return err
	}

	if _, ok := d.GetOk("acl"); ok {
		err = (*provider.Client).UpdateAcl(aclObject.Owner, aclObject.App, d.Id(), aclObject, "data", "outputs", "tcp", "server")
		if err != nil {
			return err
		}
	}

	d.SetId(name)
	return outputsTCPServerRead(d, meta)
}

func outputsTCPServerRead(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	// We first get list of outputs to get owner and app name for the specific input
	resp, err := (*provider.Client).ReadTCPServerOutputs()
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	entry, err := getOutputsTCPServerConfigByName(d.Id(), resp)
	if err != nil {
		return err
	}

	if entry == nil {
		return fmt.Errorf("Unable to find resource: %v", d.Id())
	}

	// Now we read the input configuration with proper owner and app
	resp, err = (*provider.Client).ReadTCPServerOutput(entry.Name, entry.ACL.Owner, entry.ACL.App)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	entry, err = getOutputsTCPServerConfigByName(d.Id(), resp)
	if err != nil {
		return err
	}

	if entry == nil {
		return fmt.Errorf("Unable to find resource: %v", d.Id())
	}

	if err = d.Set("name", d.Id()); err != nil {
		return err
	}

	if err = d.Set("method", entry.Content.Method); err != nil {
		return err
	}

	if err = d.Set("ssl_alt_name_to_check", entry.Content.SSLAltNameToCheck); err != nil {
		return err
	}

	if err = d.Set("ssl_cert_path", entry.Content.SSLCertPath); err != nil {
		return err
	}

	if err = d.Set("ssl_cipher", entry.Content.SSLCipher); err != nil {
		return err
	}

	if err = d.Set("disabled", entry.Content.Disabled); err != nil {
		return err
	}

	if err = d.Set("ssl_common_name_to_check", entry.Content.SSLCommonNameToCheck); err != nil {
		return err
	}

	if err = d.Set("ssl_password", entry.Content.SSLPassword); err != nil {
		return err
	}

	if err = d.Set("ssl_root_ca_path", entry.Content.SSLRootCAPath); err != nil {
		return err
	}

	if err = d.Set("ssl_verify_server_cert", entry.Content.SSLVerifyServerCert); err != nil {
		return err
	}

	err = d.Set("acl", flattenACL(&entry.ACL))
	if err != nil {
		return err
	}

	return nil
}

func outputsTCPServerUpdate(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	outputsTCPServerConfig := getOutputsTCPServerConfig(d)
	aclObject := getACLConfig(d.Get("acl").([]interface{}))
	err := (*provider.Client).UpdateTCPServerOutput(d.Id(), aclObject.Owner, aclObject.App, outputsTCPServerConfig)
	if err != nil {
		return err
	}

	//ACL update
	err = (*provider.Client).UpdateAcl(aclObject.Owner, aclObject.App, d.Id(), aclObject, "data", "outputs", "tcp", "server")
	if err != nil {
		return err
	}

	return outputsTCPServerRead(d, meta)
}

func outputsTCPServerDelete(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	aclObject := getACLConfig(d.Get("acl").([]interface{}))
	resp, err := (*provider.Client).DeleteTCPServerOutput(d.Id(), aclObject.Owner, aclObject.App)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 200, 201:
		return nil

	default:
		errorResponse := &models.OutputsTCPServerResponse{}
		_ = json.NewDecoder(resp.Body).Decode(errorResponse)
		err := errors.New(errorResponse.Messages[0].Text)
		return err
	}
}

// Helpers
func getOutputsTCPServerConfig(d *schema.ResourceData) (outputsTCPServerObj *models.OutputsTCPServerObject) {
	outputsTCPServerObj = &models.OutputsTCPServerObject{}
	outputsTCPServerObj.Disabled = d.Get("disabled").(bool)
	outputsTCPServerObj.Method = d.Get("method").(string)
	outputsTCPServerObj.SSLAltNameToCheck = d.Get("ssl_alt_name_to_check").(string)
	outputsTCPServerObj.SSLCertPath = d.Get("ssl_cert_path").(string)
	outputsTCPServerObj.SSLCipher = d.Get("ssl_cipher").(string)
	outputsTCPServerObj.SSLCommonNameToCheck = d.Get("ssl_common_name_to_check").(string)
	outputsTCPServerObj.SSLPassword = d.Get("ssl_password").(string)
	outputsTCPServerObj.SSLRootCAPath = d.Get("ssl_root_ca_path").(string)
	outputsTCPServerObj.SSLVerifyServerCert = d.Get("ssl_verify_server_cert").(bool)
	return
}

func getOutputsTCPServerConfigByName(name string, httpResponse *http.Response) (entry *models.OutputsTCPServerEntry, err error) {
	response := &models.OutputsTCPServerResponse{}
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
