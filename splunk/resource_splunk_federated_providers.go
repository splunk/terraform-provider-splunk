package splunk

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/splunk/terraform-provider-splunk/client/models"
)

// federatedProviders manages Splunk-to-Splunk federated search providers.
func federatedProviders() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
				Description:  "Unique name for the federated provider.",
			},
			"host_port": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
				Description:  "Remote Splunk management host and port, for example remote-splunk:8089.",
			},
			"service_account": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
				Description:  "Service account username on the remote Splunk deployment.",
			},
			"password": {
				Type:         schema.TypeString,
				Required:     true,
				Sensitive:    true,
				StateFunc:    federatedProviderPasswordState,
				ValidateFunc: validation.StringIsNotEmpty,
				Description:  "Remote service-account password. State stores its SHA-256 hash; configuration changes update the password in place.",
			},
			"mode": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "standard",
				ValidateFunc: validation.StringInSlice([]string{"standard", "transparent"}, false),
				Description:  "Federated-search mode. All providers on a deployment must use the same mode.",
			},
			"app_context": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "search",
				ValidateFunc: validation.StringIsNotEmpty,
				Description:  "Remote app context used by standard-mode federated searches.",
			},
		},
		Create: federatedProviderCreate,
		Read:   federatedProviderRead,
		Update: federatedProviderUpdate,
		Delete: federatedProviderDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func federatedProviderCreate(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	name := d.Get("name").(string)
	if err := (*provider.Client).CreateFederatedProvider(name, federatedProviderConfig(d)); err != nil {
		return err
	}

	d.SetId(name)
	return federatedProviderRead(d, meta)
}

func federatedProviderRead(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	resp, err := (*provider.Client).ReadFederatedProvider(d.Id())
	if resp != nil {
		defer func() { _ = resp.Body.Close() }()
	}
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return nil
		}
		return err
	}

	entry, err := federatedProviderEntry(resp)
	if err != nil {
		return err
	}
	if entry == nil {
		d.SetId("")
		return nil
	}

	if err := d.Set("name", entry.Name); err != nil {
		return err
	}
	if err := d.Set("host_port", entry.Content.HostPort); err != nil {
		return err
	}
	if err := d.Set("service_account", entry.Content.ServiceAccount); err != nil {
		return err
	}
	if err := d.Set("mode", entry.Content.Mode); err != nil {
		return err
	}
	// Transparent mode ignores app context. Preserve configuration instead of
	// replacing it with an absent API value and producing a perpetual diff.
	if entry.Content.Mode == "transparent" {
		return nil
	}
	return d.Set("app_context", entry.Content.AppContext)
}

func federatedProviderUpdate(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	if err := (*provider.Client).UpdateFederatedProvider(d.Id(), federatedProviderConfig(d)); err != nil {
		return err
	}
	return federatedProviderRead(d, meta)
}

func federatedProviderDelete(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*SplunkProvider)
	resp, err := (*provider.Client).DeleteFederatedProvider(d.Id())
	if resp != nil {
		defer func() { _ = resp.Body.Close() }()
	}
	if err != nil && (resp == nil || resp.StatusCode != http.StatusNotFound) {
		return err
	}
	d.SetId("")
	return nil
}

func federatedProviderConfig(d *schema.ResourceData) *models.FederatedProviderObject {
	config := &models.FederatedProviderObject{
		Type:           "splunk",
		Mode:           d.Get("mode").(string),
		HostPort:       d.Get("host_port").(string),
		ServiceAccount: d.Get("service_account").(string),
	}
	// Splunk rejects appContext on transparent provider writes.
	if config.Mode == "standard" {
		config.AppContext = d.Get("app_context").(string)
	}
	// The SDK preserves the original configured value for changed attributes.
	// An unchanged password may come from hashed state and must not be sent.
	if d.Id() == "" || d.HasChange("password") {
		config.Password = d.Get("password").(string)
	}
	return config
}

// federatedProviderPasswordState fingerprints the password for state and diffs.
func federatedProviderPasswordState(value interface{}) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(value.(string))))
}

func federatedProviderEntry(resp *http.Response) (*models.FederatedProviderEntry, error) {
	response := &models.FederatedProviderResponse{}
	if err := json.NewDecoder(resp.Body).Decode(response); err != nil {
		return nil, err
	}
	if len(response.Entry) == 0 {
		return nil, nil
	}
	if len(response.Entry) != 1 {
		return nil, fmt.Errorf("expected one federated provider response entry, got %d", len(response.Entry))
	}
	return &response.Entry[0], nil
}
