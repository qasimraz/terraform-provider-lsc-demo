package provider

import (
	"fmt"
	"qasimraz/terraform-provider-lsc-demo/api/client"
	"strings"

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
			"ip_address": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "IP Address of Netconf Device",
			},
			"port": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Port of the Netconf Device, Default is 830",
			},
			"username": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Username to authenticate to the device",
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Password to authenticate to the device",
			},
		},
		Create: resourceCreateItem,
		Read:   resourceReadItem,
		Update: resourceCreateItem,
		Delete: resourceDeleteItem,
	}
}

func resourceCreateItem(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	device := client.Netconf{
		Name:      d.Get("name").(string),
		Port:      d.Get("port").(string),
		IPAddress: d.Get("ip_address").(string),
		Username:  d.Get("username").(string),
		Password:  d.Get("password").(string),
	}

	err := apiClient.NetconfMount(&device)

	if err != nil {
		return err
	}
	d.SetId(device.Name)
	return nil
}

func resourceReadItem(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	deviceID := d.Id()
	device, err := apiClient.GetNetconfMount(deviceID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			d.SetId("")
		} else {
			return fmt.Errorf("error finding Item with ID %s", deviceID)
		}
	}

	d.SetId(device.Name)
	d.Set("name", device.Name)
	d.Set("port", device.Port)
	d.Set("ip_address", device.IPAddress)
	d.Set("username", device.Username)
	d.Set("password", device.Password)
	return nil
}

// func resourceUpdateItem(d *schema.ResourceData, m interface{}) error {
// 	return nil
// }

func resourceDeleteItem(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	deviceID := d.Id()

	err := apiClient.DeleteNetconfMount(deviceID)
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}
