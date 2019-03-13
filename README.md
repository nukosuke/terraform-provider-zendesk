# terraform-provider-zendesk
[![Build Status](https://travis-ci.org/nukosuke/terraform-provider-zendesk.svg?branch=master)](https://travis-ci.org/nukosuke/terraform-provider-zendesk)
[![Build status](https://ci.appveyor.com/api/projects/status/ti5il35v6a6ankcq/branch/master?svg=true)](https://ci.appveyor.com/project/nukosuke/terraform-provider-zendesk/branch/master)
[![Coverage Status](https://coveralls.io/repos/github/nukosuke/terraform-provider-zendesk/badge.svg?branch=master)](https://coveralls.io/github/nukosuke/terraform-provider-zendesk?branch=master)

Terraform provider for Zendesk

- [Available resources](https://github.com/nukosuke/terraform-provider-zendesk/wiki)

## Requirements

- Go >= v1.11
- Terraform >= v0.11.10

## Installation

Download latest version from [release page](https://github.com/nukosuke/terraform-provider-zendesk/releases)
and locate the binary `terraform-provider-zendesk(.exe)` to executable path of your system.

### Build from source

```sh
$ git clone git@github.com:nukosuke/terraform-provider-zendesk.git
$ cd terraform-provider-zendesk
$ GO111MODULE=on go mod download
$ go build
```

## License

MIT License

See the file [LICENSE](./LICENSE) for details.
