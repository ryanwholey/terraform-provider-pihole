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
    command = "until curl -sS ${local.pihole_url}/admin/api.php ; do echo waiting for Pi-hole API && sleep 1 ; done"
  }
}

resource "pihole_cname_record" "record" {
  domain = "foo.com"
  target = "bar.com"

  depends_on = [ null_resource.pihole_wait ]
}
