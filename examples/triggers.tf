# triggers.tf
#   Trigger example
#
# API reference:
#   https://developer.zendesk.com/rest_api/docs/support/triggers
#
# (C) 2018 nukosuke <nukosuke@lavabit.com>

resource "zendesk_trigger" "auto-reply-trigger" {
  title = "Auto Reply Trigger"
}
