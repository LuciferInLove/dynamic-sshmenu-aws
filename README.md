[![Go Report Card](https://goreportcard.com/badge/github.com/LuciferInLove/dynamic-sshmenu-aws)](https://goreportcard.com/report/github.com/LuciferInLove/dynamic-sshmenu-aws)
[![License](https://img.shields.io/badge/license-MIT-red.svg)](./LICENSE.md)
![Build status](https://github.com/LuciferInLove/dynamic-sshmenu-aws/workflows/Build/badge.svg)

# dynamic-sshmenu-aws

Dynamically creates a menu containing a list of AWS EC2 instances selected using tags.

## Overview

**dynamic-sshmenu-aws** generates sshmenu-style lists to connect to aws instances. It searches instances by aws instances tags that you can define as arguments. **dynamic-sshmenu-aws** executes `ssh __ip_address__` after choosing a menu item.

## Preparations for using

First of all, you should setup credentials to interact with AWS services:
* [via awscli](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-welcome.html)
* [manually](https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html#specifying-credentials)

If you are using bastion server, you can set it as proxy in ssh config as follows:

```
Host 172.31.*.*
  ProxyCommand ssh -W %h:%p 203.0.113.25
  ForwardAgent=yes
```

`172.31.*.*` - your aws instances private addresses range, `203.0.113.25` - bastion server public ip.

[Use ssh agent forwarding](https://developer.github.com/v3/guides/using-ssh-agent-forwarding/) to prevent keeping your private ssh keys on bastion servers.

AWS instances must have [tags](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/Using_Tags.html) to find by.

## Usage

You can see the **dynamic-sshmenu-aws** help by running it without arguments or with `-h` argument.

### Command Line Options

	--tags value,           -t value    instance tags in "key1:value1,value2;key2:value1" format. If undefined, full list will be shown
    --display-name value,   -d value    key of instance tag to display its values in results    (default: "Name")
    --public-ip,            -p          use public ip instead of private (default: false)
    --help,                 -h          show help
    --version,              -v          print the version

### Demo

![dynamic-sshmenu-aws](https://user-images.githubusercontent.com/34190954/87670302-2d67c600-c778-11ea-9bbd-89f72203c672.gif)

## Windows limitations

The application doesn't work in [mingw](http://www.mingw.org/) or similar terminals. You can use default cmd.exe, [windows terminal](https://github.com/microsoft/terminal) or run linux version of **dynamic-sshmenu-aws** in [wsl](https://docs.microsoft.com/en/windows/wsl/install-win10). Windows doesn't provide ssh connections ability by default. You must have `ssh.exe` installed in any of [PATH](https://docs.microsoft.com/en-us/windows/win32/shell/user-environment-variables) directories. For example, you can install [GitBash](https://gitforwindows.org/).
