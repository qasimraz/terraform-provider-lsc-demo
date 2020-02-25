provider "lsc" {
  address = "http://localhost"
  port    = "38181"
  token   = "Basic YWRtaW46YWRtaW4="
}
resource "lsc_item" "cisco_2" {
  name = "cisco_2"
  port = "830"
  ip_address = "127.0.1.2"
  username = "root"
  password = "root"
}