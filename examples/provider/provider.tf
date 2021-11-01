terraform {
  required_providers {
    pihole = {
      source = "ryanwholey/pihole"
    }
  }
}

provider "pihole" {
  url      = "https://pihole.ryanwholey.com" # PIHOLE_URL
  password = var.pihole_password             # PIHOLE_PASSWORD
}
