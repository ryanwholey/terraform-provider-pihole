---
layout: ""
page_title: "Provider: Pi-hole"
description: A Terraform provider to manage Pi-hole resources
---

# Pi-hole Provider

The [Pi-hole](https://pi-hole.net) provider is used to manage Pi-hole resources. The provider should be configured with the Pi-hole URL and the admin password (not the hashed web password).

Use the navigation to the left to read about the available resources.

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `ca_file` (String) CA file to connect to Pi-hole with TLS
- `password` (String) The admin password used to login to the admin dashboard.
- `url` (String) URL where Pi-hole is deployed

## Example Usage

### Basic

```terraform
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
```

**Note**: Authenticating via `api_token` is currently experimental and requires a Pi-hole Web Interface version of `>= 5.11.0` (see [release notes](https://github.com/pi-hole/AdminLTE/releases/tag/v5.11)). The Pi-hole API has just recently began supporting API token authentication for specific resources. Currently the following resources are manageable via API token:

- `pihole_cname_record`
- `pihole_dns_record`

### Dynamic Provider

In the case that Pi-hole is deployed in the same root module that the provider is to be used, a `null_resource` can be used to wait for the server to become ready.

```terraform
provider "docker" {
  host = "unix:///var/run/docker.sock"
}

resource "docker_image" "pihole" {
  name = "pihole/pihole:2022.05"
}

locals {
  pihole_password = "test"
  pihole_url      = "http://${docker_container.pihole.ports[0].ip}:${docker_container.pihole.ports[0].external}"
}

resource "docker_container" "pihole" {
  image = docker_image.pihole.image_id
  name  = "pihole"
  env   = ["WEBPASSWORD=${local.pihole_password}"]

  capabilities {
    add = ["NET_ADMIN"]
  }

  ports {
    internal = 80
    external = 8080
  }
}

provider "pihole" {
  url       = local.pihole_url
  api_token = sha256(sha256(local.pihole_password))
}

resource "null_resource" "pihole_wait" {
  triggers = {
    container = docker_container.pihole.id
  }

  provisioner "local-exec" {
    command = "until curl -sS ${local.pihole_url}/admin/api.php 1>/dev/null ; do echo waiting for Pi-hole API && sleep 1 ; done"
  }
}

resource "pihole_cname_record" "record" {
  domain = "foo.com"
  target = "bar.com"

  depends_on = [ null_resource.pihole_wait ]
}
```
