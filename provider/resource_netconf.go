package provider

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceNetconf() *schema.Resource {
	fmt.Print()
	return &schema.Resource{
		Create: resourceCreateItem,
		Read:   resourceReadItem,
		Update: resourceUpdateItem,
		Delete: resourceDeleteItem,
		
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the resource, also acts as it's unique ID",
			},
			"description": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "A description of an item",
			},
		},

	}
}

// func resourceCreateItem(d *schema.Resource, m interface{}) error {
// 	return nil
// }

// func resourceReadItem(d *schema.Resource, m interface{}) error {
// 	return nil
// }

// func resourceUpdateItem(d *schema.Resource, m interface{}) error {
// 	return nil
// }

// func resourceDeleteItem(d *schema.Resource, m interface{}) error {
// 	return nil
// }
