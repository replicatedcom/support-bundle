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
| [message](#message) | string | message to display before prompting |  |

*****

## specs:undefined

#### Type: object

__A list of items to include in the support bundle.__

Required: output_dir

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| description | string | A description of the file(s) that will be collected by this spec |  |
| [docker.container-cp](#docker.container-cp) | object | Copy a file from a running docker container |  |
| [docker.container-exec](#docker.container-exec) | object | Collect the stdout and stderr of `exec`-ing a command on a running docker container |  |
| [docker.container-inspect](#docker.container-inspect) | object | inspect one or more running or stopped containers |  |
| [docker.container-logs](#docker.container-logs) | object | Collect the logs from one or more docker containers |  |
| [docker.container-ls](#docker.container-ls) | object | List information about one or more containers |  |
| [docker.container-run](#docker.container-run) | object | Run a docker container, generating a file from the stdout and/or stderr |  |
| [docker.exec](#docker.exec) | object |  |  |
| [docker.image-ls](#docker.image-ls) | object |  |  |
| [docker.images](#docker.images) | object |  |  |
| [docker.info](#docker.info) | object |  |  |
| [docker.logs](#docker.logs) | object |  |  |
| [docker.node-ls](#docker.node-ls) | object |  |  |
| [docker.ps](#docker.ps) | object |  |  |
| [docker.run](#docker.run) | object |  |  |
| [docker.service-logs](#docker.service-logs) | object |  |  |
| [docker.service-ls](#docker.service-ls) | object |  |  |
| [docker.service-ps](#docker.service-ps) | object |  |  |
| [docker.stack-service-logs](#docker.stack-service-logs) | object |  |  |
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

__Copy a file from a running docker container__

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
| [exec_config](#exec_config) | object | The arguments as would be passed to `docker exec` |  |

*****

## exec_config

#### Type: object

__The arguments as would be passed to `docker exec`undefined__

Required: User, Privileged, Tty, AttachStdin, AttachStderr, AttachStdout, Detach, DetachKeys, Env, Cmd, User, Privileged, Tty, AttachStdin, AttachStderr, AttachStdout, Detach, DetachKeys, Env, Cmd

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| AttachStderr | boolean |  |  |
| AttachStdin | boolean |  |  |
| AttachStdout | boolean |  |  |
| Cmd | array |  |  |
| Detach | boolean |  |  |
| DetachKeys | string |  |  |
| Env | array |  |  |
| Privileged | boolean |  |  |
| Tty | boolean |  |  |
| User | string |  |  |

*****

## Cmd:undefined

#### Type: string

__same as would be passed to `docker exec`undefinedundefinedundefined__


*****

## Env:undefined

#### Type: string

__same as would be passed to `docker exec`undefinedundefinedundefined__


*****

## docker.container-inspect

#### Type: object

__inspect one or more running or stopped containers__

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| container | string | the container name |  |
| [container_list_options](#container_list_options) | object | Options to be used to filter the list of containers |  |

*****

## container_list_options

#### Type: object

__Options to be used to filter the list of containersundefined__

Required: Quiet, Size, All, Latest, Since, Before, Limit, Filters, Quiet, Size, All, Latest, Since, Before, Limit, Filters, Quiet, Size, All, Latest, Since, Before, Limit, Filters

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| All | boolean |  |  |
| Before | string |  |  |
| [Filters](#filters) | object |  |  |
| Latest | boolean |  |  |
| Limit | integer |  |  |
| Quiet | boolean |  |  |
| Since | string |  |  |
| Size | boolean |  |  |

*****

## Filters

#### Type: object

__Same as `--filter` option to `docker ps` undefinedundefinedundefinedundefinedundefinedundefinedundefinedundefinedundefinedundefinedundefinedundefinedundefinedundefined__

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

__Args as would be given to `docker logs`undefinedundefinedundefinedundefinedundefined__

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

__List information about one or more containers__

Required: Quiet, Size, All, Latest, Since, Before, Limit, Filters

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| All | boolean | Same as `-a` option to `docker ps`  |  |
| Before | string | Same as `--before` option to `docker logs`  |  |
| [Filters](#filters) | object | Same as `--filter` option to `docker ps`  |  |
| Latest | boolean | Same as docker cli arg |  |
| Limit | integer | Same as docker cli arg |  |
| Quiet | boolean | Same as docker cli arg |  |
| Since | string | Same as docker cli arg |  |
| Size | boolean | Same as docker cli arg |  |

*****

## docker.container-run

#### Type: object

__Run a docker container, generating a file from the stdout and/or stderr__

Required: container_create_config

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| [container_create_config](#container_create_config) | object |  |  |
| enable_pull | boolean |  |  |

*****

## container_create_config

#### Type: object

__add description to json file__

Required: Name, Config, HostConfig, NetworkingConfig, AdjustCPUShares, Name, Config, HostConfig, NetworkingConfig, AdjustCPUShares

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| AdjustCPUShares | boolean |  |  |
| [Config](#config) | object |  |  |
| [HostConfig](#hostconfig) | object |  |  |
| Name | string |  |  |
| [NetworkingConfig](#networkingconfig) | object |  |  |

*****

## Config

#### Type: object

__add description to json file__

Required: Hostname, Domainname, User, AttachStdin, AttachStdout, AttachStderr, Tty, OpenStdin, StdinOnce, Env, Cmd, Image, Volumes, WorkingDir, Entrypoint, OnBuild, Labels, Hostname, Domainname, User, AttachStdin, AttachStdout, AttachStderr, Tty, OpenStdin, StdinOnce, Env, Cmd, Image, Volumes, WorkingDir, Entrypoint, OnBuild, Labels

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| ArgsEscaped | boolean |  |  |
| AttachStderr | boolean |  |  |
| AttachStdin | boolean |  |  |
| AttachStdout | boolean |  |  |
| Cmd | array |  |  |
| Domainname | string |  |  |
| Entrypoint | array |  |  |
| Env | array |  |  |
| [ExposedPorts](#exposedports) | object |  |  |
| [Healthcheck](#healthcheck) | object |  |  |
| Hostname | string |  |  |
| Image | string |  |  |
| [Labels](#labels) | object |  |  |
| MacAddress | string |  |  |
| NetworkDisabled | boolean |  |  |
| OnBuild | array |  |  |
| OpenStdin | boolean |  |  |
| Shell | array |  |  |
| StdinOnce | boolean |  |  |
| StopSignal | string |  |  |
| StopTimeout | integer |  |  |
| Tty | boolean |  |  |
| User | string |  |  |
| [Volumes](#volumes) | object |  |  |
| WorkingDir | string |  |  |
| [.*](#.*) | string |  |  |

*****

## Entrypoint:undefined

#### Type: string

__add description to json file__


*****

## ExposedPorts

#### Type: object

__add description to json file__

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| [.*](#.*) | object |  |  |

*****

## .*

#### Type: object

__add description to json file__


*****

## Healthcheck

#### Type: object

__add description to json file__

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| Interval | integer |  |  |
| Retries | integer |  |  |
| StartPeriod | integer |  |  |
| Test | array |  |  |
| Timeout | integer |  |  |

*****

## Test:undefined

#### Type: string

__add description to json file__


*****

## Labels

#### Type: object

__add description to json file__

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| [.*](#.*) | string |  |  |

*****

## OnBuild:undefined

#### Type: string

__add description to json file__


*****

## Shell:undefined

#### Type: string

__add description to json file__


*****

## Volumes

#### Type: object

__add description to json file__

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| [.*](#.*) | object |  |  |

*****

## HostConfig

#### Type: object

__add description to json file__

Required: Binds, ContainerIDFile, LogConfig, NetworkMode, PortBindings, RestartPolicy, AutoRemove, VolumeDriver, VolumesFrom, CapAdd, CapDrop, Dns, DnsOptions, DnsSearch, ExtraHosts, GroupAdd, IpcMode, Cgroup, Links, OomScoreAdj, PidMode, Privileged, PublishAllPorts, ReadonlyRootfs, SecurityOpt, UTSMode, UsernsMode, ShmSize, ConsoleSize, Isolation, CpuShares, Memory, NanoCpus, CgroupParent, BlkioWeight, BlkioWeightDevice, BlkioDeviceReadBps, BlkioDeviceWriteBps, BlkioDeviceReadIOps, BlkioDeviceWriteIOps, CpuPeriod, CpuQuota, CpuRealtimePeriod, CpuRealtimeRuntime, CpusetCpus, CpusetMems, Devices, DeviceCgroupRules, DiskQuota, KernelMemory, MemoryReservation, MemorySwap, MemorySwappiness, OomKillDisable, PidsLimit, Ulimits, CpuCount, CpuPercent, IOMaximumIOps, IOMaximumBandwidth, Binds, ContainerIDFile, LogConfig, NetworkMode, PortBindings, RestartPolicy, AutoRemove, VolumeDriver, VolumesFrom, CapAdd, CapDrop, Dns, DnsOptions, DnsSearch, ExtraHosts, GroupAdd, IpcMode, Cgroup, Links, OomScoreAdj, PidMode, Privileged, PublishAllPorts, ReadonlyRootfs, SecurityOpt, UTSMode, UsernsMode, ShmSize, ConsoleSize, Isolation, CpuShares, Memory, NanoCpus, CgroupParent, BlkioWeight, BlkioWeightDevice, BlkioDeviceReadBps, BlkioDeviceWriteBps, BlkioDeviceReadIOps, BlkioDeviceWriteIOps, CpuPeriod, CpuQuota, CpuRealtimePeriod, CpuRealtimeRuntime, CpusetCpus, CpusetMems, Devices, DeviceCgroupRules, DiskQuota, KernelMemory, MemoryReservation, MemorySwap, MemorySwappiness, OomKillDisable, PidsLimit, Ulimits, CpuCount, CpuPercent, IOMaximumIOps, IOMaximumBandwidth

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| AutoRemove | boolean |  |  |
| Binds | array |  |  |
| BlkioDeviceReadBps | array |  |  |
| BlkioDeviceReadIOps | array |  |  |
| BlkioDeviceWriteBps | array |  |  |
| BlkioDeviceWriteIOps | array |  |  |
| BlkioWeight | integer |  |  |
| BlkioWeightDevice | array |  |  |
| CapAdd | array |  |  |
| CapDrop | array |  |  |
| Cgroup | string |  |  |
| CgroupParent | string |  |  |
| ConsoleSize | undefined |  |  |
| ContainerIDFile | string |  |  |
| CpuCount | integer |  |  |
| CpuPercent | integer |  |  |
| CpuPeriod | integer |  |  |
| CpuQuota | integer |  |  |
| CpuRealtimePeriod | integer |  |  |
| CpuRealtimeRuntime | integer |  |  |
| CpuShares | integer |  |  |
| CpusetCpus | string |  |  |
| CpusetMems | string |  |  |
| DeviceCgroupRules | array |  |  |
| Devices | array |  |  |
| DiskQuota | integer |  |  |
| Dns | array |  |  |
| DnsOptions | array |  |  |
| DnsSearch | array |  |  |
| ExtraHosts | array |  |  |
| GroupAdd | array |  |  |
| IOMaximumBandwidth | integer |  |  |
| IOMaximumIOps | integer |  |  |
| Init | boolean |  |  |
| IpcMode | string |  |  |
| Isolation | string |  |  |
| KernelMemory | integer |  |  |
| Links | array |  |  |
| [LogConfig](#logconfig) | object |  |  |
| Memory | integer |  |  |
| MemoryReservation | integer |  |  |
| MemorySwap | integer |  |  |
| MemorySwappiness | integer |  |  |
| Mounts | array |  |  |
| NanoCpus | integer |  |  |
| NetworkMode | string |  |  |
| OomKillDisable | boolean |  |  |
| OomScoreAdj | integer |  |  |
| PidMode | string |  |  |
| PidsLimit | integer |  |  |
| [PortBindings](#portbindings) | object |  |  |
| Privileged | boolean |  |  |
| PublishAllPorts | boolean |  |  |
| ReadonlyRootfs | boolean |  |  |
| [RestartPolicy](#restartpolicy) | object |  |  |
| Runtime | string |  |  |
| SecurityOpt | array |  |  |
| ShmSize | integer |  |  |
| [StorageOpt](#storageopt) | object |  |  |
| [Sysctls](#sysctls) | object |  |  |
| [Tmpfs](#tmpfs) | object |  |  |
| UTSMode | string |  |  |
| Ulimits | array |  |  |
| UsernsMode | string |  |  |
| VolumeDriver | string |  |  |
| VolumesFrom | array |  |  |

*****

## Binds:undefined

#### Type: string

__add description to json file__


*****

## CapAdd:undefined

#### Type: string

__add description to json file__


*****

## CapDrop:undefined

#### Type: string

__add description to json file__


*****

## DeviceCgroupRules:undefined

#### Type: string

__add description to json file__


*****

## Devices:undefined

#### Type: object

__add description to json file__

Required: PathOnHost, PathInContainer, CgroupPermissions, PathOnHost, PathInContainer, CgroupPermissions

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| CgroupPermissions | string |  |  |
| PathInContainer | string |  |  |
| PathOnHost | string |  |  |

*****

## Dns:undefined

#### Type: string

__add description to json file__


*****

## DnsOptions:undefined

#### Type: string

__add description to json file__


*****

## DnsSearch:undefined

#### Type: string

__add description to json file__


*****

## ExtraHosts:undefined

#### Type: string

__add description to json file__


*****

## GroupAdd:undefined

#### Type: string

__add description to json file__


*****

## Links:undefined

#### Type: string

__add description to json file__


*****

## LogConfig

#### Type: object

__add description to json file__

Required: Type, Config, Type, Config

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| [Config](#config) | object |  |  |
| Type | string |  |  |

*****

## Mounts:undefined

#### Type: object

__add description to json file__

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| [BindOptions](#bindoptions) | object |  |  |
| Consistency | string |  |  |
| ReadOnly | boolean |  |  |
| Source | string |  |  |
| Target | string |  |  |
| [TmpfsOptions](#tmpfsoptions) | object |  |  |
| Type | string |  |  |
| [VolumeOptions](#volumeoptions) | object |  |  |

*****

## BindOptions

#### Type: object

__add description to json file__

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| Propagation | string |  |  |

*****

## TmpfsOptions

#### Type: object

__add description to json file__

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| Mode | integer |  |  |
| SizeBytes | integer |  |  |

*****

## VolumeOptions

#### Type: object

__add description to json file__

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| [DriverConfig](#driverconfig) | object |  |  |
| [Labels](#labels) | object |  |  |
| NoCopy | boolean |  |  |

*****

## DriverConfig

#### Type: object

__add description to json file__

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| Name | string |  |  |
| [Options](#options) | object |  |  |

*****

## Options

#### Type: object

__add description to json file__

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| [.*](#.*) | string |  |  |

*****

## PortBindings

#### Type: object

__add description to json file__

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| [.*](#.*) | array |  |  |

*****

## RestartPolicy

#### Type: object

__add description to json file__

Required: Name, MaximumRetryCount, Name, MaximumRetryCount

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| MaximumRetryCount | integer |  |  |
| Name | string |  |  |

*****

## SecurityOpt:undefined

#### Type: string

__add description to json file__


*****

## StorageOpt

#### Type: object

__add description to json file__

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| [.*](#.*) | string |  |  |

*****

## Sysctls

#### Type: object

__add description to json file__

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| [.*](#.*) | string |  |  |

*****

## Tmpfs

#### Type: object

__add description to json file__

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| [.*](#.*) | string |  |  |

*****

## VolumesFrom:undefined

#### Type: string

__add description to json file__


*****

## NetworkingConfig

#### Type: object

__add description to json file__

Required: EndpointsConfig, EndpointsConfig

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| [EndpointsConfig](#endpointsconfig) | object |  |  |

*****

## EndpointsConfig

#### Type: object

__add description to json file__


*****

## docker.exec

#### Type: object

__add description to json file__

Required: container, exec_config

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| container | string |  |  |
| [exec_config](#exec_config) | object |  |  |

*****

## docker.image-ls

#### Type: object

__add description to json file__

Required: All, Filters

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| All | boolean |  |  |
| [Filters](#filters) | object |  |  |

*****

## docker.images

#### Type: object

__add description to json file__

Required: All, Filters

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| All | boolean |  |  |
| [Filters](#filters) | object |  |  |

*****

## docker.info

#### Type: object

__add description to json file__


*****

## docker.logs

#### Type: object

__add description to json file__

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| container | string |  |  |
| [container_list_options](#container_list_options) | object |  |  |
| [container_logs_options](#container_logs_options) | object |  |  |

*****

## docker.node-ls

#### Type: object

__add description to json file__

Required: Filters

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| [Filters](#filters) | object |  |  |

*****

## docker.ps

#### Type: object

__add description to json file__

Required: Quiet, Size, All, Latest, Since, Before, Limit, Filters

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| All | boolean |  |  |
| Before | string |  |  |
| [Filters](#filters) | object |  |  |
| Latest | boolean |  |  |
| Limit | integer |  |  |
| Quiet | boolean |  |  |
| Since | string |  |  |
| Size | boolean |  |  |

*****

## docker.run

#### Type: object

__add description to json file__

Required: container_create_config

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| [container_create_config](#container_create_config) | object |  |  |
| enable_pull | boolean |  |  |

*****

## docker.service-logs

#### Type: object

__add description to json file__

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| [container_logs_options](#container_logs_options) | object |  |  |
| service | string |  |  |
| [service_list_options](#service_list_options) | object |  |  |

*****

## service_list_options

#### Type: object

__add description to json file__

Required: Filters, Filters, Filters

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| [Filters](#filters) | object |  |  |

*****

## docker.service-ls

#### Type: object

__add description to json file__

Required: Filters

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| [Filters](#filters) | object |  |  |

*****

## docker.service-ps

#### Type: object

__add description to json file__

Required: Filters

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| [Filters](#filters) | object |  |  |

*****

## docker.stack-service-logs

#### Type: object

__add description to json file__

Required: namespace

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| [container_logs_options](#container_logs_options) | object |  |  |
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
