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
