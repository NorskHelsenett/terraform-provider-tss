terraform {
  required_version = ">= 0.12.20"
  required_providers {
    tss = {
      source  = "norskhelsenett/tss"
      version = "~> 0.3.0"
    }
  }
}

variable "tss_username" {
  type = string
}

variable "tss_domain" {
  type = string
}

variable "tss_password" {
  type = string
}

variable "tss_server_url" {
  type = string
}

variable "tss_secret_id" {
  type = string
}

provider "tss" {
  username   = var.tss_username
  password   = var.tss_password
  domain     = var.tss_domain
  server_url = var.tss_server_url
}

data "tss_secret" "my_username" {
  id    = var.tss_secret_id
  field = "username"
}

data "tss_secret" "my_password" {
  id    = var.tss_secret_id
  field = "password"
}

output "username" {
  value     = data.tss_secret.my_username.value
}

output "password" {
  value     = data.tss_secret.my_password.value
}
