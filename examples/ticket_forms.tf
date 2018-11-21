# ticket_forms.tf
#   Ticket Forms example
#
# API reference:
#   https://developer.zendesk.com/rest_api/docs/support/ticket_forms
#
# (C) 2018 nukosuke <nukosuke@lavabit.com>

resource "zendesk_ticket_form" "1-form" {
  name = "Form 1"
}
