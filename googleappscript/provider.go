package googleappscript

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/pkg/errors"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"token_file": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"googleappscript_project": resourceProject(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	var c Config
	if err := c.loadAndValidate(d.Get("token_file").(string)); err != nil {
		return nil, errors.Wrap(err, "failed to load config")
	}
	return &c, nil
}
