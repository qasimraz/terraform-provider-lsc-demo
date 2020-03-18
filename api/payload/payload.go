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

// NetconfPayload struct
type NetconfPayload struct {
	Node []Netconf `json:"node"`
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

// NetconfMountURL returns netconf mount URL - needs check for empty name?
func NetconfMountURL(name string) string {
	return fmt.Sprintf("restconf/config/network-topology:network-topology/topology/topology-netconf/node/%s", name)
}

// NetconfCiscoInterfaceURL returns netconf cisco interface URL
func NetconfCiscoInterfaceURL(device string, interfaceName string) string {
	return fmt.Sprintf("restconf/config/network-topology:network-topology/topology/topology-netconf/node/%s/yang-ext:mount/Cisco-IOS-XR-ifmgr-cfg:interface-configurations/interface-configuration/pre/%s", device, url.QueryEscape(interfaceName))
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
