# LSC Terraform Provider

This provider is to accompany Lumina SDN Controller as a Network Infrastructure as Code demo.

## Requirements

* go => 1.11
* terraform => 0.12.21

This project used Go Modules, so you will need to enable them using `export GO111MODULE=on`, otherwise your go commands (run, build and test) will fail.

## Compiling the plugin
``` sh
go build -o terraform-provider-lsc_v0.0.1 && terraform init
```

## Known Issues
* Terraform makes API calls in parallel which LSC can't handle due to locking - Might implement re-try to fix this
* Terraform makes consquent calls too fast after first mount, need to add operational checking to verify mount completion

## Example

The API is pretty simple, it mounts a device to an ODL controller running Netconf.
Name serves as the id for the Item. 

``` go
provider "lsc" {
  address = "http://localhost"
  port    = "38181"
  token   = "Basic YWRtaW46YWRtaW4="
}
resource "lsc_netconf_device" "cisco1" {
  name = "cisco1"
  port = 830
  ip_address = "10.0.100.192"
  username = "root"
  password = "root"
}

// This only creates a preconfig interface
resource "lsc_cisco_interface" "GigabitEthernet_0_0_0_4" {
  device = lsc_netconf_device.cisco1.name
  name = "GigabitEthernet0/0/0/4"
  description = "Terraform Test"
}

// This only creates a preconfig vlan
resource "lsc_cisco_vlan" "GigabitEthernet_0_0_0_4_1" {
  device = lsc_netconf_device.cisco1.name
  name = "GigabitEthernet0/0/0/4.1"
  description = "Terraform Test"
  mtu = 9216
  interface_mode = "l2-transport"
  outer_tag_type = "match-untagged"
  tag_type = "match-dot1q"
  inner_tag = 9
  outer_tag = 2
}
```

## Helpful Tools
* JSon to go struct - https://mholt.github.io/json-to-go/