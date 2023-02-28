# Marvin

Marvin is a powerful CLI tool designed to help Kubernetes cluster administrators 
ensure the security and reliability of their environments. 

Using a comprehensive set of [CEL (Common Expression Language)](https://github.com/google/cel-spec) expressions, 
Marvin performs extensive checks on cluster resources, 
identifying potential issues, misconfigurations, and vulnerabilities that could pose a risk to the system. 
It helps ensure that your Kubernetes clusters are always in compliance with best practices and industry standards.

<!-- TOC -->
* [Marvin](#marvin)
* [Installation](#installation)
  * [Manually](#manually)
  * [Install script](#install-script)
  * [From source](#from-source)
* [Usage](#usage)
* [Contributing](#contributing)
* [License](#license)
<!-- TOC -->

# Installation

The pre-compiled binaries are available in [GitHub releases page](https://github.com/undistro/marvin/releases) 
and can be installed manually or via script.

## Manually

1. Download the file for your system/architecture from the [GitHub releases page](https://github.com/undistro/marvin/releases)
2. Unpack the downloaded archive (e.g `tar -xzf marvin_0.1.0_Linux_x86_64.tar.gz`)
3. Make sure the binary has execution bit turned on (`chmod +x ./marvin`)
4. Move the binary somewhere in your `$PATH` (e.g `sudo mv ./marvin /usr/local/bin/`)

## Install script

The process above can be automated by the following script:

```shell
curl -sSfL https://raw.githubusercontent.com/undistro/marvin/main/install.sh | sh -s -- -b $HOME/.local/bin
```

## From source

```shell
go install github.com/undistro/marvin@latest
```

# Usage

Scan the current-context Kubernetes cluster performing the [builtin checks](internal/builtins):
```shell
marvin scan
```

# Contributing

We appreciate your contribution.
Please refer to our [contributing guideline](https://github.com/undistro/marvin/blob/main/CONTRIBUTING.md) for further information.
This project adheres to the Contributor Covenant [code of conduct](https://github.com/undistro/marvin/blob/main/CODE_OF_CONDUCT.md).

# License

Marvin is available under the Apache 2.0 license. See the [LICENSE](LICENSE) file for more info.
