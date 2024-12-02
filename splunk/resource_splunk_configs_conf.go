package splunk

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"

	"github.com/splunk/terraform-provider-splunk/client/models"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func configsConf() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"variables": {
				Type:        schema.TypeMap,
				Optional:    true,
				Computed:    true,
				Description: `A map of key value pairs for a stanza.`,
				Elem: &schema.Schema{
					Type:     schema.TypeString,
					Computed: true,
				},
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[a-zA-Z0-9_\-.]+/[a-zA-Z0-9_\-.:/]+`), "A '/' separated string consisting of {conf_file_name}/{stanza_name} ex. props/custom_stanza"),
				Description:  `A '/' separated string consisting of {conf_file_name}/{stanza_name} ex. props/custom_stanza`,
			},
			"acl": aclSchema(),
		},
		Read:   configsConfRead,
		Create: configsConfCreate,
		Delete: configsConfDelete,
		Update: configsConfUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

// Functions
func configsConfCreate(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	name := d.Get("name").(string)
	configsConfConfigObj := getConfigsConfConfig(d)
	aclObject := &models.ACLObject{}
	if r, ok := d.GetOk("acl"); ok {
		aclObject = getACLConfig(r.([]interface{}))
	} else {
		aclObject.Owner = "nobody"
		aclObject.App = "search"
	}
	err := (*provider.Client).CreateConfigsConfObject(name, aclObject.Owner, aclObject.App, configsConfConfigObj)
	if err != nil {
		return err
	}
	if _, ok := d.GetOk("acl"); ok {
		conf, stanza := (*provider.Client).SplitConfStanza(name)
		err := (*provider.Client).UpdateAcl(aclObject.Owner, aclObject.App, stanza, aclObject, "configs", "conf-"+conf)
		if err != nil {
			return err
		}
	}

	d.SetId(name)
	return configsConfRead(d, meta)
}

func configsConfRead(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	name := d.Id()
	_, stanza := (*provider.Client).SplitConfStanza(name)

	// We first get list of stanzas in a conf file to get owner and app name for the specific stanza
	resp, err := (*provider.Client).ReadAllConfigsConfObject(name)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	entry, err := getConfigsConfConfigByName(stanza, resp)
	if err != nil {
		return err
	}

	if entry == nil {
		return fmt.Errorf("Unable to find resource: %v", name)
	}

	// Now we read the input configuration with proper owner and app
	resp, err = (*provider.Client).ReadConfigsConfObject(name, entry.ACL.Owner, entry.ACL.App)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	contentResp, err := (*provider.Client).ReadConfigsConfObject(name, entry.ACL.Owner, entry.ACL.App)
	if err != nil {
		return err
	}
	defer contentResp.Body.Close()

	var result map[string]interface{}
	b, _ := io.ReadAll(contentResp.Body)

	err = json.Unmarshal(b, &result)
	if err != nil {
		return err
	}
	content := result["entry"].([]interface{})[0].(map[string]interface{})["content"].(map[string]interface{})

	re := regexp.MustCompile(`eai:.*`)

	for key := range content {
		if re.MatchString(key) {
			delete(content, key)
		}
	}

	// Override value to convert bool Type to string
	content["disabled"] = strconv.FormatBool(content["disabled"].(bool))

	entry, err = getConfigsConfConfigByName(stanza, resp)
	if err != nil {
		return err
	}

	if err = d.Set("name", d.Id()); err != nil {
		return err
	}

	if err = d.Set("variables", content); err != nil {
		return err
	}

	err = d.Set("acl", flattenACL(&entry.ACL))
	if err != nil {
		return err
	}

	return nil
}

func configsConfUpdate(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	configsConfConfigObj := getConfigsConfConfig(d)
	aclObject := getACLConfig(d.Get("acl").([]interface{}))
	// Update will create a new resource with private `user` permissions if resource had shared permissions set
	var owner string
	if aclObject.Sharing != "user" {
		owner = "nobody"
	} else {
		owner = aclObject.Owner
	}
	name := d.Id()
	conf, stanza := (*provider.Client).SplitConfStanza(name)

	err := (*provider.Client).UpdateConfigsConfObject(d.Id(), owner, aclObject.App, configsConfConfigObj)
	if err != nil {
		return err
	}

	//ACL update
	err = (*provider.Client).UpdateAcl(aclObject.Owner, aclObject.App, stanza, aclObject, "configs", "conf-"+conf)
	if err != nil {
		return err
	}

	return configsConfRead(d, meta)
}

func configsConfDelete(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	aclObject := getACLConfig(d.Get("acl").([]interface{}))

	resp, err := (*provider.Client).DeleteConfigsConfObject(d.Id(), aclObject.Owner, aclObject.App)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 200, 201:
		return nil

	default:
		errorResponse := &models.ConfigsConfResponse{}
		_ = json.NewDecoder(resp.Body).Decode(errorResponse)
		err := errors.New(errorResponse.Messages[0].Text)
		return err
	}
}

// Helpers
func getConfigsConfConfig(d *schema.ResourceData) (configsConfConfigObject *models.ConfigsConfObject) {
	configsConfConfigObject = &models.ConfigsConfObject{}
	mapInterface := d.Get("variables").(map[string]interface{})
	mapString := make(map[string]string)
	for key, value := range mapInterface {
		strKey := fmt.Sprintf("%v", key)
		strValue := fmt.Sprintf("%v", value)

		mapString[strKey] = strValue
	}
	configsConfConfigObject.Variables = mapString

	return configsConfConfigObject
}

func getConfigsConfConfigByName(name string, httpResponse *http.Response) (configsConfEntry *models.ConfigsConfEntry, err error) {
	response := &models.ConfigsConfResponse{}

	switch httpResponse.StatusCode {
	case 200, 201:
		_ = json.NewDecoder(httpResponse.Body).Decode(&response)
		re := regexp.MustCompile(`(.*)`)

		for _, entry := range response.Entry {
			if name == re.FindStringSubmatch(entry.Name)[1] {
				return &entry, nil
			}
		}

	default:
		_ = json.NewDecoder(httpResponse.Body).Decode(response)
		err := errors.New(response.Messages[0].Text)
		return configsConfEntry, err
	}

	return configsConfEntry, nil
}
