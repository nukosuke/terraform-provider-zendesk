# groups.tf
#   Group example
#
# API reference:
#   https://developer.zendesk.com/rest_api/docs/support/groups

resource "zendesk_group" "support-group" {
  name = "Support"
}

resource "zendesk_group" "developer-group" {
  name = "Developer"
}
