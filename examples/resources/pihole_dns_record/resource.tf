resource "pihole_dns_record" "record" {
  domain = "foo.com"
  ip     = "127.0.0.1"
}
