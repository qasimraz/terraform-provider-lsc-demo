package provider

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceItem() *schema.Resource {
	fmt.Print()
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the resource, also acts as it's unique ID",
				ForceNew:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "A description of an item",
			},
			"tags": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "An optional list of tags, represented as a key, value pair",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
		Create: resourceCreateItem,
		Read:   resourceReadItem,
		Update: resourceUpdateItem,
		Delete: resourceDeleteItem,
	}
}

func resourceCreateItem(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceReadItem(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceUpdateItem(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceDeleteItem(d *schema.ResourceData, m interface{}) error {
	return nil
}
