provider "lsc" {
  address = "http://localhost"
  port    = "38181"
  token   = "Basic YWRtaW46YWRtaW4="
}
resource "lsc_item" "cisco2" {
  name = "cisco2"
  port = 830
  ip_address = "207.226.253.52"
  username = "root"
  password = "root"
}