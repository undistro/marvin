# Checks Overview  

The table below, we can see all checks present on Marvin.   

| Framework | # ID        | Severity | Message              |
|-----------|:-------------:|:----------:|----------------------|
| General | [M-400](/internal/builtins/general/M-400_image_tag_latest.yaml) |  Medium  |Image tagged latest| 
|         | [M-401](/internal/builtins/general/M-401_unmanaged_pod.yaml)| Low | Unmanaged Pod|      
|         | [M-402](/internal/builtins/general/M-402_readiness_probe.yaml)|Medium |Readiness and startup probe not configured| 
|         | [M-403](/internal/builtins/general/M-403_liveness_probe.yaml)|  Medium  | Liveness probe not configured|
|         | [M-404](/internal/builtins/general/M-404_memory_requests.yaml)| Medium  | Memory requests not specified|
|         | [M-405](/internal/builtins/general/M-405_cpu_requests.yaml)  |  Medium  | CPU requests not specified|
|         | [M-406](/internal/builtins/general/M-406_memory_limit.yaml)  |  Medium  | Memory not limited|
|         | [M-407](/internal/builtins/general/M-407_cpu_limit.yaml)     |  Medium  | CPU not limited| 
| NSA     | [M-300](/internal/builtins/nsa/M-300_read_only_root_filesystem.yml)|  Low|   Root filesystem write allowed|   
| Mitre   | [M-200](/internal/builtins/mitre/M-200_allowed_registries.yml)|Medium|Image registry not allowed| 
|         |[M-201](/internal/builtins/mitre/M-201_app_credentials.yml)|High|Application credentials stored in configuration files|
|         |[M-202](/internal/builtins/mitre/M-202_auto_mount_service_account.yml)|Low|Automounted service account token|
|         |[M-203](/internal/builtins/mitre/M-203_ssh.yml)| Low |SSH server running inside container|
|         |                                                                    |     |        |
| PSS - Basline|[M100](/internal/builtins/pss/baseline/M-100_host_process.yml)|High|Privileged access to the Windows node|
|         |[M101](/internal/builtins/pss/baseline/M-101_host_namespaces.yml)|High|Host namespaces|
|         |[M102](/internal/builtins/pss/baseline/M-102_privileged_containers.yml)|High|Privileged container|
|         |[M103](/internal/builtins/pss/baseline/M-103_capabilities.yml)|High|Insecure capabilities|
|         |[M104](/internal/builtins/pss/baseline/M-104_host_path_volumes.yml)|High|HostPath volume|
|         |[M105](/internal/builtins/pss/baseline/M-105_host_ports.yml)|High|Not allowed hostPort|
|         |[M106](/internal/builtins/pss/baseline/M-106_apparmor.yml)|Medium|Forbidden AppArmor profile|
|         |[M107](/internal/builtins/pss/baseline/M-107_selinux.yml)|Medium|Forbidden SELinux options|
|         |[M108](/internal/builtins/pss/baseline/M-108_proc_mount.yml)|Medium|Forbidden proc mount type|
|         |[M109](/internal/builtins/pss/baseline/M-109_seccomp.yml)|Medium|Forbidden seccomp profile|
|         |[M110](/internal/builtins/pss/baseline/M-110_sysctls.yml)|Medium|Unsafe sysctls|
| CIS     |           |          |                      | | 
