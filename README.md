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

``` sh
$ docker run --rm -ti -v $(pwd)/examples:/terraform nukosuke/terraform-provider-zendesk init

Initializing provider plugins...

Terraform has been successfully initialized!

You may now begin working with Terraform. Try running "terraform plan" to see
any changes that are required for your infrastructure. All Terraform commands
should now work.

If you ever set or change modules or backend configuration for Terraform,
rerun this command to reinitialize your working directory. If you forget, other
commands will detect it and remind you to do so if necessary.


$ docker run --rm -ti -v $(pwd)/examples:/terraform nukosuke/terraform-provider-zendesk plan

Refreshing Terraform state in-memory prior to plan...
The refreshed state will be used to calculate this plan, but will not be
persisted to local or remote state storage.

data.zendesk_ticket_field.status: Refreshing state...
data.zendesk_ticket_field.group: Refreshing state...
data.zendesk_ticket_field.subject: Refreshing state...
data.zendesk_ticket_field.assignee: Refreshing state...
data.zendesk_ticket_field.description: Refreshing state...

------------------------------------------------------------------------

An execution plan has been generated and is shown below.
Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  + zendesk_attachment.logo
      id:                                   <computed>
      content_type:                         <computed>
      content_url:                          <computed>
      file_hash:                            "Vtptw0XCL7+ShQ8G3/9Q2dsYu3iofOk7J3WqTwzning="
      file_name:                            "street.jpg"
      file_path:                            "../zendesk/testdata/street.jpg"
      inline:                               <computed>
      size:                                 <computed>
      thumbnails.#:                         <computed>
      token:                                <computed>

  ...

  + zendesk_trigger.auto-reply-trigger
      id:                                   <computed>
      action.#:                             "1"
      action.14699407.field:                "notification_user"
      action.14699407.value:                "[\n  \"requester_id\",\n  \"Dear my customer\",\n  \"Hi. This message was configured by terraform-provider-zendesk.\"\n]\n"
      active:                               "true"
      all.#:                                "3"
      all.1027606754.field:                 "update_type"
      all.1027606754.operator:              "is"
      all.1027606754.value:                 "Create"
      all.2406064215.field:                 "role"
      all.2406064215.operator:              "is"
      all.2406064215.value:                 "end_user"
      all.375493961.field:                  "status"
      all.375493961.operator:               "is_not"
      all.375493961.value:                  "solved"
      title:                                "Auto Reply Trigger"


Plan: 16 to add, 0 to change, 0 to destroy.
```


## Authors
- [nukosuke](https://github.com/nukosuke)
- [tamccall](https://github.com/tamccall)

## License

MIT License

See the file [LICENSE](./LICENSE) for details.


[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fnukosuke%2Fterraform-provider-zendesk.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Fnukosuke%2Fterraform-provider-zendesk?ref=badge_large)

## See also
- [nukosuke/go-zendesk](https://github.com/nukosuke/go-zendesk)
