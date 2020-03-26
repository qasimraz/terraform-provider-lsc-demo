provider "lsc" {
  address = "http://localhost"
  port    = "38181"
  token   = "Basic YWRtaW46YWRtaW4="
}
// Creates a netconf mount
resource "lsc_netconf_device" "cisco1" {
  name = "cisco1"
  port = 830
  ip_address = "10.0.100.192"
  username = "root"
  password = "root"
}
// Creates an interface
resource "lsc_cisco_interface" "GigabitEthernet_0_0_0_4" {
  device = lsc_netconf_device.cisco1.name
  name = "GigabitEthernet0/0/0/4"
  description = "Terraform Test"
}
resource "lsc_cisco_interface" "GigabitEthernet_0_0_0_5" {
  device = lsc_netconf_device.cisco1.name
  name = "GigabitEthernet0/0/0/5"
  description = "Terraform Test"
}
// Creates a vlan
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
resource "lsc_cisco_vlan" "GigabitEthernet_0_0_0_5_1" {
  device = lsc_netconf_device.cisco1.name
  name = "GigabitEthernet0/0/0/5.1"
  description = "Terraform Test"
  mtu = 9216
  interface_mode = "l2-transport"
  outer_tag_type = "match-untagged"
  tag_type = "match-dot1q"
  inner_tag = 9
  outer_tag = 2
}
// Creates an L2VPN 
resource "lsc_cisco_l2vpn" "l2vpn_eviid_9" {
  eviid = 9
  device = lsc_netconf_device.cisco1.name
  interface_1 = lsc_cisco_vlan.GigabitEthernet_0_0_0_4_1.name
  interface_2 = lsc_cisco_vlan.GigabitEthernet_0_0_0_5_1.name
}
