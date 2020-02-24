# LSC Terraform Provider

This provider is to accompany Lumina SDN Controller as a Network Infrastructure as Code demo.

## Requirements

* go => 1.11

This project used Go Modules, so you will need to enable them using `export GO111MODULE=on`, otherwise your go commands (run, build and test) will fail.

## Example

The API is pretty simple, it just stores items which have a name, description and some tags, tags are a slice of strings. Name serves as the id for the Item. 

``` json
provider "lsc-demo" {
  address = "localhost"
  port    = "38181"
  token   = "superSecretToken"
}

resource "netconf_device" "cisco_1" {
  name = "cisco_1"
  ip_address = "127.0.1.2"
  vendor = "Cisco"
  device_type = "Cisco XR"
  os_version = "6.5.3"
}
```