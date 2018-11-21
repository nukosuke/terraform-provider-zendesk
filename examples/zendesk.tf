# zendesk.tf
#   Zendesk provider config
#
# API document about authentication:
#   https://developer.zendesk.com/rest_api/docs/support/introduction#security-and-authentication
#
# NOTE:
#   v0.0.0 supports only token authentication
#
# (C) 2018 nukosuke <nukosuke@lavabit.com>

provider "zendesk" {
  url   = ""
  email = ""
  token = ""
}
