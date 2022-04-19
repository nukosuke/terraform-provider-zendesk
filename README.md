# terraform-provider-zendesk
[![Gitter](https://badges.gitter.im/terraform-provider-zendesk/Lobby.svg)](https://gitter.im/terraform-provider-zendesk/Lobby?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge)
[![Actions Status](https://github.com/nukosuke/terraform-provider-zendesk/workflows/CI/badge.svg)](https://github.com/nukosuke/terraform-provider-zendesk/actions)
[![Build status](https://ci.appveyor.com/api/projects/status/ti5il35v6a6ankcq/branch/master?svg=true)](https://ci.appveyor.com/project/nukosuke/terraform-provider-zendesk/branch/master)
[![Coverage Status](https://coveralls.io/repos/github/nukosuke/terraform-provider-zendesk/badge.svg?branch=master)](https://coveralls.io/github/nukosuke/terraform-provider-zendesk?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/nukosuke/terraform-provider-zendesk)](https://goreportcard.com/report/github.com/nukosuke/terraform-provider-zendesk)
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fnukosuke%2Fterraform-provider-zendesk.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2Fnukosuke%2Fterraform-provider-zendesk?ref=badge_shield)

Terraform provider for Zendesk

- [Available resources](https://github.com/nukosuke/terraform-provider-zendesk/wiki)

## Requirements

- Terraform >= v0.12.7
- Go >= v1.16 (only for build)

## Installation

Download latest version from [release page](https://github.com/nukosuke/terraform-provider-zendesk/releases)
and locate the binary `terraform-provider-zendesk(.exe)` to executable path of your system.

### Build from source

```sh
$ git clone git@github.com:nukosuke/terraform-provider-zendesk.git
$ cd terraform-provider-zendesk
$ go mod download
$ go build
```

## Docker

[Docker image](https://hub.docker.com/r/nukosuke/terraform-provider-zendesk) is available. You can execute Terraform commands without preparing Zendesk provider environment by your self.

Mount your resource directory to `/terraform`.

### `terraform init`

```sh
~/workspace/github.com/nukosuke/terraform-provider-zendesk

⟩ docker run --rm -ti \
  -e TF_VAR_ZENDESK_EMAIL=agent@example.com \
  -e TF_VAR_ZENDESK_TOKEN=xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx \
  -e TF_VAR_ZENDESK_ACCOUNT=d3v-terraform-provider \
  -v $(pwd)/examples:/terraform \
  nukosuke/terraform-provider-zendesk init
```

### `terraform plan`

```sh
~/workspace/github.com/nukosuke/terraform-provider-zendesk

⟩ docker run --rm -ti \
  -e TF_VAR_ZENDESK_EMAIL=agent@example.com \
  -e TF_VAR_ZENDESK_TOKEN=xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx \
  -e TF_VAR_ZENDESK_ACCOUNT=d3v-terraform-provider \
  -v $(pwd)/examples:/terraform \
  -v $(pwd)/zendesk/testdata:/zendesk/testdata
  nukosuke/terraform-provider-zendesk plan
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
