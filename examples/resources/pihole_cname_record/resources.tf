resource "pihole_cname_record" "record" {
  domain = "foo.com"
  target = "bar.com"
}
