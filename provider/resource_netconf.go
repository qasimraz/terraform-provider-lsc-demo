package provider

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceNetconf() *schema.Resource {
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
		},
		Create: resourceCreateItem,
		Read:   resourceReadItem,
		Update: resourceUpdateItem,
		Delete: resourceDeleteItem,
		Exists: resourceExistsItem,
	}
}

func resourceCreateItem(d *schema.Resource, m interface{}) error {
	return resourceReadItem(d, m)
}

func resourceReadItem(d *schema.Resource, m interface{}) error {
	return nil
}

func resourceUpdateItem(d *schema.Resource, m interface{}) error {
	return resourceReadItem(d, m)
}

func resourceDeleteItem(d *schema.Resource, m interface{}) error {
	return nil
}

func resourceExistsItem(d *schema.Resource, m interface{}) error {
	return resourceReadItem(d, m)
}
