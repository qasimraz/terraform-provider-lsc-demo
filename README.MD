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

## Example

The API is pretty simple, it mounts a device to an ODL controller running Netconf.
Name serves as the id for the Item. 

``` go
provider "lsc" {
  address = "http://localhost"
  port    = "38181"
  token   = "superSecretToken"
}

resource "lsc_netconf_device" "cisco1" {
  name = "cisco1"
  port = 830
  ip_address = "123.226.253.52"
  username = "root"
  password = "root"
}
```