# terraform-provider-zendesk
[![Gitter](https://badges.gitter.im/terraform-provider-zendesk/Lobby.svg)](https://gitter.im/terraform-provider-zendesk/Lobby?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge)
[![Build Status](https://travis-ci.org/nukosuke/terraform-provider-zendesk.svg?branch=master)](https://travis-ci.org/nukosuke/terraform-provider-zendesk)
[![Build status](https://ci.appveyor.com/api/projects/status/ti5il35v6a6ankcq/branch/master?svg=true)](https://ci.appveyor.com/project/nukosuke/terraform-provider-zendesk/branch/master)
[![Coverage Status](https://coveralls.io/repos/github/nukosuke/terraform-provider-zendesk/badge.svg?branch=master)](https://coveralls.io/github/nukosuke/terraform-provider-zendesk?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/nukosuke/terraform-provider-zendesk)](https://goreportcard.com/report/github.com/nukosuke/terraform-provider-zendesk)

Terraform provider for Zendesk

- [Available resources](https://github.com/nukosuke/terraform-provider-zendesk/wiki)

## Requirements

- Terraform >= v0.11.10
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

## Authors
- [nukosuke](https://github.com/nukosuke)
- [tamccall](https://github.com/tamccall)

## License

MIT License

See the file [LICENSE](./LICENSE) for details.

## See also
- [nukosuke/go-zendesk](https://github.com/nukosuke/go-zendesk)
