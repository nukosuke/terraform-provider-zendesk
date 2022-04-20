# API reference:
#   https://developer.zendesk.com/rest_api/docs/support/attachments

variable "logo_file_path" {
  type    = string
  default = "../zendesk/testdata/street.jpg"
}

resource "zendesk_attachment" "logo" {
  file_name = "street.jpg"
  file_path = var.logo_file_path
  file_hash = filesha256(var.logo_file_path)
}
