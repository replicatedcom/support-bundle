
## docker.container-cp

Collect a file by copying from a running docker container

```yaml
specs:
  - description: the supergoodtool www site access logs
    docker.container-cp:
      container: supergoodtool-www
      src_path: /var/log/nginx/access.log
    output_dir: /www/access/

```
    

## docker.container-exec

Collect the stdout and stderr of `exec`-ing a command on a running docker container

```yaml
specs:
  - description: the supergoodtool debug output
    docker.container-exec:
      container: supergoodtool-www
      exec_config:
        Cmd:
          - toolctl
          - info
          - '--verbose'
    output_dir: /www/debug/

```
    

## docker.container-inspect

Collect the `docker inspect` output for one or more running or stopped containers

## docker.container-logs

Collect the logs from one or more docker containers

## docker.container-ls

Collect information about one or more containers

## docker.container-run

Collect the stderr/stdout of running a single docker container

## docker.exec

Collect the stdout/stderr of executing a command in an already running docker container

## docker.image-ls

Collect a list of docker images present on the server

## docker.images

Collect a list of docker images present on the server

## docker.info

Collect info about the Docker daemon

## docker.logs

Collect info about the Docker daemon

## docker.node-ls

Collect information about the nodes in a Docker Swarm installation

## docker.ps

Collect information about containers

## docker.run

Collect the stderr/stdout of running a single docker container

## docker.service-logs

Collect logs from a docker swarm service

## docker.service-ls

Collect a list of docker swarm services

## docker.service-ps

Collect information about the tasks run by one or more services

## docker.stack-service-logs

Collect logs from one or more services in a stack

## docker.stack-service-ls

Collect information about services in a swarm stack

## docker.stack-service-ps

Collect information about the tasks running in a service

## docker.stack-task-logs

Collect logs from a docker swarm task

## docker.task-logs

Collect logs from a docker swarm task

## docker.task-ls

TODO add description

## docker.version

TODO add description

## journald.logs

TODO add description

## kubernetes.api-versions

TODO add description

## kubernetes.cluster-info

TODO add description

## kubernetes.logs

TODO add description

## kubernetes.resource-list

TODO add description

## kubernetes.version

TODO add description

## meta.customer

TODO add description

## os.hostname

TODO add description

## os.http-request

TODO add description

## os.loadavg

TODO add description

## os.read-file

TODO add description

## os.run-command

TODO add description

## os.uptime

TODO add description

## retraced.events

TODO add description
