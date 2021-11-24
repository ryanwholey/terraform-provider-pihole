# Return all domains registered with pi-hole
data "pihole_domains" "all" {}

# Return all denied (blacklisted) domains registered with pihole
data "pihole_domains" "denied" {
  type = "deny"
}
