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

func authorizationRoles() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Required. The name of the user role to create.",
			},
			"capabilities": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "List of capabilities assigned to role.",
			},
			"cumulative_realtime_search_jobs_quota": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Maximum number of concurrently running real-time searches that all members of this role can have.",
			},
			"cumulative_search_jobs_quota": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Maximum number of concurrently running searches for all role members. Warning message logged when limit is reached.",
			},
			"default_app": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specify the folder name of the default app to use for this role. A user-specific default app overrides this.",
			},
			"imported_roles": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Specify a role to import attributes from. To import multiple roles, specify them separately. " +
					"By default a role imports no other roles. Importing other roles imports all aspects of that role, such as capabilities and allowed indexes to search. " +
					"In combining multiple roles, the effective value for each attribute is the value with the broadest permissions." +
					"nDefault roles: admin, can_delete, power, user. You can specify additional roles created. ",
			},
			"realtime_search_jobs_quota": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Specify the maximum number of concurrent real-time search jobs for this role. This count is independent from the normal search jobs limit.",
			},
			"search_disk_quota": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the maximum disk space in MB that can be used by a user's search jobs. For example, a value of 100 limits this role to 100 MB total.",
			},
			"search_filter": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: "Specify a search string that restricts the scope of searches run by this role. " +
					"Search results for this role only show events that also match the search string you specify. In the case that a user has multiple roles with different search filters, they are combined with an OR.",
			},
			"search_indexes_allowed": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "List of indexes that this role has permissions to search. " +
					"These may be wildcarded, but the index name must begin with an underscore to match internal indexes." +
					"Search indexes available by default include the following. " +
					"All internal indexes    " +
					"All non-internal indexes    " +
					"_audit    " +
					"_blocksignature    " +
					"_internal    " +
					"_thefishbucket    " +
					"history    " +
					"main    " +
					"You can also specify other search indexes added to the server. ",
			},
			"search_indexes_default": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "For this role, list of indexes to search when no index is specified. " +
					"These indexes can be wildcarded, with the exception that '*' does not match internal indexes. " +
					"To match internal indexes, start with '_'. All internal indexes are represented by '_*'. " +
					"A user with this role can search other indexes using \"index= \" For example, \"index=special_index\". " +
					"Search indexes available by default include the following.     " +
					"All internal indexes    " +
					"All non-internal indexes    " +
					"_audit    " +
					"_blocksignature    " +
					"_internal    " +
					"_thefishbucket    " +
					"history    " +
					"main    " +
					"other search indexes added to the server",
			},
			"search_jobs_quota": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The maximum number of concurrent searches a user with this role is allowed to run. For users with multiple roles, the maximum quota value among all of the roles applies.",
			},
			"search_time_win": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				Description: "Maximum time span of a search, in seconds. By default, searches are not limited to any specific time window. " +
					"To override any search time windows from imported roles, set srchTimeWin to '0', as the 'admin' role does.",
			},
		},
		Read:   authorizationRolesRead,
		Create: authorizationRolesCreate,
		Delete: authorizationRolesDelete,
		Update: authorizationRolesUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

// Functions
func authorizationRolesCreate(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	name := d.Get("name").(string)
	authenticationUserObj := getAuthorizationRolesConfig(d)
	err := (*provider.Client).CreateAuthorizationRoles(name, authenticationUserObj)
	if err != nil {
		return err
	}

	d.SetId(name)
	return authorizationRolesRead(d, meta)
}

