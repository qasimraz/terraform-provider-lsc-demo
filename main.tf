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

# resource "lsc_cisco_interface" "GigabitEthernet0/0/0/4" {
#   device = lsc_netconf_device.cisco1.name
#   interface_name = "GigabitEthernet0/0/0/4"
#   description = "Terraform Test"
# }