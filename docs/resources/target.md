---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "zendesk_target Resource - terraform-provider-zendesk"
subcategory: ""
description: |-
  Provides a target resource. (HTTP target is deprecated. See https://support.zendesk.com/hc/en-us/articles/4408826284698 for details.)
---

# zendesk_target (Resource)

Provides a target resource. (HTTP target is deprecated. See https://support.zendesk.com/hc/en-us/articles/4408826284698 for details.)

## Example Usage

```terraform
# [DEPRECATED]
# see https://support.zendesk.com/hc/en-us/articles/4408826284698 for details.
#
# API reference:
#   https://developer.zendesk.com/rest_api/docs/support/targets

resource "zendesk_target" "email-target" {
  title = "target :: email :: john.doe@example.com"
  type = "email_target"

  email = "john.doe@example.com"
  subject = "New ticket created"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `title` (String) A name for the target.
- `type` (String)

### Optional

- `active` (Boolean) Whether or not the target is activated.
- `content_type` (String, Deprecated) Content-Type for http_target
- `email` (String) Email address for "email_target"
- `id` (String) The ID of this resource.
- `method` (String) HTTP method.
- `password` (String) Password of the account which the target authenticate.
- `subject` (String) Email subject for "email_target"
- `target_url` (String) The URL for the target. Some target types commonly use this field.
- `username` (String) Username of the account which the target recognize.

### Read-Only

- `url` (String)


