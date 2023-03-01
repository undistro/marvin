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
  * [Custom checks](#custom-checks)
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

Scan the current-context Kubernetes cluster performing the [built-in checks](internal/builtins):
```shell
marvin scan
```
```
SEVERITY   ID                      CHECK                                                   STATUS   FAILED   PASSED   SKIPPED 
High       capabilities-added      Insecure capabilities                                   Failed   1        3        0         
High       app-credentials         Application credentials stored in configuration files   Failed   1        4        0         
High       privileged-containers   Privileged container                                    Failed   1        3        0         
High       host-namespaces         Host namespaces                                         Failed   1        2        0         
High       host-path-volumes       HostPath volume                                         Failed   1        2        0         
High       host-ports              Not allowed hostPort                                    Passed   0        3        0         
High       host-process            Privileged access to the Windows node                   Passed   0        3        0         
Medium     run-as-non-root         Container could be running as root user                 Failed   3        0        0         
Medium     seccomp-baseline        Forbidden seccomp profile                               Passed   0        3        0         
Medium     image-registry          Image registry not allowed                              Passed   0        3        0         
Medium     privilege-escalation    Allowed privilege escalation                            Passed   0        3        0         
Medium     sysctls                 Unsafe sysctls                                          Passed   0        3        0         
Medium     selinux                 Forbidden SELinux options                               Passed   0        3        0         
Medium     proc-mount              Forbidden proc mount type                               Passed   0        3        0         
Medium     apparmor                Forbidden AppArmor profile                              Passed   0        3        0         
Low        seccomp-restricted      Not allowed seccomp profile                             Failed   3        0        0         
Low        read-only-root-fs       Root filesystem write allowed                           Failed   3        2        0         
Low        auto-mount-sa-token     Automounted service account token                       Failed   3        0        0         
Low        capabilities            Not allowed added/dropped capabilities                  Failed   3        0        0         
Low        volume-types            Not allowed volume type                                 Failed   1        2        0         
Low        ssh-server              SSH server running inside container                     Passed   0        4        0         
Low        run-as-user             Container running as root UID                           Passed   0        3        0         
```

The default output format is `table` which represents a summary of checks result. 
You can provide `json` or `yaml` in the `-o/--output` flag to get more details.

Run `marvin scan --help` to see all available options.

## Custom checks

Marvin allows you to write your own checks by using [CEL expressions](https://github.com/google/cel-spec) in a YAML file like the example below.

```yaml
id: example-replicas
severity: Medium
message: "Replica limit"
match:
  resources:
    - group: apps
      version: v1
      resource: deployments
validations:
  - expression: >
      object.spec.replicas <= 5
    message: "Deployment with more than 5 replicas"
```

If an expression evaluates to `false`, the check fails.

Then you can provide the directory path with your custom check files in the `-f/--checks` flag:

```shell
marvin scan --disable-builtin --checks ./examples/
```
```
SEVERITY   ID                 CHECK           STATUS   FAILED   PASSED   SKIPPED 
Medium     example-replicas   Replica limit   Passed   0        1        0            
```

The flag `--disable-builtin` disables the built-in Marvin checks.

This is how built-in Marvin checks are defined as well. 
You can see all the built-in checks in the [`internal/builtins` folder](internal/builtins) for examples.

If the check matches a PodSpec (`Pod`, `ReplicationController`, `ReplicaSet`, `Deployment`, `StatefulSet`, `DaemonSet`, `Job` or `CronJob`)
the `podSpec` and `allContainers` inputs are available for expressions.

The `allContainers` input is a list of all containers including `initContainers` and `ephemeralContainers`.

# Contributing

We appreciate your contribution.
Please refer to our [contributing guideline](https://github.com/undistro/marvin/blob/main/CONTRIBUTING.md) for further information.
This project adheres to the Contributor Covenant [code of conduct](https://github.com/undistro/marvin/blob/main/CODE_OF_CONDUCT.md).

# License

Marvin is available under the Apache 2.0 license. See the [LICENSE](LICENSE) file for more info.
