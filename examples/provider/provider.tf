terraform {
  required_providers {
    pihole = {
      source = "ryanwholey/pihole"
    }
  }
}

provider "pihole" {
  url      = "https://pihole.domain.com" # PIHOLE_URL
  password = var.pihole_password         # PIHOLE_PASSWORD
}
