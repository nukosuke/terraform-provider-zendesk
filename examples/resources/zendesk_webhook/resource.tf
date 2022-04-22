# API reference:
#   https://developer.zendesk.com/api-reference/event-connectors/webhooks/webhooks/

resource "zendesk_webhook" "example-webhook" {
  name           = "Example Webhook without authentication"
  endpoint       = "https://example.com/status/200"
  http_method    = "GET"
  request_format = "json"
  subscriptions  = ["conditional_ticket_events"]
}

resource "zendesk_webhook" "example-basic-auth-webhook" {
  name           = "Example Webhook with Basic Auth"
  endpoint       = "https://example.com/status/200"
  http_method    = "POST"
  request_format = "json"
  subscriptions  = ["conditional_ticket_events"]

  authentication {
    type         = "basic_auth"
    add_position = "header"
    data = jsonencode({
      token    = "xxxxxxxxxxx"
      username = "john.doe"
      password = "password"
    })
  }
}

resource "zendesk_webhook" "example-bearer-token-webhook" {
  name           = "Example Webhook with Bearer token"
  endpoint       = "https://example.com/status/200"
  http_method    = "POST"
  request_format = "json"
  subscriptions  = ["conditional_ticket_events"]

  authentication {
    type         = "bearer_token"
    add_position = "header"
    data = jsonencode({
      token = "xxxxxxxxxx"
    })
  }
}
