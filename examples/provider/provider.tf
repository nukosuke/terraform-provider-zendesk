# API document about authentication:
#   https://developer.zendesk.com/rest_api/docs/support/introduction#security-and-authentication
#
# NOTE:
#   currently only token authentication is supported

terraform {
  required_providers {
    zendesk = {
      source  = "nukosuke/zendesk"
      version = ">= 0.0"
    }
  }
}

variable "ZENDESK_ACCOUNT" { type = string }
variable "ZENDESK_EMAIL"   { type = string }
variable "ZENDESK_TOKEN"   { type = string }

provider "zendesk" {
  account = var.ZENDESK_ACCOUNT
  email   = var.ZENDESK_EMAIL
  token   = var.ZENDESK_TOKEN
}
