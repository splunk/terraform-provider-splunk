package splunk

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"

	"github.com/nealbrown/terraform-provider-splunk/client/models"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func authenticationUsers() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Required. Unique user login name.",
			},
			"default_app": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "User default app. Overrides the default app inherited from the user roles. ",
			},
			"email": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "User email address.",
			},
			"force_change_pass": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Force user to change password indication",
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "User login password.",
			},
			"realname": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Full user name.",
			},
			"restart_background_jobs": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Restart background search job that has not completed when Splunk restarts indication.",
			},
			"roles": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Role to assign to this user. At least one existing role is required.",
			},
			"tz": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "User timezone.",
			},
		},
		Read:   authenticationUsersRead,
		Create: authenticationUsersCreate,
		Delete: authenticationUsersDelete,
		Update: authenticationUsersUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

// Functions
func authenticationUsersCreate(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	name := d.Get("name").(string)
	authenticationUserObj := getAuthenticationUserConfig(d)
	err := (*provider.Client).CreateAuthenticationUser(name, authenticationUserObj)
	if err != nil {
		return err
	}

	d.SetId(name)
	return authenticationUsersRead(d, meta)
}

func authenticationUsersRead(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	name := d.Id()
	// We first get list of inputs to get owner and app name for the specific input
	resp, err := (*provider.Client).ReadAuthenticationUsers()
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	entry, err := getAuthenticationUserByName(name, resp)
	if err != nil {
		return err
	}

	if entry == nil {
		return fmt.Errorf("Unable to find resource: %v", name)
	}

	// Now we read the input configuration with proper owner and app
	resp, err = (*provider.Client).ReadAuthenticationUser(name, entry.ACL.Owner, entry.ACL.App)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	entry, err = getAuthenticationUserByName(name, resp)
	if err != nil {
		return err
	}

	if entry == nil {
		return fmt.Errorf("Unable to find resource: %v", name)
	}

	if err = d.Set("default_app", entry.Content.DefaultApp); err != nil {
		return err
	}

	if err = d.Set("email", entry.Content.Email); err != nil {
		return err
	}

	if err = d.Set("name", entry.Name); err != nil {
		return err
	}

	if err = d.Set("realname", entry.Content.RealName); err != nil {
		return err
	}

	if err = d.Set("restart_background_jobs", entry.Content.RestartBackgroundJobs); err != nil {
		return err
	}

	if err = d.Set("roles", entry.Content.Roles); err != nil {
		return err
	}

	if err = d.Set("tz", entry.Content.TZ); err != nil {
		return err
	}

	return nil
}

func authenticationUsersUpdate(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	authenticationUserObj := getAuthenticationUserConfig(d)
	err := (*provider.Client).UpdateAuthenticationUser(d.Id(), authenticationUserObj)
	if err != nil {
		return err
	}

	return authenticationUsersRead(d, meta)
}

func authenticationUsersDelete(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	resp, err := (*provider.Client).DeleteAuthenticationUser(d.Id())
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 200, 201:
		return nil

	default:
		errorResponse := &models.AuthenticationUserResponse{}
		_ = json.NewDecoder(resp.Body).Decode(errorResponse)
		err := errors.New(errorResponse.Messages[0].Text)
		return err
	}
}

// Helpers
func getAuthenticationUserConfig(d *schema.ResourceData) (authenticationUserObject *models.AuthenticationUserObject) {
	authenticationUserObject = &models.AuthenticationUserObject{}
	authenticationUserObject.DefaultApp = d.Get("default_app").(string)
	authenticationUserObject.Email = d.Get("email").(string)
	authenticationUserObject.ForceChangePass = d.Get("force_change_pass").(bool)
	authenticationUserObject.Name = d.Get("name").(string)
	authenticationUserObject.Password = d.Get("password").(string)
	authenticationUserObject.RealName = d.Get("realname").(string)
	authenticationUserObject.RestartBackgroundJobs = d.Get("restart_background_jobs").(bool)
	authenticationUserObject.TZ = d.Get("tz").(string)
	if val, ok := d.GetOk("roles"); ok {
		for _, v := range val.([]interface{}) {
			authenticationUserObject.Roles = append(authenticationUserObject.Roles, v.(string))
		}
	}
	return authenticationUserObject
}

func getAuthenticationUserByName(name string, httpResponse *http.Response) (AuthenticationUserEntry *models.AuthenticationUserEntry, err error) {
	response := &models.AuthenticationUserResponse{}
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
		return AuthenticationUserEntry, err
	}

	return AuthenticationUserEntry, nil
}
