# API reference:
#   https://developer.zendesk.com/rest_api/docs/support/groups

resource "zendesk_group" "moderator-group" {
  name = "Moderator"
}

resource "zendesk_group" "developer-group" {
  name = "Developer"
}
