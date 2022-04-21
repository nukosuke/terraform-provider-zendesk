terraform {
  required_providers {
    zendesk = {
      source  = "nukosuke/zendesk"
      version = ">= 0.0"
    }
  }
}

provider "zendesk" {
  # configure credentials from enviroment variables
  #
  # export ZENDESK_ACCOUNT="example"
  # export ZENDESK_EMAIL="john.doe@example.com"
  # export ZENDESK_TOKEN="xxxxxxxxxx"
}
