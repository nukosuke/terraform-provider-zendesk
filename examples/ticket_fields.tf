# ticket_fields.tf
#   Ticket Field example
#
# API reference:
#   https://developer.zendesk.com/rest_api/docs/support/ticket_fields
#
# (C) 2018 nukosuke <nukosuke@lavabit.com>

resource "zendesk_ticket_field" "checkbox-field" {
  title = "Checkbox Field"
  type = "checkbox"
}

resource "zendesk_ticket_field" "date-field" {
  title = "Date Field"
  type = "date"
}

resource "zendesk_ticket_field" "decimal-field" {
  title = "Decimal Field"
  type = "decimal"
}

resource "zendesk_ticket_field" "integer-field" {
  title = "Integer Field"
  type = "integer"
}

resource "zendesk_ticket_field" "regexp-field" {
  title = "Regexp Field"
  type = "regexp"
}

resource "zendesk_ticket_field" "tagger-field" {
  title = "Tagger Field"
  type = "tagger"
  custom_field_options = []
}

resource "zendesk_ticket_field" "text-field" {
  title = "Text Field"
  type = "text"
}

resource "zendesk_ticket_field" "textarea-field" {
  title = "Textarea Field"
  type = "textarea"
}
