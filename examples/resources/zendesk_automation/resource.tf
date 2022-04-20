# API reference:
#   https://developer.zendesk.com/rest_api/docs/support/automations

resource "zendesk_automation" "auto-close-automation" {
  title  = "Close ticket 4 days after status is set to solved"
  active = true

  all {
    field = "status"
    operator = "is"
    value = "solved"
  }

  all {
    field = "SOLVED"
    operator = "greater_than"
    value = "96"
  }

  action {
    field = "status"
    value = "closed"
  }
}
