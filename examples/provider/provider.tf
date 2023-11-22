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

  # Requires Pi-hole Web Interface >= 5.11.0
  api_token = var.pihole_api_token # PIHOLE_API_TOKEN
}

provider "pihole" {
  url = "https://pihole.domain.com"

  # Pi-hole sets the API token to the admin password hashed twiced via SHA-256
  api_token = sha256(sha256(var.pihole_password))
}
