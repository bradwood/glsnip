[![Go Report Card](https://goreportcard.com/badge/github.com/bradwood/glsnip)](https://goreportcard.com/report/github.com/bradwood/glsnip)

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

Create a yaml-formatted file in your $HOME directory called `.glsnip` as
follows:
```yaml
gitlab_url: https://gitlab.com
token: xxxx
clipboard_name: glsnip
```
You may also specify an alternative location for the configuration file with the
`--config` flag.

### Environment variables

You can set environment variables as follows:
```shell
export GLSNIP_GITLAB_URL=https://gitlab.com
export GLSNIP_TOKEN=xxxx
```

## Usage

Copying examples:
```shell
glsnip copy <some_file.txt
ls | glsnip copy
ls | GLSNIP_GITLAB_URL=https://blah.com GLSNIP_TOKEN=xxx glsnip copy
```

Pasting examples:
```shell
glsnip paste   # paste to STDOUT
glsnip paste > myfile.txt
glsnip paste | less
```

## Contributions

Contributions are welcome, please feel free to raise a PR or issue.
