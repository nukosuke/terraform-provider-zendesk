# API reference:
#   https://developer.zendesk.com/api-reference/event-connectors/webhooks/webhooks/

resource "zendesk_webhook" "example-webhook" {
  name           = "Example Webhook without authentication"
  endpoint       = "https://example.com/status/200"
  http_method    = "GET"
  request_format = "json"
  status         = "active"
  subscriptions  = ["conditional_ticket_events"]
}

resource "zendesk_webhook" "example-api-key-webhook" {
  name           = "Example Webhook with Basic Auth"
  endpoint       = "https://example.com/status/200"
  http_method    = "POST"
  request_format = "xml"
  status         = "inactive"
  subscriptions  = ["conditional_ticket_events"]

  authentication {
    type         = "api_key"
    add_position = "header"
    data = {
      name  = "header_name"
      value = "xxxxxxxxxxx"
    }
  }
}

resource "zendesk_webhook" "example-basic-auth-webhook" {
  name           = "Example Webhook with Basic Auth"
  endpoint       = "https://example.com/status/200"
  http_method    = "POST"
  request_format = "form_encoded"
  status         = "active"
  subscriptions  = ["conditional_ticket_events"]

  authentication {
    type         = "basic_auth"
    add_position = "header"
    data = {
      token    = "xxxxxxxxxxx"
      username = "john.doe"
      password = "password"
    }
  }
}

resource "zendesk_webhook" "example-bearer-token-webhook" {
  name           = "Example Webhook with Bearer token"
  endpoint       = "https://example.com/status/200"
  http_method    = "POST"
  request_format = "json"
  status         = "active"
  subscriptions  = ["conditional_ticket_events"]

  authentication {
    type         = "bearer_token"
    add_position = "header"
    data = {
      token = "xxxxxxxxxx"
    }
  }
}
