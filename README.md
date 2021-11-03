# terraform-provider-pihole

![test workflow status](https://github.com/ryanwholey/terraform-provider-pihole/actions/workflows/test.yml/badge.svg?branch=main) [![terraform registry](https://img.shields.io/badge/terraform-registry-623CE4)](https://registry.terraform.io/providers/ryanwholey/pihole/latest/docs)

[Pi-hole](https://pi-hole.net/) is an ad blocking application which acts as a DNS proxy that returns empty responses when DNS requests for known advertisement domains are made from your devices. It has a number of additional capabilities like optional DHCP server capabilities, specific allow/deny profiles for specific clients, and a neat UI with a ton of information regarding your internet traffic.

Pi-hole is an open source project and can be found at https://github.com/pi-hole/pi-hole.

## Provider Development

There are a few ways to configure local providers. See the somewhat obscure [Terraform plugin installation documentation](https://www.terraform.io/docs/cli/commands/init.html#plugin-installation) for a potential recommended way. 

One way to run a local provider is to build the project, move it to the Terraform plugins directory and then use a `required_providers` block to note the address and version.

```sh
# from the project root
go build .

# Note the `/darwin_amd64/` path portion targets a Mac with an AMD64 processor, 
# see https://github.com/ryanwholey/terraform-provider-pihole/blob/main/.goreleaser.yml#L18-L27
# for possible supported combinations

mkdir -p ~/.terraform.d/plugins/terraform.local/local/pihole/0.0.1/darwin_amd64/

cp terraform-provider-pihole ~/.terraform.d/plugins/terraform.local/local/pihole/0.0.1/darwin_amd64/terraform-provider-pihole_v0.0.1
```

In the Terraform workspace, use a `required_providers` block to target the locally built provider

```tf
terraform {
  required_providers {
    pihole = {
      source  = "terraform.local/local/pihole"
      version = "0.0.1"
    }
  }
}
```

### Testing

Unit tests can be ran with a simple command

```sh
make test
```

Acceptance can run against any Pi-hole deployment given that `PIHOLE_URL` and `PIHOLE_PASSWORD` are set in the shell. A dockerized Pi-hole can be ran via the docker-compose file provided in the project root.

```sh
# from the project root
docker-compose up -d --build

export PIHOLE_URL=http://localhost:8080
export PIHOLE_PASSWORD=test

make testall
```

### Docs

Documentation is auto-generated via [tfplugindocs](https://github.com/hashicorp/terraform-plugin-docs) from description fields within the provider package, as well as examples and templates from the `examples/` and `templates/` folders respectively. 

To generate the docs, ensure that `tfplugindocs` is installed, then run

```sh
make docs
```