func authorizationRolesRead(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	name := d.Id()
	// We first get list of inputs to get owner and app name for the specific input
	resp, err := (*provider.Client).ReadAllAuthorizationRoles()
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	entry, err := getAuthorizationRolesByName(name, resp)
	if err != nil {
		return err
	}

	if entry == nil {
		return fmt.Errorf("Unable to find resource: %v", name)
	}

	// Now we read the input configuration with proper owner and app
	resp, err = (*provider.Client).ReadAuthorizationRoles(name, entry.ACL.Owner, entry.ACL.App)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	entry, err = getAuthorizationRolesByName(name, resp)
	if err != nil {
		return err
	}

	if entry == nil {
		return fmt.Errorf("Unable to find resource: %v", name)
	}

	if err = d.Set("capabilities", entry.Content.Capabilities); err != nil {
		return err
	}

	if err = d.Set("cumulative_realtime_search_jobs_quota", entry.Content.CumulativeRTSrchJobsQuota); err != nil {
		return err
	}

	if err = d.Set("cumulative_search_jobs_quota", entry.Content.CumulativeSrchJobsQuota); err != nil {
		return err
	}

	if err = d.Set("default_app", entry.Content.DefaultApp); err != nil {
		return err
	}

	if err = d.Set("imported_roles", entry.Content.ImportedRoles); err != nil {
		return err
	}

	if err = d.Set("name", entry.Name); err != nil {
		return err
	}

	if err = d.Set("realtime_search_jobs_quota", entry.Content.RtSrchJobsQuota); err != nil {
		return err
	}

	if err = d.Set("search_disk_quota", entry.Content.SrchDiskQuota); err != nil {
		return err
	}

	if err = d.Set("search_filter", entry.Content.SrchFilter); err != nil {
		return err
	}

	if err = d.Set("search_indexes_allowed", entry.Content.SrchIndexesAllowed); err != nil {
		return err
	}

	if err = d.Set("search_indexes_default", entry.Content.SrchIndexesDefault); err != nil {
		return err
	}

	if err = d.Set("search_jobs_quota", entry.Content.SrchJobsQuota); err != nil {
		return err
	}

	if err = d.Set("search_time_win", entry.Content.SrchTimeWin); err != nil {
		return err
	}

	return nil
}

func authorizationRolesUpdate(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	authenticationUserObj := getAuthorizationRolesConfig(d)
	err := (*provider.Client).UpdateAuthorizationRoles(d.Id(), authenticationUserObj)
	if err != nil {
		return err
	}

	return authorizationRolesRead(d, meta)
}

func authorizationRolesDelete(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	resp, err := (*provider.Client).DeleteAuthorizationRoles(d.Id())
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 200, 201:
		return nil

	default:
		errorResponse := &models.AuthorizationRolesResponse{}
		_ = json.NewDecoder(resp.Body).Decode(errorResponse)
		err := errors.New(errorResponse.Messages[0].Text)
		return err
	}
}

// Helpers
func getAuthorizationRolesConfig(d *schema.ResourceData) (authenticationUserObject *models.AuthorizationRolesObject) {
	authenticationUserObject = &models.AuthorizationRolesObject{}
	if val, ok := d.GetOk("capabilities"); ok {
		for _, v := range val.([]interface{}) {
			authenticationUserObject.Capabilities = append(authenticationUserObject.Capabilities, v.(string))
		}
	}
	authenticationUserObject.CumulativeRTSrchJobsQuota = d.Get("cumulative_realtime_search_jobs_quota").(int)
	authenticationUserObject.CumulativeSrchJobsQuota = d.Get("cumulative_search_jobs_quota").(int)
	authenticationUserObject.DefaultApp = d.Get("default_app").(string)
	if val, ok := d.GetOk("imported_roles"); ok {
		for _, v := range val.([]interface{}) {
			authenticationUserObject.ImportedRoles = append(authenticationUserObject.ImportedRoles, v.(string))
		}
	}
	authenticationUserObject.Name = d.Get("name").(string)
	authenticationUserObject.RtSrchJobsQuota = d.Get("realtime_search_jobs_quota").(int)
	authenticationUserObject.SrchDiskQuota = d.Get("search_disk_quota").(int)
	authenticationUserObject.SrchFilter = d.Get("search_filter").(string)
	if val, ok := d.GetOk("search_indexes_allowed"); ok {
		for _, v := range val.([]interface{}) {
			authenticationUserObject.SrchIndexesAllowed = append(authenticationUserObject.SrchIndexesAllowed, v.(string))
		}
	}
	if val, ok := d.GetOk("search_indexes_default"); ok {
		for _, v := range val.([]interface{}) {
			authenticationUserObject.SrchIndexesDefault = append(authenticationUserObject.SrchIndexesDefault, v.(string))
		}
	}
	authenticationUserObject.SrchJobsQuota = d.Get("search_jobs_quota").(int)
	authenticationUserObject.SrchTimeWin = d.Get("search_time_win").(int)
	return authenticationUserObject
}

func getAuthorizationRolesByName(name string, httpResponse *http.Response) (AuthorizationRolesEntry *models.AuthorizationRolesEntry, err error) {
	response := &models.AuthorizationRolesResponse{}
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
		return AuthorizationRolesEntry, err
	}

	return AuthorizationRolesEntry, nil
}
