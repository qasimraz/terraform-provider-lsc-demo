package provider

import (
	"fmt"
	"log"
	"qasimraz/terraform-provider-lsc-demo/api/client"
	"qasimraz/terraform-provider-lsc-demo/api/payload"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceCiscoInterface() *schema.Resource {
	fmt.Print()
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the interface resource",
				ForceNew:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Description of interface",
			},
			"device": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Device for this interface",
			},
		},
		Create: resourceCreateCiscoInterface,
		Read:   resourceReadCiscoInterface,
		Update: resourceCreateCiscoInterface,
		Delete: resourceDeleteCiscoInterface,
	}
}

func resourceCreateCiscoInterface(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	device := payload.CiscoInterface{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Active:      "pre",
	}

	url := payload.NetconfCiscoInterfaceURL(d.Get("device").(string), d.Get("name").(string))

	payloadBody, err := payload.NetconfCiscoInterfacePayload(device)
	if err != nil {
		log.Print("[Error]: ", err)
		return nil
	}

	err = apiClient.PutNetconf(url, payloadBody)
	if err != nil {
		log.Print("[Error]: ", err)
		return nil
	}

	d.SetId(device.Name)
	return nil
}

func resourceReadCiscoInterface(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	url := payload.NetconfCiscoInterfaceURL(d.Get("device").(string), d.Get("name").(string))

	bodyBytes, err := apiClient.GetNetconf(url)
	if err != nil {
		return err
	}

	device, err := payload.ParseNetconfCiscoInterfacePayload(bodyBytes)
	if err != nil {
		log.Print("[Error]: ", err)
		return nil
	}

	d.SetId(device.Name)
	d.Set("name", device.Name)
	d.Set("description", device.Description)
	return nil
}

func resourceDeleteCiscoInterface(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	url := payload.NetconfCiscoInterfaceURL(d.Get("device").(string), d.Get("name").(string))

	err := apiClient.DeleteNetconf(url)
	if err != nil {
		log.Print("[Error]: ", err)
		return nil
	}

	d.SetId("")
	return nil
}
