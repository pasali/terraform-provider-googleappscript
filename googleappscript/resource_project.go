package googleappscript

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/pkg/errors"
	"google.golang.org/api/script/v1"
)

var (
	scriptValidFileTypes = []string{"JSON", "HTML", "SERVER_JS", "ENUM_TYPE_UNSPECIFIED"}
)

func resourceProject() *schema.Resource {
	return &schema.Resource{
		Create: resourceProjectCreate,
		Read:   resourceProjectRead,
		Update: resourceProjectUpdate,
		Delete: resourceProjectDelete,

		Schema: map[string]*schema.Schema{
			"title": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"parent_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"update_time": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"script": &schema.Schema{
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},

						"type": &schema.Schema{
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice(scriptValidFileTypes, false),
						},

						"source": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func resourceProjectCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	request := &script.CreateProjectRequest{
		Title:    d.Get("title").(string),
		ParentId: d.Get("parent_id").(string)}
	project, err := config.script.Projects.Create(request).Do()
	if err != nil {
		return errors.Wrap(err, "failed to create project")
	}
	scriptId := project.ScriptId
	scripts := parseScripts(d)
	content := &script.Content{
		ScriptId: scriptId,
		Files:    scripts,
	}
	_, err = config.script.Projects.UpdateContent(scriptId, content).Do()
	if err != nil {
		return errors.Wrap(err, "couldn't update project content")
	}

	d.SetId(scriptId)
	return resourceProjectRead(d, meta)
}

func resourceProjectRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	id := d.Id()
	project, err := config.script.Projects.Get(id).Do()
	if err != nil {
		return errors.Wrap(err, "failed to read project")
	}

	d.Set("title", project.Title)
	d.Set("parent_id", project.ParentId)
	d.Set("script_id", project.ScriptId)
	d.Set("update_time", project.UpdateTime)

	content, err := config.script.Projects.GetContent(id).Do()
	if err != nil {
		return errors.Wrap(err, "failed to read project content")
	}
	d.Set("script", flattenProjectFiles(content.Files))

	return nil
}

func resourceProjectUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	id := d.Id()
	scripts := parseScripts(d)
	content := &script.Content{
		ScriptId: id,
		Files:    scripts,
	}
	_, err := config.script.Projects.UpdateContent(id, content).Do()
	if err != nil {
		return errors.Wrap(err, "couldn't update project content")
	}

	return resourceProjectRead(d, meta)
}

func resourceProjectDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	id := d.Id()
	err := config.drive.Files.Delete(id).Do()
	if err != nil {
		return errors.Wrap(err, "couldn't delete project ")
	}
	return nil
}

func parseScripts(d *schema.ResourceData) []*script.File {
	scriptsRaw := d.Get("script").(*schema.Set)
	scripts := make([]*script.File, scriptsRaw.Len())
	for i, v := range scriptsRaw.List() {
		m := v.(map[string]interface{})
		scripts[i] = &script.File{
			Name:   m["name"].(string),
			Type:   m["type"].(string),
			Source: m["source"].(string),
		}
	}
	return scripts
}

func flattenProjectFiles(list []*script.File) []map[string]interface{} {
	result := make([]map[string]interface{}, len(list))
	for i, v := range list {
		result[i] = map[string]interface{}{
			"name":   v.Name,
			"type":   v.Type,
			"source": v.Source,
		}
	}
	return result
}
