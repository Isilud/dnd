package main

import (
    "encoding/json"
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"io/ioutil"
	"os"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() *schema.Provider {
			return Provider()
		},
	})
}

// Provider -
func Provider() *schema.Provider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"dnd_character": dnd_characterResource(),
		},
		DataSourcesMap: map[string]*schema.Resource{},
	}
}

func dnd_characterResource() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCharacterCreate,
		ReadContext:   resourceCharacterRead,
		UpdateContext: resourceCharacterUpdate,
		DeleteContext: resourceCharacterDelete,
		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"class": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Paysan",
			},
			"race": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Humain",
			},
			"niveau": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
			},
			"force": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},
			"intelligence": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},
			"sagesse": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},
			"dexterite": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},
			"constitution": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},
			"charisme": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},
		},
	}
}

func resourceCharacterCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	d.SetId(d.Get("name").(string))
	file := "./Fiches/"+d.Id()+".json"
	result := charToMap(d)

	jsonMarshall, _ := json.MarshalIndent(result, "", " ")
    ioutil.WriteFile(file, jsonMarshall, 0644)
	return diags
}

func resourceCharacterRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	
	oldId := d.Id()

	file := "./Fiches/"+oldId+".json"

	if doesntExist(file) {
  		d.SetId("")
	}
	
	jsonFile, _ := os.Open(file)
    defer jsonFile.Close()

    oldData := make(map[string]interface{})
    newData := charToMap(d)

    byteValue, _ := ioutil.ReadAll(jsonFile)
    json.Unmarshal(byteValue, &oldData)
    for key, _ := range newData {
        if err := d.Set(key, oldData[key]); err != nil {
    		d.Set(key, "")
  		}
    }

	return diags
}

func resourceCharacterUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	file := "./Fiches/"+d.Id()+".json"
	result := charToMap(d)
	jsonMarshall, _ := json.MarshalIndent(result, "", " ")
    ioutil.WriteFile(file, jsonMarshall, 0644)
	return resourceCharacterRead(ctx, d, m)
}

func resourceCharacterDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	file := "./Fiches/"+d.Get("name").(string)+".json"
	os.Remove(file)

	return diags
}

func charToMap (d *schema.ResourceData) map[string]interface{}{

	result := make(map[string]interface{})

	attribut := []string{"name", "class", "race"}
	for _, att := range attribut {
		result[att] = d.Get(att).(string)
	}

	attribut = []string{"niveau", "force", "intelligence", "sagesse", "dexterite", "constitution", "charisme"}
	for _, att := range attribut {
		result[att] = d.Get(att).(int)
	}

	return result
}

func doesntExist(file string) bool{
	if _, err := os.Stat(file); err == nil {
		return false
	}
	return true
}