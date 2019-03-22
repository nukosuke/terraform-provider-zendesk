# brands.tf
#   Brand example
#
# API reference:
#   https://developer.zendesk.com/rest_api/docs/support/brands

resource "zendesk_brand" "T-800" {
  name            = "T-800"
  active          = true
  has_help_center = true
  subdomain       = "d3v-terraform-provider-t800"
  # TODO: logo
}

resource "zendesk_brand" "T-1000" {
  name            = "T-1000"
  active          = false
  has_help_center = false
  subdomain       = "d3v-terraform-provider-t1000"
}
