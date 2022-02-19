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

provider "pihole" {
  url = "https://pihole.domain.com" # PIHOLE_URL

  # Experimental, requires Pi-hole Web Interface >= 5.11
  api_token = var.pihole_api_token # PIHOLE_API_TOKEN
}
