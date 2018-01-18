# undefined

#### Type: object

__add description to json file__

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| lifecycle | array | `lifecycle` can be used to customize the end-customer's experience when generating/uploading support bundles. |  |
| specs | array | A list of items to include in the support bundle. |  |

*****

## lifecycle:undefined

#### Type: object

__`lifecycle` can be used to customize the end-customer's experience when generating/uploading support bundles.__

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| [generate](#generate) | object | Adds a step that will generate a support bundle. It has no arguments, but should be specified with an empty object `{}` |  |
| [message](#message) | object | `message` can be used to print a message to the end customer. |  |
| [upload](#upload) | object | Adds a step in which the end customer can upload a bundle to https://console.replicated.com for review by the vendor. |  |

*****

## generate

#### Type: object

__Adds a step that will generate a support bundle. It has no arguments, but should be specified with an empty object `{}`__


*****

## message

#### Type: object

__`message` can be used to print a message to the end customer.__

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| contents | string | `contents` is the message |  |

*****

## upload

#### Type: object

__Adds a step in which the end customer can upload a bundle to https://console.replicated.com for review by the vendor.__

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| [prompt](#prompt) | object | Allows configuring an interactive prompt for the user to confirm before uploading |  |

*****

## prompt

#### Type: object

__Allows configuring an interactive prompt for the user to confirm before uploading__

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| accept | string | A message to display if the user accepts the upload prompt |  |
| decline | string | A message to display if the user declines the upload prompt |  |
| [default](#default) | boolean | The default answer to the prompt |  |
| [message](#message) | string | Message to display before prompting |  |

*****

## specs:undefined

#### Type: object

__A list of items to include in the support bundle.__

Required: output_dir

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| description | string | A description of the file(s) that will be collected by this spec |  |
| [docker.container-cp](#docker.container-cp) | object | Collect a file by copying from a running docker container | `docker.container-cp:  { container: my-sweet-app }` |
| [docker.container-exec](#docker.container-exec) | object | Collect the stdout and stderr of `exec`-ing a command on a running docker container |  |
| [docker.container-inspect](#docker.container-inspect) | object | Collect the `docker inspect` output for one or more running or stopped containers |  |
| [docker.container-logs](#docker.container-logs) | object | Collect the logs from one or more docker containers |  |
| [docker.container-ls](#docker.container-ls) | object | Collect information about one or more containers |  |
| [docker.container-run](#docker.container-run) | object | Collect the stderr/stdout of running a single docker container |  |
| [docker.exec](#docker.exec) | object | Collect the stdout/stderr of executing a command in an already running docker container |  |
| [docker.image-ls](#docker.image-ls) | object | Collect a list of docker images present on the server |  |
| [docker.images](#docker.images) | object | Collect a list of docker images present on the server |  |
| [docker.info](#docker.info) | object | Collect info about the Docker daemon |  |
| [docker.logs](#docker.logs) | object | Collect info about the Docker daemon |  |
| [docker.node-ls](#docker.node-ls) | object | Collect information about the nodes in a Docker Swarm installation |  |
| [docker.ps](#docker.ps) | object | Collect information about containers |  |
| [docker.run](#docker.run) | object | Collect the stderr/stdout of running a single docker container |  |
| [docker.service-logs](#docker.service-logs) | object | Collect logs from a docker swarm service |  |
| [docker.service-ls](#docker.service-ls) | object | Collect a list of docker swarm services |  |
| [docker.service-ps](#docker.service-ps) | object | Collect information about the tasks run by one or more services |  |
| [docker.stack-service-logs](#docker.stack-service-logs) | object | Collect logs from one or more services in a stack |  |
| [docker.stack-service-ls](#docker.stack-service-ls) | object |  |  |
| [docker.stack-service-ps](#docker.stack-service-ps) | object |  |  |
| [docker.stack-task-logs](#docker.stack-task-logs) | object |  |  |
| [docker.task-logs](#docker.task-logs) | object |  |  |
| [docker.task-ls](#docker.task-ls) | object |  |  |
| [docker.version](#docker.version) | object |  |  |
| [journald.logs](#journald.logs) | object |  |  |
| [kubernetes.api-versions](#kubernetes.api-versions) | object |  |  |
| [kubernetes.cluster-info](#kubernetes.cluster-info) | object |  |  |
| [kubernetes.logs](#kubernetes.logs) | object |  |  |
| [kubernetes.resource-list](#kubernetes.resource-list) | object |  |  |
| [kubernetes.version](#kubernetes.version) | object |  |  |
| [meta](#meta) | object |  |  |
| [meta.customer](#meta.customer) | object |  |  |
| [os.hostname](#os.hostname) | object |  |  |
| [os.http-request](#os.http-request) | object |  |  |
| [os.loadavg](#os.loadavg) | object |  |  |
| [os.read-file](#os.read-file) | object |  |  |
| [os.run-command](#os.run-command) | object |  |  |
| [os.uptime](#os.uptime) | object |  |  |
| output_dir | string |  |  |
| [retraced.events](#retraced.events) | object |  |  |
| [scrub](#scrub) | object |  |  |
| timeout_seconds | integer |  |  |
| [version](#version) | object |  |  |

*****

## docker.container-cp

#### Type: object

__Collect a file by copying from a running docker container__

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| container | string | The name of the container to copy from |  |
| src_path | string | The path of the target file in the container's filesystem |  |

*****

## docker.container-exec

#### Type: object

__Collect the stdout and stderr of `exec`-ing a command on a running docker container__

Required: container, exec_config

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| container | string | The name of the container to run the command in |  |
| [exec_config](#exec_config) | object | Same as would be passed to `docker exec` |  |

*****

## exec_config

#### Type: object

__Same as would be passed to `docker exec`Config options as would be passed to `docker exec`__

Required: User, Privileged, Tty, AttachStdin, AttachStderr, AttachStdout, Detach, DetachKeys, Env, Cmd, User, Privileged, Tty, AttachStdin, AttachStderr, AttachStdout, Detach, DetachKeys, Env, Cmd

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| AttachStderr | boolean | Same as would be passed to `docker exec` |  |
| AttachStdin | boolean | Same as would be passed to `docker exec` |  |
| AttachStdout | boolean | Same as would be passed to `docker exec` |  |
| Cmd | array | Same as would be passed to `docker exec` |  |
| Detach | boolean | Same as would be passed to `docker exec` |  |
| DetachKeys | string | Same as would be passed to `docker exec` |  |
| Env | array | Same as would be passed to `docker exec` |  |
| Privileged | boolean | Same as would be passed to `docker exec` |  |
| Tty | boolean | Same as would be passed to `docker exec` |  |
| User | string | Same as would be passed to `docker exec` |  |

*****

## Cmd:undefined

#### Type: string

__Same as would be passed to `docker exec`Container create options as would be passed to `docker run`Same as would be passed to `docker exec`Container create options as would be passed to `docker run`__


*****

## Env:undefined

#### Type: string

__Same as would be passed to `docker exec`Container create options as would be passed to `docker run`Same as would be passed to `docker exec`Container create options as would be passed to `docker run`__


*****

## docker.container-inspect

#### Type: object

__Collect the `docker inspect` output for one or more running or stopped containers__

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| container | string | the container name |  |
| [container_list_options](#container_list_options) | object | Options to be used to filter the list of containers |  |

*****

## container_list_options

#### Type: object

__Options to be used to filter the list of containers__

Required: Quiet, Size, All, Latest, Since, Before, Limit, Filters, Quiet, Size, All, Latest, Since, Before, Limit, Filters, Quiet, Size, All, Latest, Since, Before, Limit, Filters

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| All | boolean | Same as `-a` option to `docker ps`  |  |
| Before | string | Same as `--before` option to `docker logs`  |  |
| [Filters](#filters) | object | Same as `--filter` option to `docker ps`  |  |
| Latest | boolean | Same as would be passed to `docker ps` |  |
| Limit | integer | Same as would be passed to `docker ps` |  |
| Quiet | boolean | Same as would be passed to `docker ps` |  |
| Since | string | Same as would be passed to `docker ps` |  |
| Size | boolean | Same as would be passed to `docker ps` |  |

*****

## Filters

#### Type: object

__Same as `--filter` option to `docker ps` Same as would be passed to `docker images`Same as would be passed to `docker images`Same as `--filter` option to `docker ps` Same as would be passed to `docker node ls`Same as `--filter` option to `docker ps` Same as would be passed to `docker service ls`Same as would be passed to `docker service ls`Same as would be passed to `docker service ps`undefinedundefinedundefinedundefinedundefinedundefined__

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| [.*](#.*) | array |  |  |

*****

## docker.container-logs

#### Type: object

__Collect the logs from one or more docker containers__

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| container | string | A container name |  |
| [container_list_options](#container_list_options) | object | Options to be used to filter the list of containers |  |
| [container_logs_options](#container_logs_options) | object | Args as would be given to `docker logs` |  |

*****

## container_logs_options

#### Type: object

__Args as would be given to `docker logs`undefinedundefined__

Required: ShowStdout, ShowStderr, Since, Until, Timestamps, Follow, Tail, Details, ShowStdout, ShowStderr, Since, Until, Timestamps, Follow, Tail, Details, ShowStdout, ShowStderr, Since, Until, Timestamps, Follow, Tail, Details, ShowStdout, ShowStderr, Since, Until, Timestamps, Follow, Tail, Details, ShowStdout, ShowStderr, Since, Until, Timestamps, Follow, Tail, Details, ShowStdout, ShowStderr, Since, Until, Timestamps, Follow, Tail, Details

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| Details | boolean |  |  |
| Follow | boolean |  |  |
| ShowStderr | boolean |  |  |
| ShowStdout | boolean |  |  |
| Since | string |  |  |
| Tail | string |  |  |
| Timestamps | boolean |  |  |
| Until | string |  |  |

*****

## docker.container-ls

#### Type: object

__Collect information about one or more containers__

Required: Quiet, Size, All, Latest, Since, Before, Limit, Filters

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| All | boolean | Same as `-a` option to `docker ps`  |  |
| Before | string | Same as `--before` option to `docker logs`  |  |
| [Filters](#filters) | object | Same as `--filter` option to `docker ps`  |  |
| Latest | boolean | Same as would be passed to `docker ps` |  |
| Limit | integer | Same as would be passed to `docker ps` |  |
| Quiet | boolean | Same as would be passed to `docker ps` |  |
| Since | string | Same as would be passed to `docker ps` |  |
| Size | boolean | Same as would be passed to `docker ps` |  |

*****

## docker.container-run

#### Type: object

__Collect the stderr/stdout of running a single docker container__

Required: container_create_config

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| [container_create_config](#container_create_config) | object | Container create options as would be passed to `docker run` |  |
| enable_pull | boolean | If `true`, allow this container to be pulled if not present |  |

*****

## container_create_config

#### Type: object

__Container create options as would be passed to `docker run`__

Required: Name, Config, HostConfig, NetworkingConfig, AdjustCPUShares, Name, Config, HostConfig, NetworkingConfig, AdjustCPUShares

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| AdjustCPUShares | boolean | Container create options as would be passed to `docker run` |  |
| [Config](#config) | object | Container create options as would be passed to `docker run` |  |
| [HostConfig](#hostconfig) | object | Container create options as would be passed to `docker run` |  |
| Name | string | Container create options as would be passed to `docker run` |  |
| [NetworkingConfig](#networkingconfig) | object | Container create options as would be passed to `docker run` |  |

*****

## Config

#### Type: object

__Container create options as would be passed to `docker run`__

Required: Hostname, Domainname, User, AttachStdin, AttachStdout, AttachStderr, Tty, OpenStdin, StdinOnce, Env, Cmd, Image, Volumes, WorkingDir, Entrypoint, OnBuild, Labels, Hostname, Domainname, User, AttachStdin, AttachStdout, AttachStderr, Tty, OpenStdin, StdinOnce, Env, Cmd, Image, Volumes, WorkingDir, Entrypoint, OnBuild, Labels

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| ArgsEscaped | boolean | Container create options as would be passed to `docker run` |  |
| AttachStderr | boolean | Container create options as would be passed to `docker run` |  |
| AttachStdin | boolean | Container create options as would be passed to `docker run` |  |
| AttachStdout | boolean | Container create options as would be passed to `docker run` |  |
| Cmd | array | Container create options as would be passed to `docker run` |  |
| Domainname | string | Container create options as would be passed to `docker run` |  |
| Entrypoint | array | Container create options as would be passed to `docker run` |  |
| Env | array | Container create options as would be passed to `docker run` |  |
| [ExposedPorts](#exposedports) | object | Container create options as would be passed to `docker run` |  |
| [Healthcheck](#healthcheck) | object | Container create options as would be passed to `docker run` |  |
| Hostname | string | Container create options as would be passed to `docker run` |  |
| Image | string | Container create options as would be passed to `docker run` |  |
| [Labels](#labels) | object | Container create options as would be passed to `docker run` |  |
| MacAddress | string | Container create options as would be passed to `docker run` |  |
| NetworkDisabled | boolean | Container create options as would be passed to `docker run` |  |
| OnBuild | array | Container create options as would be passed to `docker run` |  |
| OpenStdin | boolean | Container create options as would be passed to `docker run` |  |
| Shell | array | Container create options as would be passed to `docker run` |  |
| StdinOnce | boolean | Container create options as would be passed to `docker run` |  |
| StopSignal | string | Container create options as would be passed to `docker run` |  |
| StopTimeout | integer | Container create options as would be passed to `docker run` |  |
| Tty | boolean | Container create options as would be passed to `docker run` |  |
| User | string | Container create options as would be passed to `docker run` |  |
| [Volumes](#volumes) | object | Container create options as would be passed to `docker run` |  |
| WorkingDir | string | Container create options as would be passed to `docker run` |  |
| [.*](#.*) | string | Container create options as would be passed to `docker run` |  |

*****

## Entrypoint:undefined

#### Type: string

__Container create options as would be passed to `docker run`__


*****

## ExposedPorts

#### Type: object

__Container create options as would be passed to `docker run`__

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| [.*](#.*) | object | Container create options as would be passed to `docker run` |  |

*****

## .*

#### Type: object

__Container create options as would be passed to `docker run`__


*****

## Healthcheck

#### Type: object

__Container create options as would be passed to `docker run`__

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| Interval | integer | Container create options as would be passed to `docker run` |  |
| Retries | integer | Container create options as would be passed to `docker run` |  |
| StartPeriod | integer | Container create options as would be passed to `docker run` |  |
| Test | array | Container create options as would be passed to `docker run` |  |
| Timeout | integer | Container create options as would be passed to `docker run` |  |

*****

## Test:undefined

#### Type: string

__Container create options as would be passed to `docker run`__


*****

## Labels

#### Type: object

__Container create options as would be passed to `docker run`__

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| [.*](#.*) | string | Container create options as would be passed to `docker run` |  |

*****

## OnBuild:undefined

#### Type: string

__Container create options as would be passed to `docker run`__


*****

## Shell:undefined

#### Type: string

__Container create options as would be passed to `docker run`__


*****

## Volumes

#### Type: object

__Container create options as would be passed to `docker run`__

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| [.*](#.*) | object | Container create options as would be passed to `docker run` |  |

*****

## HostConfig

#### Type: object

__Container create options as would be passed to `docker run`__

Required: Binds, ContainerIDFile, LogConfig, NetworkMode, PortBindings, RestartPolicy, AutoRemove, VolumeDriver, VolumesFrom, CapAdd, CapDrop, Dns, DnsOptions, DnsSearch, ExtraHosts, GroupAdd, IpcMode, Cgroup, Links, OomScoreAdj, PidMode, Privileged, PublishAllPorts, ReadonlyRootfs, SecurityOpt, UTSMode, UsernsMode, ShmSize, ConsoleSize, Isolation, CpuShares, Memory, NanoCpus, CgroupParent, BlkioWeight, BlkioWeightDevice, BlkioDeviceReadBps, BlkioDeviceWriteBps, BlkioDeviceReadIOps, BlkioDeviceWriteIOps, CpuPeriod, CpuQuota, CpuRealtimePeriod, CpuRealtimeRuntime, CpusetCpus, CpusetMems, Devices, DeviceCgroupRules, DiskQuota, KernelMemory, MemoryReservation, MemorySwap, MemorySwappiness, OomKillDisable, PidsLimit, Ulimits, CpuCount, CpuPercent, IOMaximumIOps, IOMaximumBandwidth, Binds, ContainerIDFile, LogConfig, NetworkMode, PortBindings, RestartPolicy, AutoRemove, VolumeDriver, VolumesFrom, CapAdd, CapDrop, Dns, DnsOptions, DnsSearch, ExtraHosts, GroupAdd, IpcMode, Cgroup, Links, OomScoreAdj, PidMode, Privileged, PublishAllPorts, ReadonlyRootfs, SecurityOpt, UTSMode, UsernsMode, ShmSize, ConsoleSize, Isolation, CpuShares, Memory, NanoCpus, CgroupParent, BlkioWeight, BlkioWeightDevice, BlkioDeviceReadBps, BlkioDeviceWriteBps, BlkioDeviceReadIOps, BlkioDeviceWriteIOps, CpuPeriod, CpuQuota, CpuRealtimePeriod, CpuRealtimeRuntime, CpusetCpus, CpusetMems, Devices, DeviceCgroupRules, DiskQuota, KernelMemory, MemoryReservation, MemorySwap, MemorySwappiness, OomKillDisable, PidsLimit, Ulimits, CpuCount, CpuPercent, IOMaximumIOps, IOMaximumBandwidth

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| AutoRemove | boolean | Container create options as would be passed to `docker run` |  |
| Binds | array | Container create options as would be passed to `docker run` |  |
| BlkioDeviceReadBps | array | Container create options as would be passed to `docker run` |  |
| BlkioDeviceReadIOps | array | Container create options as would be passed to `docker run` |  |
| BlkioDeviceWriteBps | array | Container create options as would be passed to `docker run` |  |
| BlkioDeviceWriteIOps | array | Container create options as would be passed to `docker run` |  |
| BlkioWeight | integer | Container create options as would be passed to `docker run` |  |
| BlkioWeightDevice | array | Container create options as would be passed to `docker run` |  |
| CapAdd | array | Container create options as would be passed to `docker run` |  |
| CapDrop | array | Container create options as would be passed to `docker run` |  |
| Cgroup | string | Container create options as would be passed to `docker run` |  |
| CgroupParent | string | Container create options as would be passed to `docker run` |  |
| ConsoleSize | undefined | Container create options as would be passed to `docker run` |  |
| ContainerIDFile | string | Container create options as would be passed to `docker run` |  |
| CpuCount | integer | Container create options as would be passed to `docker run` |  |
| CpuPercent | integer | Container create options as would be passed to `docker run` |  |
| CpuPeriod | integer | Container create options as would be passed to `docker run` |  |
| CpuQuota | integer | Container create options as would be passed to `docker run` |  |
| CpuRealtimePeriod | integer | Container create options as would be passed to `docker run` |  |
| CpuRealtimeRuntime | integer | Container create options as would be passed to `docker run` |  |
| CpuShares | integer | Container create options as would be passed to `docker run` |  |
| CpusetCpus | string | Container create options as would be passed to `docker run` |  |
| CpusetMems | string | Container create options as would be passed to `docker run` |  |
| DeviceCgroupRules | array | Container create options as would be passed to `docker run` |  |
| Devices | array | Container create options as would be passed to `docker run` |  |
| DiskQuota | integer | Container create options as would be passed to `docker run` |  |
| Dns | array | Container create options as would be passed to `docker run` |  |
| DnsOptions | array | Container create options as would be passed to `docker run` |  |
| DnsSearch | array | Container create options as would be passed to `docker run` |  |
| ExtraHosts | array | Container create options as would be passed to `docker run` |  |
| GroupAdd | array | Container create options as would be passed to `docker run` |  |
| IOMaximumBandwidth | integer | Container create options as would be passed to `docker run` |  |
| IOMaximumIOps | integer | Container create options as would be passed to `docker run` |  |
| Init | boolean | Container create options as would be passed to `docker run` |  |
| IpcMode | string | Container create options as would be passed to `docker run` |  |
| Isolation | string | Container create options as would be passed to `docker run` |  |
| KernelMemory | integer | Container create options as would be passed to `docker run` |  |
| Links | array | Container create options as would be passed to `docker run` |  |
| [LogConfig](#logconfig) | object | Container create options as would be passed to `docker run` |  |
| Memory | integer | Container create options as would be passed to `docker run` |  |
| MemoryReservation | integer | Container create options as would be passed to `docker run` |  |
| MemorySwap | integer | Container create options as would be passed to `docker run` |  |
| MemorySwappiness | integer | Container create options as would be passed to `docker run` |  |
| Mounts | array | Container create options as would be passed to `docker run` |  |
| NanoCpus | integer | Container create options as would be passed to `docker run` |  |
| NetworkMode | string | Container create options as would be passed to `docker run` |  |
| OomKillDisable | boolean | Container create options as would be passed to `docker run` |  |
| OomScoreAdj | integer | Container create options as would be passed to `docker run` |  |
| PidMode | string | Container create options as would be passed to `docker run` |  |
| PidsLimit | integer | Container create options as would be passed to `docker run` |  |
| [PortBindings](#portbindings) | object | Container create options as would be passed to `docker run` |  |
| Privileged | boolean | Container create options as would be passed to `docker run` |  |
| PublishAllPorts | boolean | Container create options as would be passed to `docker run` |  |
| ReadonlyRootfs | boolean | Container create options as would be passed to `docker run` |  |
| [RestartPolicy](#restartpolicy) | object | Container create options as would be passed to `docker run` |  |
| Runtime | string | Container create options as would be passed to `docker run` |  |
| SecurityOpt | array | Container create options as would be passed to `docker run` |  |
| ShmSize | integer | Container create options as would be passed to `docker run` |  |
| [StorageOpt](#storageopt) | object | Container create options as would be passed to `docker run` |  |
| [Sysctls](#sysctls) | object | Container create options as would be passed to `docker run` |  |
| [Tmpfs](#tmpfs) | object | Container create options as would be passed to `docker run` |  |
| UTSMode | string | Container create options as would be passed to `docker run` |  |
| Ulimits | array | Container create options as would be passed to `docker run` |  |
| UsernsMode | string | Container create options as would be passed to `docker run` |  |
| VolumeDriver | string | Container create options as would be passed to `docker run` |  |
| VolumesFrom | array | Container create options as would be passed to `docker run` |  |

*****

## Binds:undefined

#### Type: string

__Container create options as would be passed to `docker run`__


*****

## CapAdd:undefined

#### Type: string

__Container create options as would be passed to `docker run`__


*****

## CapDrop:undefined

#### Type: string

__Container create options as would be passed to `docker run`__


*****

## DeviceCgroupRules:undefined

#### Type: string

__Container create options as would be passed to `docker run`__


*****

## Devices:undefined

#### Type: object

__Container create options as would be passed to `docker run`__

Required: PathOnHost, PathInContainer, CgroupPermissions, PathOnHost, PathInContainer, CgroupPermissions

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| CgroupPermissions | string | Container create options as would be passed to `docker run` |  |
| PathInContainer | string | Container create options as would be passed to `docker run` |  |
| PathOnHost | string | Container create options as would be passed to `docker run` |  |

*****

## Dns:undefined

#### Type: string

__Container create options as would be passed to `docker run`__


*****

## DnsOptions:undefined

#### Type: string

__Container create options as would be passed to `docker run`__


*****

## DnsSearch:undefined

#### Type: string

__Container create options as would be passed to `docker run`__


*****

## ExtraHosts:undefined

#### Type: string

__Container create options as would be passed to `docker run`__


*****

## GroupAdd:undefined

#### Type: string

__Container create options as would be passed to `docker run`__


*****

## Links:undefined

#### Type: string

__Container create options as would be passed to `docker run`__


*****

## LogConfig

#### Type: object

__Container create options as would be passed to `docker run`__

Required: Type, Config, Type, Config

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| [Config](#config) | object | Container create options as would be passed to `docker run` |  |
| Type | string | Container create options as would be passed to `docker run` |  |

*****

## Mounts:undefined

#### Type: object

__Container create options as would be passed to `docker run`__

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| [BindOptions](#bindoptions) | object | Container create options as would be passed to `docker run` |  |
| Consistency | string | Container create options as would be passed to `docker run` |  |
| ReadOnly | boolean | Container create options as would be passed to `docker run` |  |
| Source | string | Container create options as would be passed to `docker run` |  |
| Target | string | Container create options as would be passed to `docker run` |  |
| [TmpfsOptions](#tmpfsoptions) | object | Container create options as would be passed to `docker run` |  |
| Type | string | Container create options as would be passed to `docker run` |  |
| [VolumeOptions](#volumeoptions) | object | Container create options as would be passed to `docker run` |  |

*****

## BindOptions

#### Type: object

__Container create options as would be passed to `docker run`__

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| Propagation | string | Container create options as would be passed to `docker run` |  |

*****

## TmpfsOptions

#### Type: object

__Container create options as would be passed to `docker run`__

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| Mode | integer | Container create options as would be passed to `docker run` |  |
| SizeBytes | integer | Container create options as would be passed to `docker run` |  |

*****

## VolumeOptions

#### Type: object

__Container create options as would be passed to `docker run`__

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| [DriverConfig](#driverconfig) | object | Container create options as would be passed to `docker run` |  |
| [Labels](#labels) | object | Container create options as would be passed to `docker run` |  |
| NoCopy | boolean | Container create options as would be passed to `docker run` |  |

*****

## DriverConfig

#### Type: object

__Container create options as would be passed to `docker run`__

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| Name | string | Container create options as would be passed to `docker run` |  |
| [Options](#options) | object | Container create options as would be passed to `docker run` |  |

*****

## Options

#### Type: object

__Container create options as would be passed to `docker run`__

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| [.*](#.*) | string | Container create options as would be passed to `docker run` |  |

*****

## PortBindings

#### Type: object

__Container create options as would be passed to `docker run`__

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| [.*](#.*) | array | Container create options as would be passed to `docker run` |  |

*****

## RestartPolicy

#### Type: object

__Container create options as would be passed to `docker run`__

Required: Name, MaximumRetryCount, Name, MaximumRetryCount

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| MaximumRetryCount | integer | Container create options as would be passed to `docker run` |  |
| Name | string | Container create options as would be passed to `docker run` |  |

*****

## SecurityOpt:undefined

#### Type: string

__Container create options as would be passed to `docker run`__


*****

## StorageOpt

#### Type: object

__Container create options as would be passed to `docker run`__

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| [.*](#.*) | string | Container create options as would be passed to `docker run` |  |

*****

## Sysctls

#### Type: object

__Container create options as would be passed to `docker run`__

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| [.*](#.*) | string | Container create options as would be passed to `docker run` |  |

*****

## Tmpfs

#### Type: object

__Container create options as would be passed to `docker run`__

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| [.*](#.*) | string | Container create options as would be passed to `docker run` |  |

*****

## VolumesFrom:undefined

#### Type: string

__Container create options as would be passed to `docker run`__


*****

## NetworkingConfig

#### Type: object

__Container create options as would be passed to `docker run`__

Required: EndpointsConfig, EndpointsConfig

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| [EndpointsConfig](#endpointsconfig) | object | Container create options as would be passed to `docker run` |  |

*****

## EndpointsConfig

#### Type: object

__Container create options as would be passed to `docker run`__


*****

## docker.exec

#### Type: object

__Collect the stdout/stderr of executing a command in an already running docker container__

Required: container, exec_config

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| container | string | The container name |  |
| [exec_config](#exec_config) | object | Config options as would be passed to `docker exec`  |  |

*****

## docker.image-ls

#### Type: object

__Collect a list of docker images present on the server__

Required: All, Filters

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| All | boolean | Same as would be passed to `docker images` |  |
| [Filters](#filters) | object | Same as would be passed to `docker images` |  |

*****

## docker.images

#### Type: object

__Collect a list of docker images present on the server__

Required: All, Filters

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| All | boolean | Same as would be passed to `docker images` |  |
| [Filters](#filters) | object | Same as would be passed to `docker images` |  |

*****

## docker.info

#### Type: object

__Collect info about the Docker daemon__


*****

## docker.logs

#### Type: object

__Collect info about the Docker daemon__

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| container | string | A container name |  |
| [container_list_options](#container_list_options) | object | Options to be used to filter the list of containers |  |
| [container_logs_options](#container_logs_options) | object | Args as would be given to `docker logs` |  |

*****

## docker.node-ls

#### Type: object

__Collect information about the nodes in a Docker Swarm installation__

Required: Filters

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| [Filters](#filters) | object | Same as would be passed to `docker node ls` |  |

*****

## docker.ps

#### Type: object

__Collect information about containers__

Required: Quiet, Size, All, Latest, Since, Before, Limit, Filters

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| All | boolean | Same as `-a` option to `docker ps`  |  |
| Before | string | Same as `--before` option to `docker logs`  |  |
| [Filters](#filters) | object | Same as `--filter` option to `docker ps`  |  |
| Latest | boolean | Same as would be passed to `docker ps` |  |
| Limit | integer | Same as would be passed to `docker ps` |  |
| Quiet | boolean | Same as would be passed to `docker ps` |  |
| Since | string | Same as would be passed to `docker ps` |  |
| Size | boolean | Same as would be passed to `docker ps` |  |

*****

## docker.run

#### Type: object

__Collect the stderr/stdout of running a single docker container__

Required: container_create_config

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| [container_create_config](#container_create_config) | object | Container create options as would be passed to `docker run` |  |
| enable_pull | boolean | If `true`, allow this container to be pulled if not present |  |

*****

## docker.service-logs

#### Type: object

__Collect logs from a docker swarm service__

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| [container_logs_options](#container_logs_options) | object | Args as would be given to `docker logs` |  |
| service | string | A single docker swarm service for which to collect logs |  |
| [service_list_options](#service_list_options) | object | Options for filtering all swarm services |  |

*****

## service_list_options

#### Type: object

__Options for filtering all swarm servicesundefinedundefined__

Required: Filters, Filters, Filters

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| [Filters](#filters) | object |  |  |

*****

## docker.service-ls

#### Type: object

__Collect a list of docker swarm services__

Required: Filters

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| [Filters](#filters) | object | Same as would be passed to `docker service ls` |  |

*****

## docker.service-ps

#### Type: object

__Collect information about the tasks run by one or more services__

Required: Filters

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| [Filters](#filters) | object | Same as would be passed to `docker service ps` |  |

*****

## docker.stack-service-logs

#### Type: object

__Collect logs from one or more services in a stack__

Required: namespace

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| [container_logs_options](#container_logs_options) | object | Args as would be given to `docker logs` |  |
| namespace | string |  |  |
| [service_list_options](#service_list_options) | object |  |  |

*****

## docker.stack-service-ls

#### Type: object

__add description to json file__

Required: namespace

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| namespace | string |  |  |
| [service_list_options](#service_list_options) | object |  |  |

*****

## docker.stack-service-ps

#### Type: object

__add description to json file__

Required: namespace

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| namespace | string |  |  |
| [task_list_options](#task_list_options) | object |  |  |

*****

## task_list_options

#### Type: object

__add description to json file__

Required: Filters, Filters, Filters

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| [Filters](#filters) | object |  |  |

*****

## docker.stack-task-logs

#### Type: object

__add description to json file__

Required: namespace

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| [container_logs_options](#container_logs_options) | object |  |  |
| namespace | string |  |  |
| [task_list_options](#task_list_options) | object |  |  |

*****

## docker.task-logs

#### Type: object

__add description to json file__

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| [container_logs_options](#container_logs_options) | object |  |  |
| id | string |  |  |
| [task_list_options](#task_list_options) | object |  |  |

*****

## docker.task-ls

#### Type: object

__add description to json file__

Required: Filters

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| [Filters](#filters) | object |  |  |

*****

## docker.version

#### Type: object

__add description to json file__


*****

## journald.logs

#### Type: object

__add description to json file__

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| reverse | boolean |  |  |
| since | string |  |  |
| unit | string |  |  |

*****

## kubernetes.api-versions

#### Type: object

__add description to json file__


*****

## kubernetes.cluster-info

#### Type: object

__add description to json file__


*****

## kubernetes.logs

#### Type: object

__add description to json file__

Required: pod

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| namespace | string |  |  |
| pod | string |  |  |
| [pod_log_options](#pod_log_options) | object |  |  |

*****

## pod_log_options

#### Type: object

__add description to json file__

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| apiVersion | string |  |  |
| container | string |  |  |
| follow | boolean |  |  |
| kind | string |  |  |
| limitBytes | integer |  |  |
| previous | boolean |  |  |
| sinceSeconds | integer |  |  |
| [sinceTime](#sincetime) | object |  |  |
| tailLines | integer |  |  |
| timestamps | boolean |  |  |

*****

## sinceTime

#### Type: object

__add description to json file__


*****

## kubernetes.resource-list

#### Type: object

__add description to json file__

Required: kind

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| kind | string |  |  |
| namespace | string |  |  |
| [resource_list_options](#resource_list_options) | object |  |  |

*****

## resource_list_options

#### Type: object

__add description to json file__

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| apiVersion | string |  |  |
| continue | string |  |  |
| fieldSelector | string |  |  |
| includeUninitialized | boolean |  |  |
| kind | string |  |  |
| labelSelector | string |  |  |
| limit | integer |  |  |
| resourceVersion | string |  |  |
| timeoutSeconds | integer |  |  |
| watch | boolean |  |  |

*****

## kubernetes.version

#### Type: object

__add description to json file__


*****

## meta

#### Type: object

__add description to json file__

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| [labels](#labels) | object |  |  |
| name | string |  |  |

*****

## labels

#### Type: object

__add description to json file__

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| [.*](#.*) | string |  |  |

*****

## meta.customer

#### Type: object

__add description to json file__

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| customer_id | string |  |  |
| customer_name | string |  |  |
| installation_id | string |  |  |

*****

## os.hostname

#### Type: object

__add description to json file__


*****

## os.http-request

#### Type: object

__add description to json file__

Required: url

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| body | string |  |  |
| [header](#header) | object |  |  |
| insecure | boolean |  |  |
| method | string |  |  |
| url | string |  |  |

*****

## header

#### Type: object

__add description to json file__

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| [.*](#.*) | array |  |  |

*****

## os.loadavg

#### Type: object

__add description to json file__

Required: template

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| template | string |  |  |

*****

## os.read-file

#### Type: object

__add description to json file__

Required: filepath

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| filepath | string |  |  |

*****

## os.run-command

#### Type: object

__add description to json file__

Required: name

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| args | array |  |  |
| dir | string |  |  |
| env | array |  |  |
| name | string |  |  |

*****

## args:undefined

#### Type: string

__add description to json file__


*****

## env:undefined

#### Type: string

__add description to json file__


*****

## os.uptime

#### Type: object

__add description to json file__

Required: template

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| template | string |  |  |

*****

## retraced.events

#### Type: object

__add description to json file__

Required: api_endpoint

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| api_endpoint | string |  |  |
| api_token | string |  |  |
| insecure | boolean |  |  |
| [mask](#mask) | object |  |  |
| project_id | string |  |  |
| [query](#query) | object |  |  |

*****

## mask

#### Type: object

__add description to json file__

Required: ID, Action, CRUD, Description, IsFailure, IsAnonymous, SourceIP, Country, LocSubdiv1, LocSubdiv2, Received, Created, CanonicalTime, Component, Version, Fields, Raw, GroupID, GroupName, ActorID, ActorName, ActorHref, ActorFields, TargetID, TargetName, TargetHref, TargetType, TargetFields, DisplayMarkdown

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| Action | boolean |  |  |
| ActorFields | boolean |  |  |
| ActorHref | boolean |  |  |
| ActorID | boolean |  |  |
| ActorName | boolean |  |  |
| CRUD | boolean |  |  |
| CanonicalTime | boolean |  |  |
| Component | boolean |  |  |
| Country | boolean |  |  |
| Created | boolean |  |  |
| Description | boolean |  |  |
| DisplayMarkdown | boolean |  |  |
| Fields | boolean |  |  |
| GroupID | boolean |  |  |
| GroupName | boolean |  |  |
| ID | boolean |  |  |
| IsAnonymous | boolean |  |  |
| IsFailure | boolean |  |  |
| LocSubdiv1 | boolean |  |  |
| LocSubdiv2 | boolean |  |  |
| Raw | boolean |  |  |
| Received | boolean |  |  |
| SourceIP | boolean |  |  |
| TargetFields | boolean |  |  |
| TargetHref | boolean |  |  |
| TargetID | boolean |  |  |
| TargetName | boolean |  |  |
| TargetType | boolean |  |  |
| Version | boolean |  |  |

*****

## query

#### Type: object

__add description to json file__

Required: Action, CRUD, ReceivedStart, ReceivedEnd, CreatedStart, CreatedEnd, ActorName, ActorID, Description, Location

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| Action | string |  |  |
| ActorID | string |  |  |
| ActorName | string |  |  |
| CRUD | string |  |  |
| CreatedEnd | string |  |  |
| CreatedStart | string |  |  |
| Description | string |  |  |
| Location | string |  |  |
| ReceivedEnd | string |  |  |
| ReceivedStart | string |  |  |

*****

## scrub

#### Type: object

__add description to json file__

Required: regex, replace

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| regex | string |  |  |
| replace | string |  |  |

*****

## version

#### Type: object

__add description to json file__


*****
