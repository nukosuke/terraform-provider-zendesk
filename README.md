# terraform-provider-zendesk
[![Gitter](https://badges.gitter.im/terraform-provider-zendesk/Lobby.svg)](https://gitter.im/terraform-provider-zendesk/Lobby?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge)
[![Build Status](https://travis-ci.org/nukosuke/terraform-provider-zendesk.svg?branch=master)](https://travis-ci.org/nukosuke/terraform-provider-zendesk)
[![Build status](https://ci.appveyor.com/api/projects/status/ti5il35v6a6ankcq/branch/master?svg=true)](https://ci.appveyor.com/project/nukosuke/terraform-provider-zendesk/branch/master)
[![Coverage Status](https://coveralls.io/repos/github/nukosuke/terraform-provider-zendesk/badge.svg?branch=master)](https://coveralls.io/github/nukosuke/terraform-provider-zendesk?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/nukosuke/terraform-provider-zendesk)](https://goreportcard.com/report/github.com/nukosuke/terraform-provider-zendesk)
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fnukosuke%2Fterraform-provider-zendesk.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2Fnukosuke%2Fterraform-provider-zendesk?ref=badge_shield)

Terraform provider for Zendesk

- [Available resources](https://github.com/nukosuke/terraform-provider-zendesk/wiki)

## Requirements

- Terraform >= v0.12.7
- Go >= v1.11 (only for build)

## Installation

Download latest version from [release page](https://github.com/nukosuke/terraform-provider-zendesk/releases)
and locate the binary `terraform-provider-zendesk(.exe)` to executable path of your system.

### Build from source

```sh
$ git clone git@github.com:nukosuke/terraform-provider-zendesk.git
$ cd terraform-provider-zendesk
$ export GO111MODULE=on
$ go mod download
$ go build
```

## Docker

[Docker image](https://hub.docker.com/r/nukosuke/terraform-provider-zendesk) is available. You can execute Terraform commands without preparing Zendesk provider environment by your self.

Mount your resource directory to `/terraform`.

<details>
  <summary><b>terraform init</b></summary>
  
```sh
~/workspace/github.com/nukosuke/terraform-provider-zendesk
⟩ docker run --rm -ti \
  -e TF_VAR_ZENDESK_EMAIL=agent@example.com \
  -e TF_VAR_ZENDESK_TOKEN=xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx \
  -e TF_VAR_ZENDESK_ACCOUNT=d3v-terraform-provider \
  -v (pwd)/examples:/terraform \
  nukosuke/terraform-provider-zendesk init

Initializing the backend...

Initializing provider plugins...

Terraform has been successfully initialized!

You may now begin working with Terraform. Try running "terraform plan" to see
any changes that are required for your infrastructure. All Terraform commands
should now work.

If you ever set or change modules or backend configuration for Terraform,
rerun this command to reinitialize your working directory. If you forget, other
commands will detect it and remind you to do so if necessary.
```

</details>

<details>
  <summary><b>terraform plan</b></summary>

```sh
~/workspace/github.com/nukosuke/terraform-provider-zendesk
⟩ docker run --rm -ti \
  -e TF_VAR_ZENDESK_EMAIL=agent@example.com \
  -e TF_VAR_ZENDESK_TOKEN=xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx \
  -e TF_VAR_ZENDESK_ACCOUNT=d3v-terraform-provider \
  -v $(pwd)/examples:/terraform \
  -v $(pwd)/zendesk/testdata:/zendesk/testdata
  nukosuke/terraform-provider-zendesk plan

Refreshing Terraform state in-memory prior to plan...
The refreshed state will be used to calculate this plan, but will not be
persisted to local or remote state storage.

data.zendesk_ticket_field.assignee: Refreshing state...
data.zendesk_ticket_field.description: Refreshing state...
data.zendesk_ticket_field.status: Refreshing state...
data.zendesk_ticket_field.group: Refreshing state...
data.zendesk_ticket_field.subject: Refreshing state...

------------------------------------------------------------------------

An execution plan has been generated and is shown below.
Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # zendesk_attachment.logo will be created
  + resource "zendesk_attachment" "logo" {
      + content_type = (known after apply)
      + content_url  = (known after apply)
      + file_hash    = "56da6dc345c22fbf92850f06dfff50d9db18bb78a87ce93b2775aa4f0ce78a78"
      + file_name    = "street.jpg"
      + file_path    = "../zendesk/testdata/street.jpg"
      + id           = (known after apply)
      + inline       = (known after apply)
      + size         = (known after apply)
      + thumbnails   = (known after apply)
      + token        = (known after apply)
    }

  # zendesk_brand.T-1000 will be created
  + resource "zendesk_brand" "T-1000" {
      + active            = false
      + brand_url         = (known after apply)
      + has_help_center   = (known after apply)
      + help_center_state = (known after apply)
      + id                = (known after apply)
      + name              = "T-1000"
      + subdomain         = "d3v-terraform-provider-t1000"
      + ticket_form_ids   = (known after apply)
      + url               = (known after apply)
    }

  # zendesk_brand.T-800 will be created
  + resource "zendesk_brand" "T-800" {
      + active            = true
      + brand_url         = (known after apply)
      + has_help_center   = (known after apply)
      + help_center_state = (known after apply)
      + id                = (known after apply)
      + name              = "T-800"
      + subdomain         = "d3v-terraform-provider-t800"
      + ticket_form_ids   = (known after apply)
      + url               = (known after apply)
    }

  # zendesk_group.developer-group will be created
  + resource "zendesk_group" "developer-group" {
      + id   = (known after apply)
      + name = "Developer"
      + url  = (known after apply)
    }

  # zendesk_group.moderator-group will be created
  + resource "zendesk_group" "moderator-group" {
      + id   = (known after apply)
      + name = "Moderator"
      + url  = (known after apply)
    }

  # zendesk_target.email-target will be created
  + resource "zendesk_target" "email-target" {
      + active  = true
      + email   = "john.doe@example.com"
      + id      = (known after apply)
      + subject = "New ticket created"
      + title   = "target :: email :: john.doe@example.com"
      + type    = "email_target"
      + url     = (known after apply)
    }

  # zendesk_ticket_field.checkbox-field will be created
  + resource "zendesk_ticket_field" "checkbox-field" {
      + active                = true
      + description           = (known after apply)
      + id                    = (known after apply)
      + position              = (known after apply)
      + regexp_for_validation = (known after apply)
      + removable             = (known after apply)
      + system_field_options  = (known after apply)
      + title                 = "Checkbox Field"
      + title_in_portal       = (known after apply)
      + type                  = "checkbox"
      + url                   = (known after apply)
    }

  # zendesk_ticket_field.date-field will be created
  + resource "zendesk_ticket_field" "date-field" {
      + active                = true
      + description           = (known after apply)
      + id                    = (known after apply)
      + position              = (known after apply)
      + regexp_for_validation = (known after apply)
      + removable             = (known after apply)
      + system_field_options  = (known after apply)
      + title                 = "Date Field"
      + title_in_portal       = (known after apply)
      + type                  = "date"
      + url                   = (known after apply)
    }

  # zendesk_ticket_field.decimal-field will be created
  + resource "zendesk_ticket_field" "decimal-field" {
      + active                = true
      + description           = (known after apply)
      + id                    = (known after apply)
      + position              = (known after apply)
      + regexp_for_validation = (known after apply)
      + removable             = (known after apply)
      + system_field_options  = (known after apply)
      + title                 = "Decimal Field"
      + title_in_portal       = (known after apply)
      + type                  = "decimal"
      + url                   = (known after apply)
    }

  # zendesk_ticket_field.integer-field will be created
  + resource "zendesk_ticket_field" "integer-field" {
      + active                = true
      + description           = (known after apply)
      + id                    = (known after apply)
      + position              = (known after apply)
      + regexp_for_validation = (known after apply)
      + removable             = (known after apply)
      + system_field_options  = (known after apply)
      + title                 = "Integer Field"
      + title_in_portal       = (known after apply)
      + type                  = "integer"
      + url                   = (known after apply)
    }

  # zendesk_ticket_field.regexp-field will be created
  + resource "zendesk_ticket_field" "regexp-field" {
      + active                = true
      + description           = (known after apply)
      + id                    = (known after apply)
      + position              = (known after apply)
      + regexp_for_validation = "^[0-9]+-[0-9]+-[0-9]+$"
      + removable             = (known after apply)
      + system_field_options  = (known after apply)
      + title                 = "Regexp Field"
      + title_in_portal       = (known after apply)
      + type                  = "regexp"
      + url                   = (known after apply)
    }

  # zendesk_ticket_field.tagger-field will be created
  + resource "zendesk_ticket_field" "tagger-field" {
      + active                = true
      + description           = (known after apply)
      + id                    = (known after apply)
      + position              = (known after apply)
      + regexp_for_validation = (known after apply)
      + removable             = (known after apply)
      + system_field_options  = (known after apply)
      + title                 = "Tagger Field"
      + title_in_portal       = (known after apply)
      + type                  = "tagger"
      + url                   = (known after apply)

      + custom_field_option {
          + id    = (known after apply)
          + name  = "Option 1"
          + value = "opt1"
        }
      + custom_field_option {
          + id    = (known after apply)
          + name  = "Option 2"
          + value = "opt2"
        }
    }

  # zendesk_ticket_field.text-field will be created
  + resource "zendesk_ticket_field" "text-field" {
      + active                = true
      + description           = (known after apply)
      + id                    = (known after apply)
      + position              = (known after apply)
      + regexp_for_validation = (known after apply)
      + removable             = (known after apply)
      + system_field_options  = (known after apply)
      + title                 = "Text Field"
      + title_in_portal       = (known after apply)
      + type                  = "text"
      + url                   = (known after apply)
    }

  # zendesk_ticket_field.textarea-field will be created
  + resource "zendesk_ticket_field" "textarea-field" {
      + active                = true
      + description           = (known after apply)
      + id                    = (known after apply)
      + position              = (known after apply)
      + regexp_for_validation = (known after apply)
      + removable             = (known after apply)
      + system_field_options  = (known after apply)
      + title                 = "Textarea Field"
      + title_in_portal       = (known after apply)
      + type                  = "textarea"
      + url                   = (known after apply)
    }

  # zendesk_ticket_form.form-1 will be created
  + resource "zendesk_ticket_form" "form-1" {
      + active               = true
      + id                   = (known after apply)
      + in_all_brands        = true
      + name                 = "Form 1"
      + restricted_brand_ids = (known after apply)
      + ticket_field_ids     = (known after apply)
      + url                  = (known after apply)
    }

  # zendesk_ticket_form.form-2 will be created
  + resource "zendesk_ticket_form" "form-2" {
      + active               = true
      + id                   = (known after apply)
      + in_all_brands        = true
      + name                 = "Form 2"
      + restricted_brand_ids = (known after apply)
      + ticket_field_ids     = (known after apply)
      + url                  = (known after apply)
    }

  # zendesk_trigger.auto-reply-trigger will be created
  + resource "zendesk_trigger" "auto-reply-trigger" {
      + active   = true
      + id       = (known after apply)
      + position = (known after apply)
      + title    = "Auto Reply Trigger"

      + action {
          + field = "notification_user"
          + value = jsonencode(
                [
                  + "requester_id",
                  + "Dear my customer",
                  + "Hi. This message was configured by terraform-provider-zendesk.",
                ]
            )
        }

      + all {
          + field    = "role"
          + operator = "is"
          + value    = "end_user"
        }
      + all {
          + field    = "status"
          + operator = "is_not"
          + value    = "solved"
        }
      + all {
          + field    = "update_type"
          + operator = "is"
          + value    = "Create"
        }
    }

Plan: 17 to add, 0 to change, 0 to destroy.

------------------------------------------------------------------------

Note: You didn't specify an "-out" parameter to save this plan, so Terraform
can't guarantee that exactly these actions will be performed if
"terraform apply" is subsequently run.
```

</details>


## Authors
- [nukosuke](https://github.com/nukosuke)
- [tamccall](https://github.com/tamccall)

## License

MIT License

See the file [LICENSE](./LICENSE) for details.


[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fnukosuke%2Fterraform-provider-zendesk.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Fnukosuke%2Fterraform-provider-zendesk?ref=badge_large)

## See also
- [nukosuke/go-zendesk](https://github.com/nukosuke/go-zendesk)
