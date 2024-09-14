resource "pihole_dns_record" "dns_record" {
  for_each = {for record in var.dns_records: record.domain => record}

  domain = each.value.domain
  ip     = each.value.ip
}

variable "dns_records" {
  default = {
    aa = {
      domain = "aa.com",
      ip     = "10.0.0.1"
    },
  }
}

terraform {
  required_providers {
    pihole = {
      source  = "terraform.local/local/pihole"
      version = "0.0.1"
    }
  }
}