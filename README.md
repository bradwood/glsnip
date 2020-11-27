[![Go Report Card](https://goreportcard.com/badge/github.com/bradwood/glsnip)](https://goreportcard.com/report/github.com/bradwood/glsnip)
![Go](https://github.com/bradwood/glsnip/workflows/Go/badge.svg)
![Release](https://github.com/bradwood/glsnip/workflows/Release/badge.svg)
[![Coverage Status](https://coveralls.io/repos/github/bradwood/glsnip/badge.svg)](https://coveralls.io/github/bradwood/glsnip)

![glsnip logo](.github/glsnip-logo.png?raw=true)

Copy-paste across machines using [GitLab
Snippets](https://docs.gitlab.com/ee/user/snippets.html) as a storage backend.

This is a simple CLI tool inspired by the usability of `pbcopy` and `pbpaste` or `xclip`
but designed to work across machines.

## Installation

If you have Go installed `go get github.com/bradwood/glsnip`.

Alternatively, you can download a binary from the [Releases
page](https://github.com/bradwood/glsnip/releases)

## Configuration

You can configure `glsnip` via a configuration file or via environment
variables. Environment variables will always override configuration file
settings.

### Configuration file

Create a YAML-formatted config file (default location `$HOME/.glsnip`). You must
include at least a single server profile YAML block called `default`, like this:

```yaml
---
default:
  gitlab_url: https://url.of.gitlab.server/
  token: USERTOKEN
  clipboard_name: glsnip
...
```

Multiple additional server profile blocks can be added using any block name,
like this:

```yaml
...
work:
  gitlab_url: https://url.of.work.server/
  token: USERTOKENWORK
  clipboard_name: glsnip
...
```

You may also specify an alternative location for the configuration file with the
`--config` flag.

### Environment variables

Instead of using a configuration file, you may set environment variables by
prefixing the key in a configuration file block with `GLSNIP_` and then converting
all alphabetic characters to UPPERCASE. Note that environment variables will
override any configuration specified in the configuration file, regardless of
the profile specified. You may specify a server profile by setting
`GLSNIP_PROFILE`.

You can set environment variables as follows:
```shell
export GLSNIP_GITLAB_URL=https://gitlab.com
export GLSNIP_TOKEN=xxxx
```

## Usage

To specify a non-`default` server profile use the `--profile` flag.

Copying examples:
```shell
glsnip copy <some_file.txt
glsnip copy --profile work <some_file.txt  # copy to Snippet at on "work" GitLab server
ls | glsnip copy
ls | GLSNIP_GITLAB_URL=https://blah.com GLSNIP_TOKEN=xxx glsnip copy
```

Pasting examples:
```shell
glsnip paste   # paste to STDOUT
glsnip paste > myfile.txt
glsnip paste --profile public > myfile.txt  # paste from public GitLab server
glsnip paste | less
```

## Contributions

Contributions are welcome, please feel free to raise a PR or issue.
