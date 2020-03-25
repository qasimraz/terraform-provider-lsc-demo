package payload

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
)

// Netconf struct represents a Netconf Device Details
type Netconf struct {
	Name      string `json:"node-id"`
	IPAddress string `json:"netconf-node-topology:host"`
	Port      int    `json:"netconf-node-topology:port"`
	Username  string `json:"netconf-node-topology:username"`
	Password  string `json:"netconf-node-topology:password"`
}

// NetconfOperational struct represents a Netconf Opertaional Device Details
type NetconfOperational struct {
	Name      string `json:"node-id"`
	IPAddress string `json:"netconf-node-topology:host"`
	Port      int    `json:"netconf-node-topology:port"`
	Status    string `json:"netconf-node-topology:connection-status"`
}

// NetconfPayload struct
type NetconfPayload struct {
	Node []Netconf `json:"node"`
}

// NetconfPayloadOperational struct
type NetconfPayloadOperational struct {
	Node []NetconfOperational `json:"node"`
}

// NetconfMountURL returns netconf mount URL - needs check for empty name?
func NetconfMountURL(name string) string {
	return fmt.Sprintf("restconf/config/network-topology:network-topology/topology/topology-netconf/node/%s", name)
}

// NetconfMountURLOperational returns netconf mount Operational URL - needs check for empty name?
func NetconfMountURLOperational(name string) string {
	return fmt.Sprintf("restconf/operational/network-topology:network-topology/topology/topology-netconf/node/%s", name)
}

// NetconfMountPayload forms a json payload for Netconf Mount
func NetconfMountPayload(device Netconf) (bytes.Buffer, error) {
	payloadBody := NetconfPayload{ // Make into seperate function
		Node: []Netconf{device},
	}

	buf := bytes.Buffer{}
	err := json.NewEncoder(&buf).Encode(payloadBody)
	if err != nil {
		return buf, err
	}
	return buf, nil
}

// ParseNetconfMountPayload parses the json netconf mount payload to a struct
func ParseNetconfMountPayload(bodyBytes []byte) (Netconf, error) {
	item := &NetconfPayload{}
	err := json.Unmarshal(bodyBytes, item)
	if err != nil {
		log.Print("[Error]: ", err)
		return Netconf{}, err
	}

	log.Printf("[DEBUG] Parsed Body: ", item)

	var device Netconf = item.Node[0]
	return device, nil
}

// ParseNetconfOperationalMountPayload parses the json netconf mount payload to a struct
func ParseNetconfOperationalMountPayload(bodyBytes []byte) (NetconfOperational, error) {
	item := &NetconfPayloadOperational{}
	err := json.Unmarshal(bodyBytes, item)
	if err != nil {
		log.Print("[Error]: ", err)
		return NetconfOperational{}, err
	}

	log.Printf("[DEBUG] Parsed Body: ", item)

	var device NetconfOperational = item.Node[0]
	return device, nil
}

// CiscoInterface struct
type CiscoInterface struct {
	Active      string `json:"active"`
	Name        string `json:"interface-name"`
	Description string `json:"description"`
}

// CiscoInterfacePayload struct
type CiscoInterfacePayload struct {
	Node []CiscoInterface `json:"interface-configuration"`
}

// NetconfCiscoInterfaceURL returns netconf cisco interface URL
func NetconfCiscoInterfaceURL(device string, interfaceName string) string {
	return fmt.Sprintf("restconf/config/network-topology:network-topology/topology/topology-netconf/node/%s/yang-ext:mount/Cisco-IOS-XR-ifmgr-cfg:interface-configurations/interface-configuration/pre/%s", device, url.QueryEscape(interfaceName))
}

// NetconfCiscoInterfacePayload forms a json payload for cisco interface
func NetconfCiscoInterfacePayload(device CiscoInterface) (bytes.Buffer, error) {
	payloadBody := CiscoInterfacePayload{ // Make into seperate function
		Node: []CiscoInterface{device},
	}

	buf := bytes.Buffer{}
	err := json.NewEncoder(&buf).Encode(payloadBody)
	if err != nil {
		return buf, err
	}
	return buf, nil
}

// ParseNetconfCiscoInterfacePayload parses json payload for cisco interface to a struct
func ParseNetconfCiscoInterfacePayload(bodyBytes []byte) (CiscoInterface, error) {
	item := &CiscoInterfacePayload{}
	err := json.Unmarshal(bodyBytes, item)
	if err != nil {
		log.Print("[Error]: ", err)
		return CiscoInterface{}, err
	}

	log.Printf("[DEBUG] Parsed Body: ", item)

	var device CiscoInterface = item.Node[0]
	return device, nil
}

