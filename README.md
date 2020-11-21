# glsnip

Copy-paste across machines using [GitLab
Snippets](https://docs.gitlab.com/ee/user/snippets.html) as a storage backend.

This is a simple CLI tool inspired by the usability of `pbcopy` and `pbpaste` or `xclip`
but designed to work across machines.

## Installation

TBC

## Configuration

You can configure `glsnip` via a configuration file or via environment
variables. Environment variables will always override configuration file
settings.

### Configuration file

Create a yaml-formatted file in your $HOME directory called `.pbsnip` as
follows:
```yaml
gitlab_url: https://gitlab.com
token: xxxx
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
pbsnip copy <some_file.txt
ls | pbsnip copy
```

Pasting examples:
```shell
pbsnip paste   # paste to STDOUT
pbsnip paste > myfile.txt
pbsnip paste | less
```

