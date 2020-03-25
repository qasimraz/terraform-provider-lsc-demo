package provider

import (
	"errors"
	"fmt"
	"log"
	"qasimraz/terraform-provider-lsc-demo/api/client"
	"qasimraz/terraform-provider-lsc-demo/api/payload"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
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
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(45 * time.Second),
		}}
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

	return resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		err = apiClient.PutNetconf(url, payloadBody)

		if err != nil {
			return resource.RetryableError(fmt.Errorf("Error from controller: %s", err))
		}

		return resource.NonRetryableError(resourceReadCiscoInterface(d, m))
	})
}

func resourceReadCiscoInterface(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	url := payload.NetconfCiscoInterfaceURL(d.Get("device").(string), d.Get("name").(string))

	bodyBytes, err := apiClient.GetNetconf(url)
	if err != nil {
		if errors.Is(err, client.ErrNotFound) {
			d.SetId("")
			return nil
		}
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
