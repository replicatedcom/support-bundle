# 

#### Type: object

__add description to json file__

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| lifecycle | array | `lifecycle` can be used to customize the end-customer's experience when generating/uploading support bundles |  |
| specs | array | A list of items to include in the support bundle |  |

*****

## lifecycle

#### Type: object

__`lifecycle` can be used to customize the end-customer's experience when generating/uploading support bundles__

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| [generate](#generate) | object | Adds a step that will generate a support bundle. It has no arguments, but should be specified with an empty object `{}` |  |
| [message](#message) | object | `message` can be used to print a message to the end customer |  |
| [upload](#upload) | object | Adds a step in which the end customer can upload a bundle to https://console.replicated.com for review by the vendor |  |

*****

## generate

#### Type: object

__Adds a step that will generate a support bundle. It has no arguments, but should be specified with an empty object `{}`__


*****

## message

#### Type: object

__`message` can be used to print a message to the end customer__

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| contents | string | `contents` is the message |  |

*****

## upload

#### Type: object

__Adds a step in which the end customer can upload a bundle to https://console.replicated.com for review by the vendor__

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

## specs

#### Type: object

__A list of items to include in the support bundle__

Required: 

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| description | string | A description of the file(s) that will be collected by this spec |  |
| [docker.container-cp](#docker.container-cp) | object | Collect a file by copying from a running docker container |  |
| [docker.container-exec](#docker.container-exec) | object | Collect the stdout and stderr of `exec`-ing a command on a running docker container |  |
| [docker.container-inspect](#docker.container-inspect) | object | Collect the `docker inspect` output for one or more running or stopped containers. One of `container` or `container_list_options` is required. |  |
| [docker.container-logs](#docker.container-logs) | object | Collect the stdout/stderr logs from one or more docker containers. One of `container` or `container_list_options` is required. |  |
| [docker.container-ls](#docker.container-ls) | object | Collect information about one or more containers |  |
| [docker.container-run](#docker.container-run) | object | Collect the stderr/stdout of running a single docker container |  |
| [docker.exec](#docker.exec) | object | Collect the stdout/stderr of executing a command in an already running docker container |  |
| [docker.image-ls](#docker.image-ls) | object | Collect a list of docker images present on the server |  |
| [docker.images](#docker.images) | object | Collect a list of docker images present on the server |  |
| [docker.info](#docker.info) | object | Collect info about the Docker daemon |  |
| [docker.logs](#docker.logs) | object | Collect the stdout/stderr logs from one or more docker containers. One of `container` or `container_list_options` is required. |  |
| [docker.node-ls](#docker.node-ls) | object | Collect information about the nodes in a Docker Swarm installation |  |
| [docker.ps](#docker.ps) | object | Collect information about containers |  |
| [docker.run](#docker.run) | object | Collect the stderr/stdout of running a single docker container |  |
| [docker.service-logs](#docker.service-logs) | object | Collect logs from a docker swarm service. One of `service` or `service_list_options` is required. |  |
| [docker.service-ls](#docker.service-ls) | object | Collect a list of docker swarm services |  |
| [docker.service-ps](#docker.service-ps) | object | Collect information about the tasks run by one or more services |  |
| [docker.stack-service-logs](#docker.stack-service-logs) | object | Collect logs from one or more services in a stack |  |
| [docker.stack-service-ls](#docker.stack-service-ls) | object | Collect information about services in a swarm stack |  |
| [docker.stack-service-ps](#docker.stack-service-ps) | object | Collect information about the tasks running in a service |  |
| [docker.stack-task-logs](#docker.stack-task-logs) | object | Collect logs from a docker swarm task |  |
| [docker.task-logs](#docker.task-logs) | object | Collect logs from a docker swarm task |  |
| [docker.task-ls](#docker.task-ls) | object | List docker swarm tasks |  |
| [docker.version](#docker.version) | object | Get the version of the docker server |  |
| [journald.logs](#journald.logs) | object | Collect logs from journald |  |
| [meta](#meta) | object |  |  |
| [os.hostname](#os.hostname) | object | Get the hostname |  |
| [os.http-request](#os.http-request) | object | Collect information about the response from making an HTTP request |  |
| [os.loadavg](#os.loadavg) | object | Collect information about CPU load |  |
| [os.read-file](#os.read-file) | object | Read a file from the filesystem |  |
| [os.run-command](#os.run-command) | object | Run a command. Will be run in the support bundle container if `docker run replicated/support-bundle` is used. To run in a custom container, use `docker.run` instead |  |
| [os.uptime](#os.uptime) | object | Collect information about the system's uptime |  |
| output_dir | string |  |  |
| [retraced.events](#retraced.events) | object | Collect Audit Log events from a running Retraced instance |  |
| [scrub](#scrub) | object |  |  |
| timeout_seconds | integer |  |  |

*****

## docker.container-cp

#### Type: object

__Collect a file by copying from a running docker container__

Required: container, src_path, output_dir

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| container | string | The name of the container to copy from | my-sweet-app |
| description | string |  |  |
| [meta](#meta) | object |  |  |
| output_dir | string |  |  |
| [scrub](#scrub) | object |  |  |
| src_path | string | The path of the target file in the container's filesystem | /var/log/supergoodtool/access.log |
| timeout_seconds | integer |  |  |

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

## scrub

#### Type: object

__add description to json file__

Required: regex, replace, regex, replace, regex, replace, regex, replace, regex, replace, regex, replace, regex, replace, regex, replace, regex, replace, regex, replace, regex, replace, regex, replace, regex, replace, regex, replace, regex, replace, regex, replace, regex, replace, regex, replace, regex, replace, regex, replace, regex, replace, regex, replace, regex, replace, regex, replace, regex, replace, regex, replace, regex, replace, regex, replace, regex, replace, regex, replace, regex, replace, regex, replace, regex, replace, regex, replace, regex, replace

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| regex | string |  |  |
| replace | string |  |  |

*****

## docker.container-exec

#### Type: object

__Collect the stdout and stderr of `exec`-ing a command on a running docker container__

Required: output_dir, container, exec_config

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| container | string | The name of the container to run the command in |  |
| description | string |  |  |
| [exec_config](#exec_config) | object | Same as would be passed to `docker exec`, as in [The Docker API](https://github.com/moby/moby/blob/master/api/types/configs.go#L43) |  |
| [meta](#meta) | object |  |  |
| output_dir | string |  |  |
| [scrub](#scrub) | object |  |  |
| timeout_seconds | integer |  |  |

*****

## exec_config

#### Type: object

__Same as would be passed to `docker exec`, as in [The Docker API](https://github.com/moby/moby/blob/master/api/types/configs.go#L43)Config options as would be passed to `docker exec`, as in [The Docker API](https://github.com/moby/moby/blob/master/api/types/configs.go)__

Required: Cmd, Cmd

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

## Cmd

#### Type: string

__Same as would be passed to `docker exec`Container create options as would be passed to `docker run`Same as would be passed to `docker exec`Container create options as would be passed to `docker run`__


*****

## Env

#### Type: string

__Same as would be passed to `docker exec`Container create options as would be passed to `docker run`Same as would be passed to `docker exec`Container create options as would be passed to `docker run`__


*****

## docker.container-inspect

#### Type: object

__Collect the `docker inspect` output for one or more running or stopped containers. One of `container` or `container_list_options` is required.__

Required: output_dir

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| container | string | the container name |  |
| [container_list_options](#container_list_options) | object | Options to be used to filter the list of containers, as in [The Docker API](https://github.com/moby/moby/blob/master/api/types/client.go#L61) |  |
| description | string |  |  |
| [meta](#meta) | object |  |  |
| output_dir | string |  |  |
| [scrub](#scrub) | object |  |  |
| timeout_seconds | integer |  |  |

*****

## container_list_options

#### Type: object

__Options to be used to filter the list of containers, as in [The Docker API](https://github.com/moby/moby/blob/master/api/types/client.go#L61)__

Required: 

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| All | boolean | Option to be used to filter the list of containers, as in [The Docker API](https://github.com/moby/moby/blob/master/api/types/client.go#L61) |  |
| Before | string | Same as `--before` option to `docker logs`  |  |
| [Filters](#filters) | object | Option to be used to filter the list of containers, as in [The Docker API](https://github.com/moby/moby/blob/master/api/types/client.go#L61) |  |
| Latest | boolean | Option to be used to filter the list of containers, as in [The Docker API](https://github.com/moby/moby/blob/master/api/types/client.go#L61) |  |
| Limit | integer | Option to be used to filter the list of containers, as in [The Docker API](https://github.com/moby/moby/blob/master/api/types/client.go#L61) |  |
| Quiet | boolean | Option to be used to filter the list of containers, as in [The Docker API](https://github.com/moby/moby/blob/master/api/types/client.go#L61) |  |
| Since | string | Option to be used to filter the list of containers, as in [The Docker API](https://github.com/moby/moby/blob/master/api/types/client.go#L61) |  |
| Size | boolean | Option to be used to filter the list of containers, as in [The Docker API](https://github.com/moby/moby/blob/master/api/types/client.go#L61) |  |

*****

## Filters

#### Type: object

__Option to be used to filter the list of containers, as in [The Docker API](https://github.com/moby/moby/blob/master/api/types/client.go#L61)Same as would be passed to `docker images`Same as would be passed to `docker images`Option to be used to filter the list of containers, as in [The Docker API](https://github.com/moby/moby/blob/master/api/types/client.go#L61)Same as would be passed to `docker node ls`Option to be used to filter the list of containers, as in [The Docker API](https://github.com/moby/moby/blob/master/api/types/client.go#L61)Same as would be passed to `docker service ls`Same as would be passed to `docker service ls`Same as would be passed to `docker service ps`Filters for tasks__

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| [.*](#.*) | array |  |  |

*****

## docker.container-logs

#### Type: object

__Collect the stdout/stderr logs from one or more docker containers. One of `container` or `container_list_options` is required.__

Required: output_dir

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| container | string | A container name |  |
| [container_list_options](#container_list_options) | object | Options to be used to filter the list of containers, as in [The Docker API](https://github.com/moby/moby/blob/master/api/types/client.go#L61) |  |
| [container_logs_options](#container_logs_options) | object | Args as would be given to `docker logs`, as in [The Docker API](https://github.com/moby/moby/blob/master/api/types/client.go#L73) |  |
| description | string |  |  |
| [meta](#meta) | object |  |  |
| output_dir | string |  |  |
| [scrub](#scrub) | object |  |  |
| timeout_seconds | integer |  |  |

*****

## container_logs_options

#### Type: object

__Args as would be given to `docker logs`, as in [The Docker API](https://github.com/moby/moby/blob/master/api/types/client.go#L73)__

Required: ShowStdout, ShowStderr, Since, Until, Timestamps, Follow, Tail, Details, ShowStdout, ShowStderr, Since, Until, Timestamps, Follow, Tail, Details, ShowStdout, ShowStderr, Since, Until, Timestamps, Follow, Tail, Details, ShowStdout, ShowStderr, Since, Until, Timestamps, Follow, Tail, Details

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

Required: 

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| All | boolean | Option to be used to filter the list of containers, as in [The Docker API](https://github.com/moby/moby/blob/master/api/types/client.go#L61) |  |
| Before | string | Option to be used to filter the list of containers, as in [The Docker API](https://github.com/moby/moby/blob/master/api/types/client.go#L61) |  |
| [Filters](#filters) | object | Option to be used to filter the list of containers, as in [The Docker API](https://github.com/moby/moby/blob/master/api/types/client.go#L61) |  |
| Latest | boolean | Option to be used to filter the list of containers, as in [The Docker API](https://github.com/moby/moby/blob/master/api/types/client.go#L61) |  |
| Limit | integer | Option to be used to filter the list of containers, as in [The Docker API](https://github.com/moby/moby/blob/master/api/types/client.go#L61) |  |
| Quiet | boolean | Option to be used to filter the list of containers, as in [The Docker API](https://github.com/moby/moby/blob/master/api/types/client.go#L61) |  |
| Since | string | Option to be used to filter the list of containers, as in [The Docker API](https://github.com/moby/moby/blob/master/api/types/client.go#L61) |  |
| Size | boolean | Option to be used to filter the list of containers, as in [The Docker API](https://github.com/moby/moby/blob/master/api/types/client.go#L61) |  |
| description | string |  |  |
| [meta](#meta) | object |  |  |
| output_dir | string |  |  |
| [scrub](#scrub) | object |  |  |
| timeout_seconds | integer |  |  |

*****

## docker.container-run

#### Type: object

__Collect the stderr/stdout of running a single docker container__

Required: output_dir, container_create_config

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| [container_create_config](#container_create_config) | object | Container create options as would be passed to `docker run` |  |
| description | string |  |  |
| enable_pull | boolean | If `true`, allow this container to be pulled if not present |  |
| [meta](#meta) | object |  |  |
| output_dir | string |  |  |
| [scrub](#scrub) | object |  |  |
| timeout_seconds | integer |  |  |

*****

## container_create_config

#### Type: object

__Container create options as would be passed to `docker run`Same as would be passed to `docker run`, as in [The Docker API](https://github.com/moby/moby/blob/master/api/types/configs.go#L13)__

Required: Config, Config

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

Required: Image, Image

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

## Entrypoint

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

## Test

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

## OnBuild

#### Type: string

__Container create options as would be passed to `docker run`__


*****

## Shell

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

Required: 

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
| ConsoleSize |  | Container create options as would be passed to `docker run` |  |
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

## Binds

#### Type: string

__Container create options as would be passed to `docker run`__


*****

## CapAdd

#### Type: string

__Container create options as would be passed to `docker run`__


*****

## CapDrop

#### Type: string

__Container create options as would be passed to `docker run`__


*****

## DeviceCgroupRules

#### Type: string

__Container create options as would be passed to `docker run`__


*****

## Devices

#### Type: object

__Container create options as would be passed to `docker run`__

Required: PathOnHost, PathInContainer, CgroupPermissions, PathOnHost, PathInContainer, CgroupPermissions

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| CgroupPermissions | string | Container create options as would be passed to `docker run` |  |
| PathInContainer | string | Container create options as would be passed to `docker run` |  |
| PathOnHost | string | Container create options as would be passed to `docker run` |  |

*****

## Dns

#### Type: string

__Container create options as would be passed to `docker run`__


*****

## DnsOptions

#### Type: string

__Container create options as would be passed to `docker run`__


*****

## DnsSearch

#### Type: string

__Container create options as would be passed to `docker run`__


*****

## ExtraHosts

#### Type: string

__Container create options as would be passed to `docker run`__


*****

## GroupAdd

#### Type: string

__Container create options as would be passed to `docker run`__


*****

## Links

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

## Mounts

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

## SecurityOpt

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

## VolumesFrom

#### Type: string

__Container create options as would be passed to `docker run`__


*****

## NetworkingConfig

#### Type: object

__Container create options as would be passed to `docker run`__

Required: EndpointsConfig

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

Required: output_dir, container, exec_config

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| container | string | The container name |  |
| description | string |  |  |
| [exec_config](#exec_config) | object | Config options as would be passed to `docker exec`, as in [The Docker API](https://github.com/moby/moby/blob/master/api/types/configs.go) |  |
| [meta](#meta) | object |  |  |
| output_dir | string |  |  |
| [scrub](#scrub) | object |  |  |
| timeout_seconds | integer |  |  |

*****

## docker.image-ls

#### Type: object

__Collect a list of docker images present on the server__

Required: 

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| All | boolean | Same as would be passed to `docker images` |  |
| [Filters](#filters) | object | Same as would be passed to `docker images` |  |
| description | string |  |  |
| [meta](#meta) | object |  |  |
| output_dir | string |  |  |
| [scrub](#scrub) | object |  |  |
| timeout_seconds | integer |  |  |

*****

## docker.images

#### Type: object

__Collect a list of docker images present on the server__

Required: 

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| All | boolean | Same as would be passed to `docker images` |  |
| [Filters](#filters) | object | Same as would be passed to `docker images` |  |
| description | string |  |  |
| [meta](#meta) | object |  |  |
| output_dir | string |  |  |
| [scrub](#scrub) | object |  |  |
| timeout_seconds | integer |  |  |

*****

## docker.info

#### Type: object

__Collect info about the Docker daemon__

Required: output_dir

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| description | string |  |  |
| [meta](#meta) | object |  |  |
| output_dir | string |  |  |
| [scrub](#scrub) | object |  |  |
| timeout_seconds | integer |  |  |

*****

## docker.logs

#### Type: object

__Collect the stdout/stderr logs from one or more docker containers. One of `container` or `container_list_options` is required.__

Required: output_dir

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| container | string | A container name |  |
| [container_list_options](#container_list_options) | object | Options to be used to filter the list of containers, as in [The Docker API](https://github.com/moby/moby/blob/master/api/types/client.go#L61) |  |
| [container_logs_options](#container_logs_options) | object | Args as would be given to `docker logs`, as in [The Docker API](https://github.com/moby/moby/blob/master/api/types/client.go#L73) |  |
| description | string |  |  |
| [meta](#meta) | object |  |  |
| output_dir | string |  |  |
| [scrub](#scrub) | object |  |  |
| timeout_seconds | integer |  |  |

*****

## docker.node-ls

#### Type: object

__Collect information about the nodes in a Docker Swarm installation__

Required: 

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| [Filters](#filters) | object | Same as would be passed to `docker node ls` |  |
| description | string |  |  |
| [meta](#meta) | object |  |  |
| output_dir | string |  |  |
| [scrub](#scrub) | object |  |  |
| timeout_seconds | integer |  |  |

*****

## docker.ps

#### Type: object

__Collect information about containers__

Required: 

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| All | boolean | Option to be used to filter the list of containers, as in [The Docker API](https://github.com/moby/moby/blob/master/api/types/client.go#L61) |  |
| Before | string | Option to be used to filter the list of containers, as in [The Docker API](https://github.com/moby/moby/blob/master/api/types/client.go#L61) |  |
| [Filters](#filters) | object | Option to be used to filter the list of containers, as in [The Docker API](https://github.com/moby/moby/blob/master/api/types/client.go#L61) |  |
| Latest | boolean | Option to be used to filter the list of containers, as in [The Docker API](https://github.com/moby/moby/blob/master/api/types/client.go#L61) |  |
| Limit | integer | Option to be used to filter the list of containers, as in [The Docker API](https://github.com/moby/moby/blob/master/api/types/client.go#L61) |  |
| Quiet | boolean | Option to be used to filter the list of containers, as in [The Docker API](https://github.com/moby/moby/blob/master/api/types/client.go#L61) |  |
| Since | string | Option to be used to filter the list of containers, as in [The Docker API](https://github.com/moby/moby/blob/master/api/types/client.go#L61) |  |
| Size | boolean | Option to be used to filter the list of containers, as in [The Docker API](https://github.com/moby/moby/blob/master/api/types/client.go#L61) |  |
| description | string |  |  |
| [meta](#meta) | object |  |  |
| output_dir | string |  |  |
| [scrub](#scrub) | object |  |  |
| timeout_seconds | integer |  |  |

*****

## docker.run

#### Type: object

__Collect the stderr/stdout of running a single docker container__

Required: output_dir, container_create_config

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| [container_create_config](#container_create_config) | object | Same as would be passed to `docker run`, as in [The Docker API](https://github.com/moby/moby/blob/master/api/types/configs.go#L13) |  |
| description | string |  |  |
| enable_pull | boolean | If `true`, allow this container to be pulled if not present |  |
| [meta](#meta) | object |  |  |
| output_dir | string |  |  |
| [scrub](#scrub) | object |  |  |
| timeout_seconds | integer |  |  |

*****

## docker.service-logs

#### Type: object

__Collect logs from a docker swarm service. One of `service` or `service_list_options` is required.__

Required: 

Required one of: service, service_list_options

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| [container_logs_options](#container_logs_options) | object | Args as would be given to `docker logs`, as in [The Docker API](https://github.com/moby/moby/blob/master/api/types/client.go#L73) |  |
| description | string |  |  |
| [meta](#meta) | object |  |  |
| output_dir | string |  |  |
| [scrub](#scrub) | object |  |  |
| service | string | A single docker swarm service for which to collect logs |  |
| [service_list_options](#service_list_options) | object | Options for filtering all swarm services |  |
| timeout_seconds | integer |  |  |

*****

## service_list_options

#### Type: object

__Options for filtering all swarm servicesOptions as would be passed to `docker stack servicesOptions as would be passed to `docker stack services__

Required: Filters, output_dir, Filters

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| [Filters](#filters) | object |  |  |
| description | string |  |  |
| [meta](#meta) | object |  |  |
| output_dir | string |  |  |
| [scrub](#scrub) | object |  |  |
| timeout_seconds | integer |  |  |

*****

## docker.service-ls

#### Type: object

__Collect a list of docker swarm services__

Required: 

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| [Filters](#filters) | object | Same as would be passed to `docker service ls` |  |
| description | string |  |  |
| [meta](#meta) | object |  |  |
| output_dir | string |  |  |
| [scrub](#scrub) | object |  |  |
| timeout_seconds | integer |  |  |

*****

## docker.service-ps

#### Type: object

__Collect information about the tasks run by one or more services__

Required: output_dir, output_dir, Filters

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| [Filters](#filters) | object | Same as would be passed to `docker service ps` |  |
| description | string |  |  |
| [meta](#meta) | object |  |  |
| output_dir | string |  |  |
| [scrub](#scrub) | object |  |  |
| timeout_seconds | integer |  |  |

*****

## docker.stack-service-logs

#### Type: object

__Collect logs from one or more services in a stack__

Required: output_dir, namespace

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| [container_logs_options](#container_logs_options) | object | Args as would be given to `docker logs`, as in [The Docker API](https://github.com/moby/moby/blob/master/api/types/client.go#L73) |  |
| description | string |  |  |
| [meta](#meta) | object |  |  |
| namespace | string | The stack's namespace |  |
| output_dir | string |  |  |
| [scrub](#scrub) | object |  |  |
| [service_list_options](#service_list_options) | object | Options as would be passed to `docker stack services |  |
| timeout_seconds | integer |  |  |

*****

## docker.stack-service-ls

#### Type: object

__Collect information about services in a swarm stack__

Required: output_dir, namespace

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| description | string |  |  |
| [meta](#meta) | object |  |  |
| namespace | string | The stack's namespace |  |
| output_dir | string |  |  |
| [scrub](#scrub) | object |  |  |
| [service_list_options](#service_list_options) | object | Options as would be passed to `docker stack services |  |
| timeout_seconds | integer |  |  |

*****

## docker.stack-service-ps

#### Type: object

__Collect information about the tasks running in a service__

Required: output_dir, namespace

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| description | string |  |  |
| [meta](#meta) | object |  |  |
| namespace | string | The stack's namespace |  |
| output_dir | string |  |  |
| [scrub](#scrub) | object |  |  |
| [task_list_options](#task_list_options) | object | Options as would be passed to `docker stack ps` |  |
| timeout_seconds | integer |  |  |

*****

## task_list_options

#### Type: object

__Options as would be passed to `docker stack ps`Options for filtering stack tasksOptions for filtering stack tasks__

Required: 

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| [Filters](#filters) | object |  |  |

*****

## docker.stack-task-logs

#### Type: object

__Collect logs from a docker swarm task__

Required: output_dir, namespace

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| [container_logs_options](#container_logs_options) | object | Args as would be given to `docker logs`, as in [The Docker API](https://github.com/moby/moby/blob/master/api/types/client.go#L73) |  |
| description | string |  |  |
| [meta](#meta) | object |  |  |
| namespace | string | The stack namespace |  |
| output_dir | string |  |  |
| [scrub](#scrub) | object |  |  |
| [task_list_options](#task_list_options) | object | Options for filtering stack tasks |  |
| timeout_seconds | integer |  |  |

*****

## docker.task-logs

#### Type: object

__Collect logs from a docker swarm task__

Required: output_dir

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| [container_logs_options](#container_logs_options) | object | Args as would be given to `docker logs`, as in [The Docker API](https://github.com/moby/moby/blob/master/api/types/client.go#L73) |  |
| description | string |  |  |
| id | string | The id of a single task |  |
| [meta](#meta) | object |  |  |
| output_dir | string |  |  |
| [scrub](#scrub) | object |  |  |
| [task_list_options](#task_list_options) | object | Options for filtering stack tasks |  |
| timeout_seconds | integer |  |  |

*****

## docker.task-ls

#### Type: object

__List docker swarm tasks__

Required: 

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| [Filters](#filters) | object | Filters for tasks |  |
| description | string |  |  |
| [meta](#meta) | object |  |  |
| output_dir | string |  |  |
| [scrub](#scrub) | object |  |  |
| timeout_seconds | integer |  |  |

*****

## docker.version

#### Type: object

__Get the version of the docker server__

Required: output_dir

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| description | string |  |  |
| [meta](#meta) | object |  |  |
| output_dir | string |  |  |
| [scrub](#scrub) | object |  |  |
| timeout_seconds | integer |  |  |

*****

## journald.logs

#### Type: object

__Collect logs from journald__

Required: unit, output_dir

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| description | string |  |  |
| [meta](#meta) | object |  |  |
| output_dir | string |  |  |
| reverse | boolean | Store logs in reverse chronological order |  |
| [scrub](#scrub) | object |  |  |
| since | string | Since date for log collection |  |
| timeout_seconds | integer |  |  |
| unit | string | Systemd unit from which to collect logs |  |

*****

## os.hostname

#### Type: object

__Get the hostname__

Required: output_dir

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| description | string |  |  |
| [meta](#meta) | object |  |  |
| output_dir | string |  |  |
| [scrub](#scrub) | object |  |  |
| timeout_seconds | integer |  |  |

*****

## os.http-request

#### Type: object

__Collect information about the response from making an HTTP request__

Required: url, method

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| body | string | The request body |  |
| description | string |  |  |
| [header](#header) | object | One or more headers to send |  |
| insecure | boolean | Set to `true` to skip TLS verification |  |
| [meta](#meta) | object |  |  |
| method | string | HTTP request method |  |
| output_dir | string |  |  |
| [scrub](#scrub) | object |  |  |
| timeout_seconds | integer |  |  |
| url | string | The HTTP request URL |  |

*****

## header

#### Type: object

__One or more headers to send__

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| [.*](#.*) | array |  |  |

*****

## os.loadavg

#### Type: object

__Collect information about CPU load__

Required: 

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| description | string |  |  |
| [meta](#meta) | object |  |  |
| output_dir | string |  |  |
| [scrub](#scrub) | object |  |  |
| template | string | (Optional) a template for the human-readable output |  |
| timeout_seconds | integer |  |  |

*****

## os.read-file

#### Type: object

__Read a file from the filesystem__

Required: output_dir, filepath

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| description | string |  |  |
| filepath | string | The file path on the host. If running Support Bundle via docker, this will work, but symlinks are not supported |  |
| [meta](#meta) | object |  |  |
| output_dir | string |  |  |
| [scrub](#scrub) | object |  |  |
| timeout_seconds | integer |  |  |

*****

## os.run-command

#### Type: object

__Run a command. Will be run in the support bundle container if `docker run replicated/support-bundle` is used. To run in a custom container, use `docker.run` instead__

Required: output_dir, name

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| args | array | The command arguments |  |
| description | string |  |  |
| dir | string | The working directory of the command |  |
| env | array | Specifies additional environment variables |  |
| [meta](#meta) | object |  |  |
| name | string | The command to run |  |
| output_dir | string |  |  |
| [scrub](#scrub) | object |  |  |
| timeout_seconds | integer |  |  |

*****

## args

#### Type: string

__add description to json file__


*****

## env

#### Type: string

__add description to json file__


*****

## os.uptime

#### Type: object

__Collect information about the system's uptime__

Required: output_dir

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| description | string |  |  |
| [meta](#meta) | object |  |  |
| output_dir | string |  |  |
| [scrub](#scrub) | object |  |  |
| template | string | Template for the human-readable output |  |
| timeout_seconds | integer |  |  |

*****

## retraced.events

#### Type: object

__Collect Audit Log events from a running Retraced instance__

Required: output_dir, api_token, api_endpoint, project_id

| Name    | Type    | Description | Example |
| ------- | ------- | ----------- | ------- |
| api_endpoint | string | The location of the Audit Logs API endpoint in the form `http(s)://<host>:<port>` |  |
| api_token | string | An API token for communicating with the audit log instance |  |
| description | string |  |  |
| insecure | boolean | Set to `true` to skip SSL verification on the endpoint |  |
| [mask](#mask) | object | A mask to determine what fields to include in the output file, as in https://github.com/retracedhq/retraced-go/blob/master/graphql.go#L160 |  |
| [meta](#meta) | object |  |  |
| output_dir | string |  |  |
| project_id | string | The Audit Log Project ID |  |
| [query](#query) | object | A structured query for filtering events, as in https://github.com/retracedhq/retraced-go/blob/master/graphql.go#L16 |  |
| [scrub](#scrub) | object |  |  |
| timeout_seconds | integer |  |  |

*****

## mask

#### Type: object

__A mask to determine what fields to include in the output file, as in https://github.com/retracedhq/retraced-go/blob/master/graphql.go#L160__

Required: 

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

__A structured query for filtering events, as in https://github.com/retracedhq/retraced-go/blob/master/graphql.go#L16__

Required: 

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
