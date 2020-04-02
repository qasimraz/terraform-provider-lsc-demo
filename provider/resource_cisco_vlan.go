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

func resourceCiscoVlan() *schema.Resource {
	fmt.Print()
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the vlan resource",
				ForceNew:    true,
			},
			"interface": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "interface for dependency",
			},
			"description": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Description of vlan",
			},
			"device": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Device for this vlan",
			},
			"mtu": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "MTU size",
			},
			"interface_mode": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "interface-mode-non-physical, ie l2-transport",
			},
			"outer_tag_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "outer-tag-type for this vlan ie match-untagged",
			},
			"tag_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "inner-tag-type and outer-tag-type ie match-dot1",
			},
			"inner_tag": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "outer-tag-type for this vlan ie match-untagged",
			},
			"outer_tag": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "outer-tag-type for this vlan ie match-untagged",
			},
		},
		Create: resourceCreateCiscoVlan,
		Read:   resourceReadCiscoVlan,
		Update: resourceCreateCiscoVlan,
		Delete: resourceDeleteCiscoVlan,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(45 * time.Second),
		}}
}

func resourceCreateCiscoVlan(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	var MTUbody = []payload.Mtu{
		{
			Owner: "sub_vlan",
			Mtu:   d.Get("mtu").(int),
		},
	}

	device := payload.CiscoVlan{
		InterfaceName:            d.Get("name").(string),
		Description:              d.Get("description").(string),
		Active:                   "pre",
		InterfaceModeNonPhysical: d.Get("interface_mode").(string),
		CiscoIOSXRL2EthInfraCfgEthernetService: payload.CiscoIOSXRL2EthInfraCfgEthernetService{
			Encapsulation: payload.Encapsulation{
				OuterTagType: d.Get("outer_tag_type").(string),
			},
			Rewrite: payload.Rewrite{
				InnerTagType:  d.Get("tag_type").(string),
				OuterTagType:  d.Get("tag_type").(string),
				InnerTagValue: d.Get("inner_tag").(int),
				OuterTagValue: d.Get("outer_tag").(int),
				RewriteType:   "push2",
			},
		},
		Mtus: payload.Mtus{
			Mtu: MTUbody,
		},
	}

	url := payload.NetconfCiscoVlanURL(d.Get("device").(string), d.Get("name").(string))

	payloadBody, err := payload.NetconfCiscoVlanPayload(device)
	if err != nil {
		log.Print("[Error]: ", err)
		return nil
	}

	return resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		err = apiClient.PutNetconf(url, payloadBody)

		if err != nil {
			return resource.RetryableError(fmt.Errorf("Error from controller: %s", err))
		}

		return resource.NonRetryableError(resourceReadCiscoVlan(d, m))
	})
}

func resourceReadCiscoVlan(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	url := payload.NetconfCiscoVlanURL(d.Get("device").(string), d.Get("name").(string))

	bodyBytes, err := apiClient.GetNetconf(url)
	if err != nil {
		log.Print("[Error] GET: ", err)
		if errors.Is(err, client.ErrNotFound) {
			d.SetId("")
			return nil
		}
		return err
	}

	device, err := payload.ParseNetconfCiscoVlanPayload(bodyBytes)
	if err != nil {
		log.Print("[Error]: ", err)
		return nil
	}

	d.SetId(device.InterfaceName)
	d.Set("name", device.InterfaceName)
	d.Set("description", device.Description)
	d.Set("mtu", device.Mtus.Mtu[0].Mtu)
	d.Set("interface_mode", device.InterfaceModeNonPhysical)
	d.Set("description", device.Description)
	d.Set("outer_tag_type", device.CiscoIOSXRL2EthInfraCfgEthernetService.Encapsulation.OuterTagType)
	d.Set("tag_type", device.CiscoIOSXRL2EthInfraCfgEthernetService.Rewrite.InnerTagType)
	d.Set("inner_tag", device.CiscoIOSXRL2EthInfraCfgEthernetService.Rewrite.InnerTagValue)
	d.Set("outer_tag", device.CiscoIOSXRL2EthInfraCfgEthernetService.Rewrite.OuterTagValue)
	return nil
}

func resourceDeleteCiscoVlan(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	url := payload.NetconfCiscoVlanURL(d.Get("device").(string), d.Get("name").(string))

	err := apiClient.DeleteNetconf(url)
	if err != nil {
		log.Print("[Error]: ", err)
		return nil
	}

	d.SetId("")
	return nil
}
