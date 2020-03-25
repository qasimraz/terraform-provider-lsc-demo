package provider

import (
	"fmt"
	"log"
	"qasimraz/terraform-provider-lsc-demo/api/client"
	"qasimraz/terraform-provider-lsc-demo/api/payload"
	"strconv"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceCiscoL2VPN() *schema.Resource {
	fmt.Print()
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"eviid": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Eviid for l2vpn",
				ForceNew:    true,
			},
			"device": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Device for this interface",
			},
			"interface_1": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "vlan-aware-fxc-attachment-circuit",
			},
			"interface_2": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "vlan-aware-fxc-attachment-circuit",
			},
		},
		Create: resourceCreateCiscoL2VPN,
		Read:   resourceReadCiscoL2VPN,
		Update: resourceCreateCiscoL2VPN,
		Delete: resourceDeleteCiscoL2VPN,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(45 * time.Second),
		}}
}

func resourceCreateCiscoL2VPN(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	var Circuit = []payload.VlanAwareFxcAttachmentCircuit{
		{
			Name: d.Get("interface_1").(string),
		},
		{
			Name: d.Get("interface_2").(string),
		},
	}

	device := payload.CiscoL2VPN{
		Eviid: d.Get("eviid").(int),
		VlanAwareFxcAttachmentCircuits: payload.VlanAwareFxcAttachmentCircuits{
			VlanAwareFxcAttachmentCircuit: Circuit,
		},
	}

	url := payload.NetconfCiscoL2VPNURL(d.Get("device").(string), d.Get("eviid").(int))

	payloadBody, err := payload.NetconfCiscoL2VPNPayload(device)
	if err != nil {
		log.Print("[Error]: ", err)
		return nil
	}

	return resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		err = apiClient.PutNetconf(url, payloadBody)

		if err != nil {
			return resource.RetryableError(fmt.Errorf("Error from controller: %s", err))
		}

		return resource.NonRetryableError(resourceReadCiscoL2VPN(d, m))
	})
}

func resourceReadCiscoL2VPN(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	url := payload.NetconfCiscoL2VPNURL(d.Get("device").(string), d.Get("eviid").(int))

	bodyBytes, err := apiClient.GetNetconf(url)
	if err != nil {
		return err
	}

	device, err := payload.ParseNetconfCiscoL2VPNPayload(bodyBytes)
	if err != nil {
		log.Print("[Error]: ", err)
		return nil
	}

	d.SetId(strconv.Itoa(device.Eviid))
	d.Set("eviid", device.Eviid)
	d.Set("interface_1", device.VlanAwareFxcAttachmentCircuits.VlanAwareFxcAttachmentCircuit[0].Name)
	d.Set("interface_2", device.VlanAwareFxcAttachmentCircuits.VlanAwareFxcAttachmentCircuit[1].Name)
	return nil
}

func resourceDeleteCiscoL2VPN(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	url := payload.NetconfCiscoL2VPNURL(d.Get("device").(string), d.Get("eviid").(int))

	err := apiClient.DeleteNetconf(url)
	if err != nil {
		log.Print("[Error]: ", err)
		return nil
	}

	d.SetId("")
	return nil
}
