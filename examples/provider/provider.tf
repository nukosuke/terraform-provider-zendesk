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

provider "zendesk" {
  # example.zendesk.com
  account = "example"
  email   = "john.doe@example.com"
  token   = "xxxxxxxxxx"

  # or configure from enviroment variables
  # if you don't want to hardcode the credentials.
  #
  # export ZENDESK_ACCOUNT="example"
  # export ZENDESK_EMAIL="john.doe@example.com"
  # export ZENDESK_TOKEN="xxxxxxxxxx"
}
