<div align="center">

<picture>
  <source media="(prefers-color-scheme: dark)" srcset="assets/banner-dark.png">
  <img alt="Zora logo" src="assets/banner-light.png">
</picture>

[![Go Reference](https://pkg.go.dev/badge/github.com/undistro/marvin.svg)](https://pkg.go.dev/github.com/undistro/marvin)
[![Artifact Hub](https://img.shields.io/endpoint?url=https://artifacthub.io/badge/repository/marvin)](https://artifacthub.io/packages/krew/krew-index/marvin)
[![Test](https://github.com/undistro/marvin/actions/workflows/test.yml/badge.svg?branch=main&event=push)](https://github.com/undistro/marvin/actions/workflows/test.yml)
![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/undistro/marvin?sort=semver&color=brightgreen)
![GitHub](https://img.shields.io/github/license/undistro/marvin?color=brightgreen)
[![Go Report Card](https://goreportcard.com/badge/github.com/undistro/marvin)](https://goreportcard.com/report/github.com/undistro/marvin)
![GitHub all releases](https://img.shields.io/github/downloads/undistro/marvin/total)

</div>

Marvin is a CLI tool designed to help Kubernetes cluster administrators 
ensure the security and reliability of their environments. 

Using a comprehensive set of [CEL (Common Expression Language)](https://github.com/google/cel-spec) expressions, 
Marvin performs extensive checks on cluster resources, 
identifying potential issues, misconfigurations, and vulnerabilities that could pose a risk to the system. 
It helps ensure that your Kubernetes clusters are always in compliance with best practices and industry standards.

Marvin is also used as a plugin in [Zora](https://zora-docs.undistro.io/latest/).

<!-- TOC -->
* [Installation](#installation)
  * [Manually](#manually)
  * [Install via script](#install-via-script)
  * [Install via Krew](#install-via-krew)
  * [Install from source](#install-from-source)
* [Usage](#usage)
  * [Built-in checks](#built-in-checks)
  * [Custom checks](#custom-checks)
  * [Skipping resources](#skipping-resources)
  * [RBAC](#rbac)
* [Contributing](#contributing)
* [License](#license)
<!-- TOC -->

# Installation

The pre-compiled binaries are available in [GitHub releases page](https://github.com/undistro/marvin/releases) 
and can be installed manually, via script or as a `kubectl` plugin with [Krew](https://krew.sigs.k8s.io).

## Manually

1. Download the file for your system/architecture from the [GitHub releases page](https://github.com/undistro/marvin/releases)
2. Unpack the downloaded archive (e.g `tar -xzf marvin_Linux_x86_64.tar.gz`)
3. Make sure the binary has execution bit turned on (`chmod +x ./marvin`)
4. Move the binary somewhere in your `$PATH` (e.g `sudo mv ./marvin /usr/local/bin/`)

## Install via script

The process above can be automated by the following script:

```shell
curl -sSfL https://raw.githubusercontent.com/undistro/marvin/main/install.sh | sh -s -- -b $HOME/.local/bin
```

## Install via [Krew](https://krew.sigs.k8s.io)

You can install Marvin as a `kubectl` plugin via [Krew](https://krew.sigs.k8s.io):
```shell
kubectl krew install marvin
```

Then you can use Marvin with `kubectl` prefix:
```shell
kubectl marvin version
```

## Install from source

```shell
go install github.com/undistro/marvin@latest
```

# Usage

## Built-in checks

Scan the current-context Kubernetes cluster performing the [built-in checks](internal/builtins):
```shell
marvin scan
```
```
SEVERITY   ID      CHECK                                                   STATUS   FAILED   PASSED   SKIPPED 
High       M-101   Host namespaces                                         Failed   8        25       0         
High       M-104   HostPath volume                                         Failed   8        25       0         
High       M-201   Application credentials stored in configuration files   Failed   2        45       0         
High       M-102   Privileged container                                    Failed   2        31       0         
High       M-103   Insecure capabilities                                   Failed   2        31       0         
High       M-100   Privileged access to the Windows node                   Passed   0        33       0         
High       M-105   Not allowed hostPort                                    Passed   0        33       0         
Medium     M-113   Container could be running as root user                 Failed   33       0        0         
Medium     M-407   CPU not limited                                         Failed   31       2        0         
Medium     M-406   Memory not limited                                      Failed   27       6        0         
Medium     M-404   Memory requests not specified                           Failed   26       7        0         
Medium     M-402   Readiness and startup probe not configured              Failed   25       8        0         
Medium     M-403   Liveness probe not configured                           Failed   25       8        0         
Medium     M-405   CPU requests not specified                              Failed   23       10       0         
Medium     M-106   Forbidden AppArmor profile                              Passed   0        33       0         
Medium     M-107   Forbidden SELinux options                               Passed   0        33       0         
Medium     M-108   Forbidden proc mount type                               Passed   0        33       0         
Medium     M-109   Forbidden seccomp profile                               Passed   0        33       0         
Medium     M-110   Unsafe sysctls                                          Passed   0        33       0         
Medium     M-112   Allowed privilege escalation                            Passed   0        33       0         
Medium     M-200   Image registry not allowed                              Passed   0        33       0         
Medium     M-400   Image tagged latest                                     Passed   0        33       0         
Medium     M-408   Sudo in container entrypoint                            Passed   0        33       0         
Medium     M-409   Deprecated image registry                               Passed   0        33       0         
Medium     M-500   Workload in default namespace                           Passed   0        33       0         
Medium     M-410   Not allowed restartPolicy                               Passed   0        18       0         
Low        M-116   Not allowed added/dropped capabilities                  Failed   33       0        0         
Low        M-202   Automounted service account token                       Failed   33       0        0         
Low        M-115   Not allowed seccomp profile                             Failed   29       4        0         
Low        M-300   Root filesystem write allowed                           Failed   29       4        0         
Low        M-111   Not allowed volume type                                 Failed   8        25       0         
Low        M-203   SSH server running inside container                     Passed   0        39       0         
Low        M-114   Container running as root UID                           Passed   0        33       0         
Low        M-401   Unmanaged Pod                                           Passed   0        15       0         
```

The default output format is `table` which represents a summary of checks result. 
You can provide `json` or `yaml` in the `-o/--output` flag to get more details.

Run `marvin scan --help` to see all available options.

## Custom checks

Marvin allows you to write your own checks by using [CEL expressions](https://github.com/google/cel-spec) in a YAML file like the example below.

```yaml
id: CUSTOM-001
severity: Medium
message: "Replicas limit"
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

This is how built-in Marvin checks are defined as well.
You can see all the built-in checks in the [`internal/builtins` folder](internal/builtins) for examples.

If you want to quickly test CEL expressions from your browser, check out the [CEL Playground](https://playcel.undistro.io/).

Then provide the directory path with your custom check files in the `-f/--checks` flag:

```shell
marvin scan --disable-builtin --checks ./examples/
```
```
SEVERITY   ID           CHECK            STATUS   FAILED   PASSED   SKIPPED 
Medium     CUSTOM-001   Replicas limit   Passed   0        2        0         
```

The flag `--disable-builtin` disables the built-in Marvin checks.

If the check matches a PodSpec (`Pod`, `ReplicationController`, `ReplicaSet`, `Deployment`, `StatefulSet`, `DaemonSet`, `Job` or `CronJob`)
the `podSpec` and `allContainers` inputs are available for expressions.

The `allContainers` input is a list of all containers including `initContainers` and `ephemeralContainers`.

## Skipping resources

You can use annotations to skip certain checks for specific resources in your cluster.
By adding the `marvin.undistro.io/skip` annotation to a resource, 
you can specify a comma-separated list of check IDs to skip.

Example:
```shell
kubectl annotate deployment nginx marvin.undistro.io/skip='M-202, M-111'
```

By default, Marvin will respect the `marvin.undistro.io/skip` annotation when performing checks. 
However, you can disable this behavior by using the `--disable-annotation-skip` flag.
This flag will cause Marvin to perform all checks on all resources.

If you prefer to use a different annotation to skip checks, 
you can use the `--skip-annotation` flag to specify the annotation name. 
Example: `--skip-annotation='my-company.com/skip-checks'`

## RBAC

Currently, the built-in checks look for the below resources
and Marvin needs view (`get` and `list`) permission to verify them.

- `v1/pods`
- `v1/configmaps`
- `v1/services`
- `apps/v1/deployments`
- `apps/v1/daemonsets`
- `apps/v1/statefulsets`
- `apps/v1/replicasets`
- `batch/v1/cronjobs`
- `batch/v1/jobs`

<details>

<summary> Here is a sample `ClusterRole` for Marvin: </summary>

```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: marvin
rules:
  - apiGroups: [ "" ]
    resources:
      - configmaps
      - pods
      - services
    verbs: [ "get", "list" ]
  - apiGroups: [ "apps" ]
    resources:
      - daemonsets
      - deployments
      - statefulsets
      - replicasets
    verbs: [ "get", "list" ]
  - apiGroups: [ batch ]
    resources:
      - jobs
      - cronjobs
    verbs: [ "get", "list" ]
```

</details>

> **Note**
> You can write a custom check to look at any resource. 
> But Marvin needs view permission. 
> Remember to update RBAC for new resources you want to check.

# Contributing

We appreciate your contribution.
Please refer to our [contributing guideline](https://github.com/undistro/marvin/blob/main/CONTRIBUTING.md) for further information.
This project adheres to the Contributor Covenant [code of conduct](https://github.com/undistro/marvin/blob/main/CODE_OF_CONDUCT.md).

# License

Marvin is available under the Apache 2.0 license. See the [LICENSE](LICENSE) file for more info.
