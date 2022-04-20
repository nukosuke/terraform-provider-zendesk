# API reference:
#   https://developer.zendesk.com/rest_api/docs/support/sla_policies

resource "zendesk_sla_policy" "incidents_sla_policy" {
  title  = "Incidents"
  active = true

  all {
    field    = "type"
    operator = "is"
    value    = "incident"
  }

  policy_metrics {
    priority = "normal"
    metric = "first_reply_time"
    target = 30
    business_hours = false
  }
}
