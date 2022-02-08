resource "pihole_group" "group" {
  name        = "relaxed"
  description = "A group for clients with more relaxed allow/deny rules"
}
