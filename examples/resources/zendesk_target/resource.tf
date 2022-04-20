# [DEPRECATED]
# see https://support.zendesk.com/hc/en-us/articles/4408826284698 for details.
#
# API reference:
#   https://developer.zendesk.com/rest_api/docs/support/targets

resource "zendesk_target" "email-target" {
  title = "target :: email :: john.doe@example.com"
  type = "email_target"

  email = "john.doe@example.com"
  subject = "New ticket created"
}
