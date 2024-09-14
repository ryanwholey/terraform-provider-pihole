# terraform-provider-pihole

![test workflow status](https://github.com/ryanwholey/terraform-provider-pihole/actions/workflows/test.yml/badge.svg?branch=main) [![terraform registry](https://img.shields.io/badge/terraform-registry-623CE4)](https://registry.terraform.io/providers/ryanwholey/pihole/latest/docs)

[Pi-hole](https://pi-hole.net/) is an ad blocking application which acts as a DNS proxy that returns empty responses when DNS requests for known advertisement domains are made from your devices. It has a number of additional capabilities like optional DHCP server capabilities, specific allow/deny profiles for specific clients, and a neat UI with a ton of information regarding your internet traffic.

Pi-hole is an open source project and can be found at https://github.com/pi-hole/pi-hole.

## Usage

This provider is published to the Terraform Provider registry.

```tf
terraform {
  required_providers {
    pihole = {
      source  = "ryanwholey/pihole"
      version = "x.x.x"
    }
  }
}
```

Configure the provider with credentials, or pass environment variables.

```tf
provider "pihole" {
  url       = "https://pihole.domain.com" # PIHOLE_URL
  password  = var.pihole_password         # PIHOLE_PASSWORD

  # api_token = var.pihole_api_token      # PIHOLE_API_TOKEN (experimental, requires Web Interface >= 5.11)
}
```

See the [provider documentation](https://registry.terraform.io/providers/ryanwholey/pihole/latest/docs) for more details.

## Provider Development

There are a few ways to configure local providers. See the somewhat obscure [Terraform plugin installation documentation](https://www.terraform.io/docs/cli/commands/init.html#plugin-installation) for a potential recommended way.

One way to run a local provider is to build the project, move it to the Terraform plugins directory and then use a `required_providers` block to note the address and version.

> [!NOTE]
> Note the `/darwin_arm64/` path portion targets a Mac with an ARM64 processor,
> see https://github.com/ryanwholey/terraform-provider-pihole/blob/main/.goreleaser.yml#L18-L27
> for possible supported combinations.

```sh
# from the project root
go build .

mkdir -p ~/.terraform.d/plugins/terraform.local/local/pihole/0.0.1/darwin_arm64/

cp terraform-provider-pihole ~/.terraform.d/plugins/terraform.local/local/pihole/0.0.1/darwin_arm64/terraform-provider-pihole_v0.0.1
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

Testing a Terraform provider comes in several forms. This chapter will attempt to explain the differences, where to find documentation, and how to contribute.

> [!NOTE]
> For the current tests in this repository the SDKv2 is used. In issue [#4](https://github.com/ryanwholey/terraform-provider-pihole/issues/38) an upgrade from SDKv2 to [plugin-testing](https://developer.hashicorp.com/terraform/plugin/framework) can be tracked.

#### Unit testing
```sh
make test
```

#### Acceptance testing

The `make testall` command is prefixed with the `TF_ACC=1`. This tells go to include the tests that utilise the `helper/resource.Test()` functions.

For further reading, please see Hashicorp's [documenation](https://developer.hashicorp.com/terraform/plugin/sdkv2/testing/acceptance-tests) on acceptance tests.

To setup a proper environment combining an instance of Pihole in a docker container with tests, some environment variables need to be set for the tests to make their requests to the correct location.

Run the following commands to test against a local Pi-hole server via [docker](https://docs.docker.com/engine/install/)
```sh
# Set the local Terraform provider environment variables
export PIHOLE_URL=http://localhost:8080
export PIHOLE_PASSWORD=test

# Start the pi-hole server
make docker-run

# Run Terraform tests against the server
make testall
```

To test against a specific Pi-hole image tag, specify the tag via the `TAG` env var

```sh
TAG=nightly make docker-run
```

For further reading about Terraform acceptance tests, see Hashicorp's [documenation](https://developer.hashicorp.com/terraform/plugin/sdkv2/testing/acceptance-tests) on acceptance tests.

#### TFTest

To assert that resources are created by the planned result of Terraform, the [Terraform tests chapter](https://developer.hashicorp.com/terraform/language/tests) is a good introduction on the topic.

No such tests are yet implemented.

### Docs

Documentation is auto-generated via [tfplugindocs](https://github.com/hashicorp/terraform-plugin-docs) from description fields within the provider package, as well as examples and templates from the `examples/` and `templates/` folders respectively.

To generate the docs, ensure that `tfplugindocs` is installed, then run

```sh
make docs
```
