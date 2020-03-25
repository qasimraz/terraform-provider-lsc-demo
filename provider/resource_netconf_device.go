package provider

import (
	"errors"
	"fmt"
	"log"
	"qasimraz/terraform-provider-lsc-demo/api/client"
	"qasimraz/terraform-provider-lsc-demo/api/payload"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceNetconfDevice() *schema.Resource {
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
				Type:        schema.TypeInt,
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
		Create: resourceCreateNetconfDevice,
		Read:   resourceReadNetconfDevice,
		Update: resourceCreateNetconfDevice,
		Delete: resourceDeleteNetconfDevice,
	}
}

func resourceCreateNetconfDevice(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	device := payload.Netconf{
		Name:      d.Get("name").(string),
		Port:      d.Get("port").(int),
		IPAddress: d.Get("ip_address").(string),
		Username:  d.Get("username").(string),
		Password:  d.Get("password").(string),
	}

	url := payload.NetconfMountURL(d.Get("name").(string))

	payloadBody, err := payload.NetconfMountPayload(device)
	if err != nil {
		log.Print("[Error]: ", err)
		return nil
	}

	err = apiClient.PutNetconf(url, payloadBody)
	if err != nil {
		log.Print("[Error]: ", err)
		return nil
	}

	// Verify netconf mount connects succesfully
	for {
		time.Sleep(5 * time.Second)
		url := payload.NetconfMountURLOperational(d.Get("name").(string))
		bodyBytes, err := apiClient.GetNetconf(url)
		if err != nil {
			log.Print("[Error]: ", err)
			return nil
		}

		device, err := payload.ParseNetconfOperationalMountPayload(bodyBytes)
		if err != nil {
			log.Print("[Error]: ", err)
			return nil
		}
		log.Print("[Status]: ", device.Status)
		if device.Status == "connected" {
			break
		}
	}

	d.SetId(device.Name)
	return nil
}

func resourceReadNetconfDevice(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	url := payload.NetconfMountURL(d.Get("name").(string))

	bodyBytes, err := apiClient.GetNetconf(url)
	if err != nil {
		if errors.Is(err, client.ErrNotFound) {
			d.SetId("")
			return nil
		}
		return err
	}

	device, err := payload.ParseNetconfMountPayload(bodyBytes)
	if err != nil {
		log.Print("[Error]: ", err)
		return nil
	}

	d.SetId(device.Name)
	d.Set("name", device.Name)
	d.Set("port", device.Port)
	d.Set("ip_address", device.IPAddress)
	d.Set("username", device.Username)
	d.Set("password", device.Password)
	return nil
}

func resourceDeleteNetconfDevice(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	url := payload.NetconfMountURL(d.Get("name").(string))

	err := apiClient.DeleteNetconf(url)
	if err != nil {
		log.Print("[Error]: ", err)
		return nil
	}

	d.SetId("")
	return nil
}
