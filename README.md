# mracek 

[![Release](https://img.shields.io/github/release/mchlumsky/mracek.svg)](https://github.com/mchlumsky/mracek/releases/latest)
[![codecov](https://codecov.io/gh/mchlumsky/mracek/branch/main/graph/badge.svg?token=YHCWIP3V43)](https://codecov.io/gh/mchlumsky/mracek)
[![Software License](https://img.shields.io/badge/license-MIT-brightgreen.svg)](/LICENSE.md)
[![Build status](https://img.shields.io/github/workflow/status/mchlumsky/mracek/build)](https://github.com/mchlumsky/mracek/actions?workflow=build)
[![Powered By: GoReleaser](https://img.shields.io/badge/powered%20by-goreleaser-green.svg)](https://github.com/goreleaser)

mracek (Czech word meaning "little cloud") is a small command line tool to manage your OpenStack [configuration files](https://docs.openstack.org/os-client-config/latest/user/configuration.html#config-files).

mracek is inspired by the awesome [kubectx/kubens](https://github.com/ahmetb/kubectx).

<img src="demo.gif" width="1300"  alt=""/>

## Features

* mracek supports auto-completion under bash, fish and zsh shells.
* mracek is opinionated about where it puts secrets (passwords, application credential secrets) and always puts them in secrets.yaml
* the directory where the openstack config files are stored is configurable (defaults to `$HOME/.config/openstack/`). See configuration section below.

## Examples
```shell
# Create a cloud
$ mracek create-cloud --username user1 --password very_secure --verify --auth-url https://cloud1.example.com:5000/v3 --project-name project1 --domain-name domain1 --region-name region1 cloud1

# Use a cloud (reads from openstack configuration files and sets OS_* environment variables)
$ mracek cloud1
Switching to cloud cloud1
$ env|grep OS_
OS_REGION_NAME=region1
OS_CLOUD=cloud1
OS_AUTH_URL=https://cloud1.example.com:5000/v3
OS_TENANT_NAME=project1
OS_USERNAME=user1
OS_DOMAIN_NAME=domain1
OS_PROJECT_NAME=project1
OS_PASSWORD=very_secure

# Create a profile (a profile is a cloud stored in clouds-public.yaml)
$ mracek create-profile --username user1 --password very_secure --verify --auth-url https://cloud1.example.com:5000/v3 --project-name project1 --domain-name domain1 --region-name region1 profile1

# Delete a cloud
$ mracek delete-cloud cloud1

# Delete a profile
$ mracek delete-profile profile1

# List profiles
$ mracek list-profiles
profile1

# Set cloud details
$ mracek set-cloud  --project-name project1 cloud1

# Set profile details
$ mracek set-profile  --project-name project1 profile1

# Show cloud details
$ mracek show-cloud cloud1
---
auth:
    auth_url: https://cloud1.example.com:5000/v3
    username: user1
    password: <masked>
    project_name: project1
    domain_name: domain1
region_name: region1
verify: true

# Show profile details
$ mracek show-profile profile1
---
auth:
    auth_url: https://cloud1.example.com:5000/v3
    username: user1
    password: <masked>
    project_name: project1
    domain_name: domain1
region_name: region1
verify: true
```

## Installation

```shell
go install github.com/mchlumsky/mracek@latest
```

## Configuration

mracek supports configuration through the configuration file `$HOME/.mracek.yaml` by default and can be changed with the `--config` command line flag.

Example:
```yaml
---
# Can be overridden by environment variable MRACEK_OS_CONFIG_DIR
os-config-dir: /path/to/openstack/config

# Can be overridden by environment variable MRACEK_SHELL
shell: /usr/bin/zsh
```
