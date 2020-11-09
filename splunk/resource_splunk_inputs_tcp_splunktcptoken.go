package splunk

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"terraform-provider-splunk/client/models"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func inputsTCPSplunkTCPToken() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Required. Name for the token to create.",
			},
			"token": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Optional. Token value to use. If unspecified, a token is generated automatically.",
			},
			"acl": aclSchema(),
		},
		Read:   inputsSplunkTCPTokenRead,
		Create: inputsSplunkTCPTokenCreate,
		Update: inputsSplunkTCPTokenUpdate,
		Delete: inputsSplunkTCPTokenDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

// Functions
func inputsSplunkTCPTokenCreate(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	name := d.Get("name").(string)
	splunkTcpTokenConfigObj := getSplunkTcpTokenConfig(d)
	aclObject := &models.ACLObject{}
	if r, ok := d.GetOk("acl"); ok {
		aclObject = getACLConfig(r.([]interface{}))
	} else {
		aclObject.Owner = "nobody"
		aclObject.App = "search"
	}
	err := (*provider.Client).CreateSplunkTCPTokenInput(aclObject.Owner, aclObject.App, splunkTcpTokenConfigObj)
	if err != nil {
		return err
	}

	if _, ok := d.GetOk("acl"); ok {
		err = (*provider.Client).UpdateAcl(aclObject.Owner, aclObject.App, name, aclObject, "data", "inputs", "tcp", "splunktcptoken")
		if err != nil {
			return err
		}
	}

	d.SetId(name)
	return inputsSplunkTCPTokenRead(d, meta)
}

func inputsSplunkTCPTokenRead(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	// We first get list of inputs to get owner and app name for the specific input
	resp, err := (*provider.Client).ReadSplunkTCPTokenInputs()
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	entry, err := getSplunkTCPTokenConfigByName(d.Id(), resp)
	if err != nil {
		return err
	}

	if entry == nil {
		return fmt.Errorf("Unable to find resource: %v", d.Id())
	}

	// Now we read the input configuration with proper owner and app
	resp, err = (*provider.Client).ReadSplunkTCPTokenInput(entry.Name, entry.ACL.Owner, entry.ACL.App)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	entry, err = getSplunkTCPTokenConfigByName(d.Id(), resp)
	if err != nil {
		return err
	}

	if entry == nil {
		return fmt.Errorf("Unable to find resource: %v", d.Id())
	}

	if err = d.Set("name", d.Id()); err != nil {
		return err
	}

	if err = d.Set("token", entry.Content.Token); err != nil {
		return err
	}

	err = d.Set("acl", flattenACL(&entry.ACL))
	if err != nil {
		return err
	}

	return nil
}

func inputsSplunkTCPTokenUpdate(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	name := "splunktcptoken://" + d.Id()
	splunkTcpTokenConfig := getSplunkTcpTokenConfig(d)
	aclObject := getACLConfig(d.Get("acl").([]interface{}))
	err := (*provider.Client).UpdateSplunkTCPTokenInput(name, aclObject.Owner, aclObject.App, splunkTcpTokenConfig)
	if err != nil {
		return err
	}

	//ACL update
	err = (*provider.Client).UpdateAcl(aclObject.Owner, aclObject.App, d.Id(), aclObject, "data", "inputs", "tcp", "splunktcptoken")
	if err != nil {
		return err
	}

	return inputsSplunkTCPTokenRead(d, meta)
}

func inputsSplunkTCPTokenDelete(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	name := "splunktcptoken://" + d.Id()
	aclObject := getACLConfig(d.Get("acl").([]interface{}))
	resp, err := (*provider.Client).DeleteSplunkTCPTokenInput(name, aclObject.Owner, aclObject.App)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 200, 201:
		return nil

	default:
		errorResponse := &models.InputsSplunkTCPTokenResponse{}
		_ = json.NewDecoder(resp.Body).Decode(errorResponse)
		err := errors.New(errorResponse.Messages[0].Text)
		return err
	}
}

// Helpers
func getSplunkTcpTokenConfig(d *schema.ResourceData) (inputsSplunkTCPTokenObject *models.InputsSplunkTCPTokenObject) {
	inputsSplunkTCPTokenObject = &models.InputsSplunkTCPTokenObject{}
	inputsSplunkTCPTokenObject.Name = d.Get("name").(string)
	inputsSplunkTCPTokenObject.Token = d.Get("token").(string)
	return inputsSplunkTCPTokenObject
}

func getSplunkTCPTokenConfigByName(name string, httpResponse *http.Response) (splunkTCPTokenEntry *models.InputsSplunkTCPTokenEntry, err error) {
	response := &models.InputsSplunkTCPTokenResponse{}
	switch httpResponse.StatusCode {
	case 200, 201:
		_ = json.NewDecoder(httpResponse.Body).Decode(&response)
		re := regexp.MustCompile(`splunktcptoken://(.*)`)
		for _, entry := range response.Entry {
			if name == re.FindStringSubmatch(entry.Name)[1] {
				return &entry, nil
			}
		}

	default:
		_ = json.NewDecoder(httpResponse.Body).Decode(response)
		err := errors.New(response.Messages[0].Text)
		return splunkTCPTokenEntry, err
	}

	return splunkTCPTokenEntry, nil
}
