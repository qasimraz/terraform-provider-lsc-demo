provider "example" {
  address = "http://localhost"
  port    = "3001"
  token   = "superSecretToken"
}

resource "netconf_device" "test" {
  name = "this_is_an_item"
  description = "this is an item"
}