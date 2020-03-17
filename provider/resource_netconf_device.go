package provider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"qasimraz/terraform-provider-lsc-demo/api/client"
	"qasimraz/terraform-provider-lsc-demo/api/payload"
	"strings"

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

	payloadBody := payload.NetconfPayload{ // Make into seperate function
		Node: []payload.Netconf{device},
	}

	buf := bytes.Buffer{}
	err := json.NewEncoder(&buf).Encode(payloadBody)
	if err != nil {
		return err
	}

	err = apiClient.PutNetconf(payload.NetconfMountURL(d.Get("name").(string)), buf)

	if err != nil {
		return err
	}
	d.SetId(device.Name)
	return nil
}

func resourceReadNetconfDevice(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	deviceID := d.Id()
	bodyBytes, err := apiClient.GetNetconf(payload.NetconfMountURL(d.Get("name").(string)))
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			d.SetId("")
		} else {
			return fmt.Errorf("error finding Item with ID %s", deviceID)
		}
	}

	item := &payload.NetconfPayload{}
	err = json.Unmarshal(bodyBytes, item)
	if err != nil {
		log.Print("[Error]: ", err)
		return nil
	}

	var device = &item.Node[0]

	log.Printf("[DEBUG] Parsed Body: ", item)

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

	deviceID := d.Id()

	err := apiClient.DeleteNetconf(payload.NetconfMountURL(deviceID))
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}
