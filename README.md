# Thycotic Secret Server - Terraform Provider

The [Thycotic](https://thycotic.com/) [Secret Server](https://thycotic.com/products/secret-server/) [Terraform](https://www.terraform.io/) Provider allows you to access and reference Secrets in your vault for use in Terraform configurations.

This is a fork of Thycotic's own provider incorporating the changes made by Dan Hale to support Secret Server installations using the Domain-field.
It currently depends on Dan's TSS SDK repository, but will use Thycotic's once Dan's PR is merged, sometime in the future.

## Install via Registry

> Preferred way to install

The latest release can be [downloaded from the terraform registry](https://registry.terraform.io/providers/norskhelsenett/tss/latest). The documentation can be found [here](https://registry.terraform.io/providers/norskhelsenett/tss/latest/docs).

If wish to install straight from source, follow the steps below.

## Install form Source

### Terraform 0.12 and earlier

Extract the specific file for your OS and Architecture to the plugins directory
of the user's profile. You may have to create the directory.

| OS      | Default Path                    |
| ------- | ------------------------------- |
| Linux   | `~/.terraform.d/plugins`        |
| Windows | `%APPDATA%\terraform.d\plugins` |

### Terraform 0.13 and later

Terraform 0.13 uses a different file system layout for 3rd party providers. More information on this can be found [here](https://www.terraform.io/upgrade-guides/0-13.html#new-filesystem-layout-for-local-copies-of-providers). The following folder path will need to be created in the plugins directory of the user's profile.

#### Windows

```text
%APPDATA%\TERRAFORM.D\PLUGINS
└───terraform.nhn.no
    └───norskhelsenett
        └───tss
            └───0.1.2
                └───windows_amd64
```

#### Linux

```text
~/.terraform.d/plugins
└───terraform.nhn.no
    └───norskhelsenett
        └───tss
            └───0.1.2
                ├───linux_amd64
```

## Usage

For Terraform 0.13+, include the `terraform` block in your configuration, or plan, that specifies the provider:

```terraform
terraform {
  required_providers {
    tss = {
      source = "norskhelsenett/tss"
      version = "0.1.2"
    }
  }
}
```

To run the example, create a `terraform.tfvars`:

```json
tss_username   = "my_app_user"
tss_password   = "Passw0rd."
tss_domain     = "foo.bar"
tss_server_url = "https://example/SecretServer"
tss_secret_id  = "1"
```
