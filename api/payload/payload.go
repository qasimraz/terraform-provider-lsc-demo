package payload

import "fmt"

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

// NetconfMountURL returns netconf mount URL - needs check for empty name?
func NetconfMountURL(name string) string {
	return fmt.Sprintf("restconf/config/network-topology:network-topology/topology/topology-netconf/node/%s", name)
}
