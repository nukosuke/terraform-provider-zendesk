# API reference:
#   https://developer.zendesk.com/api-reference/ticketing/ticket-management/dynamic_content/

resource "zendesk_dynamic_content_item" "loc-lang" {
  name           = "language"
  default_locale = "en-US"

  variant {
    locale  = "en-US"
    content = "English (US)"
  }

  variant {
    locale  = "ja"
    content = "日本語"
    default = true
  }

  variant {
    locale  = "zh-tw"
    content = "繁體中文"
  }
}
