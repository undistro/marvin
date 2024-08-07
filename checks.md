# Checks Overview  

In the table below, you can view all checks present on Marvin. Click on the #ID column item for more details about each check.


| Framework        | #ID                                                                                                  | Severity | Message                                                          |
|------------------|------------------------------------------------------------------------------------------------------|----------|------------------------------------------------------------------|
| CIS Benchmarks   | [M-500](/internal/builtins/cis/M-500_default_namespace.yaml)                                         | Medium   | Workloads in default namespace                                   |
| General          | [M-400](/internal/builtins/general/M-400_image_tag_latest.yaml)                                      | Medium   | Image tagged latest                                              |
|                  | [M-401](/internal/builtins/general/M-401_unmanaged_pod.yaml)                                         | Low      | Unmanaged Pod                                                    |
|                  | [M-402](/internal/builtins/general/M-402_readiness_probe.yaml)                                       | Medium   | Readiness and startup probe not configured                       |
|                  | [M-403](/internal/builtins/general/M-403_liveness_probe.yaml)                                        | Medium   | Liveness probe not configured                                    |
|                  | [M-404](/internal/builtins/general/M-404_memory_requests.yaml)                                       | Medium   | Memory requests not specified                                    |
|                  | [M-405](/internal/builtins/general/M-405_cpu_requests.yaml)                                          | Medium   | CPU requests not specified                                       |
|                  | [M-406](/internal/builtins/general/M-406_memory_limit.yaml)                                          | Medium   | Memory not limited                                               |
|                  | [M-407](/internal/builtins/general/M-407_cpu_limit.yaml)                                             | Medium   | CPU not limited                                                  |
|                  | [M-408](/internal/builtins/general/M-408_sudo_container_entrypoint.yaml)                             | Medium   | Sudo in container entrypoint                                     |
|                  | [M-409](/internal/builtins/general/M-409_deprecated_image_registry.yaml)                             | Medium   | Deprecated image registry                                        |
|                  | [M-410](/internal/builtins/general/M-410_resource_using_invalid_restartpolicy.yaml)                  | Medium   | Resource is using an invalid restartPolicy                       |
|                  | [M-411](/internal/builtins/general/M-411_role_binding_referencing_anonymous_or_unauthanticated.yaml) | Medium   | Role Binding referencing anonymous user or unauthenticated group |
| NSA-CISA         | [M-300](/internal/builtins/nsa/M-300_read_only_root_filesystem.yml)                                  | Low      | Root filesystem write allowed                                    |
| MITRE ATT&CK     | [M-200](/internal/builtins/mitre/M-200_allowed_registries.yml)                                       | Medium   | Image registry not allowed                                       |
|                  | [M-201](/internal/builtins/mitre/M-201_app_credentials.yml)                                          | High     | Application credentials stored in configuration files            |
|                  | [M-202](/internal/builtins/mitre/M-202_auto_mount_service_account.yml)                               | Low      | Automounted service account token                                |
|                  | [M-203](/internal/builtins/mitre/M-203_ssh.yml)                                                      | Low      | SSH server running inside container                              |
| PSS - Baseline   | [M-100](/internal/builtins/pss/baseline/M-100_host_process.yml)                                      | High     | Privileged access to the Windows node                            |
|                  | [M-101](/internal/builtins/pss/baseline/M-101_host_namespaces.yml)                                   | High     | Host namespaces                                                  |
|                  | [M-102](/internal/builtins/pss/baseline/M-102_privileged_containers.yml)                             | High     | Privileged container                                             |
|                  | [M-103](/internal/builtins/pss/baseline/M-103_capabilities.yml)                                      | High     | Insecure capabilities                                            |
|                  | [M-104](/internal/builtins/pss/baseline/M-104_host_path_volumes.yml)                                 | High     | HostPath volume                                                  |
|                  | [M-105](/internal/builtins/pss/baseline/M-105_host_ports.yml)                                        | High     | Not allowed hostPort                                             |
|                  | [M-106](/internal/builtins/pss/baseline/M-106_apparmor.yml)                                          | Medium   | Forbidden AppArmor profile                                       |
|                  | [M-107](/internal/builtins/pss/baseline/M-107_selinux.yml)                                           | Medium   | Forbidden SELinux options                                        |
|                  | [M-108](/internal/builtins/pss/baseline/M-108_proc_mount.yml)                                        | Medium   | Forbidden proc mount type                                        |
|                  | [M-109](/internal/builtins/pss/baseline/M-109_seccomp.yml)                                           | Medium   | Forbidden seccomp profile                                        |
|                  | [M-110](/internal/builtins/pss/baseline/M-110_sysctls.yml)                                           | Medium   | Unsafe sysctls                                                   |
| PSS - Restricted | [M-111](/internal/builtins/pss/restricted/M-111_volume_types.yml)                                    | Low      | Not allowed volume type                                          |
|                  | [M-112](/internal/builtins/pss/restricted/M-112_privilege_escalation.yml)                            | Medium   | Allowed privilege escalation                                     |
|                  | [M-113](/internal/builtins/pss/restricted/M-113_run_as_non_root.yml)                                 | Medium   | Container could be running as root user                          |
|                  | [M-114](/internal/builtins/pss/restricted/M-114_run_as_user.yml)                                     | Medium   | Container running as root UID                                    |
|                  | [M-115](/internal/builtins/pss/restricted/M-115_seccomp.yml)                                         | Low      | Not allowed seccomp profile                                      |
|                  | [M-116](/internal/builtins/pss/restricted/M-116_capabilities.yml)                                    | Low      | Not allowed added/dropped capabilities                           |