type CiscoVlanPayload struct {
	Node []CiscoVlan `json:"interface-configuration"`
}
type Mtu struct {
	Owner string `json:"owner"`
	Mtu   int    `json:"mtu"`
}
type Mtus struct {
	Mtu []Mtu `json:"mtu"`
}
type Encapsulation struct {
	OuterTagType string `json:"outer-tag-type"`
}
type Rewrite struct {
	InnerTagType  string `json:"inner-tag-type"`
	InnerTagValue int    `json:"inner-tag-value"`
	OuterTagType  string `json:"outer-tag-type"`
	RewriteType   string `json:"rewrite-type"`
	OuterTagValue int    `json:"outer-tag-value"`
}
type CiscoIOSXRL2EthInfraCfgEthernetService struct {
	Encapsulation Encapsulation `json:"encapsulation"`
	Rewrite       Rewrite       `json:"rewrite"`
}
type CiscoVlan struct {
	Active                                 string                                 `json:"active"`
	InterfaceName                          string                                 `json:"interface-name"`
	Description                            string                                 `json:"description"`
	Mtus                                   Mtus                                   `json:"mtus"`
	InterfaceModeNonPhysical               string                                 `json:"interface-mode-non-physical"`
	CiscoIOSXRL2EthInfraCfgEthernetService CiscoIOSXRL2EthInfraCfgEthernetService `json:"Cisco-IOS-XR-l2-eth-infra-cfg:ethernet-service"`
}

// NetconfCiscoVlanURL returns netconf cisco interface URL
func NetconfCiscoVlanURL(device string, interfaceName string) string {
	return fmt.Sprintf("restconf/config/network-topology:network-topology/topology/topology-netconf/node/%s/yang-ext:mount/Cisco-IOS-XR-ifmgr-cfg:interface-configurations/interface-configuration/pre/%s", device, url.QueryEscape(interfaceName))
}

// NetconfCiscoVlanPayload forms a json payload for cisco interface
func NetconfCiscoVlanPayload(device CiscoVlan) (bytes.Buffer, error) {
	payloadBody := CiscoVlanPayload{ // Make into seperate function
		Node: []CiscoVlan{device},
	}

	buf := bytes.Buffer{}
	err := json.NewEncoder(&buf).Encode(payloadBody)
	if err != nil {
		return buf, err
	}
	return buf, nil
}

// ParseNetconfCiscoVlanPayload parses json payload for cisco interface to a struct
func ParseNetconfCiscoVlanPayload(bodyBytes []byte) (CiscoVlan, error) {
	item := &CiscoVlanPayload{}
	err := json.Unmarshal(bodyBytes, item)
	if err != nil {
		log.Print("[Error]: ", err)
		return CiscoVlan{}, err
	}

	log.Printf("[DEBUG] Parsed Body: ", item)

	var device CiscoVlan = item.Node[0]
	return device, nil
}

type CiscoL2VPNPayload struct {
	Node []CiscoL2VPN `json:"vlan-aware-flexible-xconnect-service"`
}
type VlanAwareFxcAttachmentCircuit struct {
	Name string `json:"name"`
}
type VlanAwareFxcAttachmentCircuits struct {
	VlanAwareFxcAttachmentCircuit []VlanAwareFxcAttachmentCircuit `json:"vlan-aware-fxc-attachment-circuit"`
}
type CiscoL2VPN struct {
	Eviid                          int                            `json:"eviid"`
	VlanAwareFxcAttachmentCircuits VlanAwareFxcAttachmentCircuits `json:"vlan-aware-fxc-attachment-circuits"`
}

// NetconfCiscoL2VPNURL returns netconf cisco interface URL
func NetconfCiscoL2VPNURL(device string, eviid int) string {
	return fmt.Sprintf("restconf/config/network-topology:network-topology/topology/topology-netconf/node/%s/yang-ext:mount/Cisco-IOS-XR-l2vpn-cfg:l2vpn/database/flexible-xconnect-service-table/vlan-aware-flexible-xconnect-services/vlan-aware-flexible-xconnect-service/%d", device, eviid)
}

// NetconfCiscoL2VPNPayload forms a json payload for cisco interface
func NetconfCiscoL2VPNPayload(device CiscoL2VPN) (bytes.Buffer, error) {
	payloadBody := CiscoL2VPNPayload{ // Make into seperate function
		Node: []CiscoL2VPN{device},
	}

	buf := bytes.Buffer{}
	err := json.NewEncoder(&buf).Encode(payloadBody)
	if err != nil {
		return buf, err
	}
	return buf, nil
}

// ParseNetconfCiscoL2VPNPayload parses json payload for cisco interface to a struct
func ParseNetconfCiscoL2VPNPayload(bodyBytes []byte) (CiscoL2VPN, error) {
	item := &CiscoL2VPNPayload{}
	err := json.Unmarshal(bodyBytes, item)
	if err != nil {
		log.Print("[Error]: ", err)
		return CiscoL2VPN{}, err
	}

	log.Printf("[DEBUG] Parsed Body: ", item)

	var device CiscoL2VPN = item.Node[0]
	return device, nil
}
