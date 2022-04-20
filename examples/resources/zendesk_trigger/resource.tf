# API reference:
#   https://developer.zendesk.com/rest_api/docs/support/triggers

resource "zendesk_trigger" "auto-reply-trigger" {
  title  = "Auto Reply Trigger"
  active = true

  all {
    field    = "role"
    operator = "is"
    value    = "end_user"
  }

  all {
    field    = "update_type"
    operator = "is"
    value    = "Create"
  }

  all {
    field    = "status"
    operator = "is_not"
    value    = "solved"
  }

  action {
    field = "notification_user"
    value = jsonencode([
      "requester_id",
      "Dear my customer",
      "Hi. This message was configured by terraform-provider-zendesk."
    ])
  }
}
