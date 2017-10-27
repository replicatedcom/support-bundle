package bundle

import (
	"github.com/replicatedcom/support-bundle/pkg/spec"
	"github.com/replicatedcom/support-bundle/pkg/types"
)

var unconditionalYaml = `
specs:
  - id: Audit log events TODO Internal to Replicated
  - builtin: core.read-command
    raw: /daemon/commands/date
    config:
      command: "date"
  - builtin: core.read-command
    raw: /daemon/commands/df
    config:
      command: "df"
      args:
      - "-al"
  - builtin: core.read-command
    raw: /daemon/commands/df_inodes
    config:
      command: "df"
      args:
      - "-ali"
  - builtin: core.read-command
    raw: /daemon/commands/dmesg
    config:
      command: "dmesg"
  - builtin: core.read-command
    raw: /daemon/commands/free
    config:
      command: "free"
      args:
      - "-m"
  - builtin: core.hostname
    raw: /daemon/commands/hostname
  - builtin: core.read-command
    raw: /daemon/commands/ip_addr_show
    config:
      command: "ip"
      args:
      - "-o"
      - "addr"
      - "show"
  - builtin: core.read-command
    raw: /daemon/commands/ip_link_show
    config:
      command: "ip"
      args:
      - "-o"
      - "link"
      - "show"
  - builtin: core.read-command
    raw: /daemon/commands/ip_route_show
    config:
      command: "ip"
      args:
      - "-o"
      - "route"
      - "show"
  - builtin: core.read-command
    raw: /daemon/commands/ps
    config:
      command: "ps"
      args:
      - "fauxwww"
  - builtin: core.uptime
    raw: /daemon/commands/uptime
  - builtin: docker.daemon
    json: /daemon/docker/
  - builtin: core.read-file
    raw: /daemon/etc/default/docker
    config:
      file_path: /etc/default/docker
  - builtin: core.read-file
    raw: /daemon/etc/default/replicated
    config:
      file_path: /etc/default/replicated
  - builtin: core.read-file
    raw: /daemon/etc/default/replicated-operator
    config:
      file_path: /etc/default/replicated-operator
  - builtin: core.read-file
    raw: /daemon/etc/hostname
    config:
      file_path: /etc/hostname
  - builtin: core.read-file
    raw: /daemon/etc/hosts
    config:
      file_path: /etc/hosts
  - builtin: core.read-file
    raw: /daemon/etc/os-release
    config:
      file_path: /etc/os-release
  - builtin: core.read-file
    raw: /daemon/etc/replicated.conf
    config:
      file_path: /etc/replicated.conf
  - builtin: core.read-file
    raw: /daemon/etc/sysconfig/docker
    config:
      file_path: /etc/sysconfig/docker
  - builtin: core.read-file
    raw: /daemon/etc/sysconfig/replicated
    config:
      file_path: /etc/sysconfig/replicated
  - builtin: core.read-file
    raw: /daemon/etc/sysconfig/replicated-operator
    config:
      file_path: /etc/sysconfig/replicated-operator
  - builtin: core.read-file
    raw: /daemon/etc/system-release
    config:
      file_path: /etc/system-release
  - builtin: core.read-file
    raw: /daemon/etc/systemd/system/docker.service.d/http-proxy.conf
    config:
      file_path: /etc/systemd/system/docker.service.d/http-proxy.conf
  - builtin: core.read-file
    raw: /daemon/etc/timezone
    config:
      file_path: /etc/timezone
  - builtin: core.read-command
    raw: /daemon/journald/replicated.log
    config:
      command: journalctl
      args:
      - "-u"
      - "replicated"
  - builtin: core.read-command
    raw: /daemon/journald/replicated-operator.log
    config:
      command: journalctl
      args:
      - "-u"
      - "replicated-operator"
  - builtin: core.read-command
    raw: /daemon/journald/replicated-ui.log
    config:
      command: journalctl
      args:
      - "-u"
      - "replicated-ui"
  - builtin: core.read-file
    raw: /daemon/proc/cpuinfo
    config:
      file_path: /proc/cpuinfo
  - builtin: core.read-file
    raw: /daemon/proc/meminfo
    config:
      file_path: /proc/meminfo
  - builtin: core.read-file
    raw: /daemon/proc/mounts
    config:
      file_path: /proc/mounts
  - builtin: core.read-file
    raw: /daemon/proc/uptime
    config:
      file_path: /proc/uptime
  - builtin: core.read-file
    raw: /daemon/proc/version
    config:
      file_path: /proc/version
  - builtin: core.read-file
    raw: /daemon/proc/vmstat
    config:
      file_path: /proc/vmstat
  - id: config-commands.txt TODO Internal to Replicated
  - id: daemon.json TODO Internal to Replicated
  - id: ledis-app.dump TODO Internal to Replicated
  - id: ledis-registry.dump TODO Internal to Replicated
  - id: ledis-snapshots.dump TODO Internal to Replicated
  - id: params.json TODO Internal to Replicated
  - builtin: docker.inspect
    json: /daemon/replicated/replicated-inspect.json
    config:
      container_name: replicated
  - builtin: docker.logs
    json: /daemon/replicated/replicated.log
    config:
      container_name: replicated
  - builtin: docker.inspect
    json: /daemon/replicated/replicated-operator-inspect.json
    config:
      container_name: replicated-operator
  - builtin: docker.logs
    id: "TODO limit to last 10000"
    json: /daemon/replicated/replicated-operator.log
    config:
      container_name: replicated-operator
  - builtin: docker.inspect
    json: /daemon/replicated/replicated-ui-inspect.json
    config:
      container_name: replicated-ui
  - builtin: docker.logs
    id: "TODO limit to last 10000"
    json: /daemon/replicated/replicated-ui.log
    config:
      container_name: replicated-ui
  - id: replicated-versions.txt TODO Internal to Replicated
  - id: goroutines.txt TODO Internal to Replicated
  - id: tasks.txt TODO Internal to Replicated
  - builtin: core.read-file
    raw: /daemon/var/log/upstart/replicated.log
    config:
      file_path: /var/log/upstart/replicated.log
  - builtin: core.read-file
    raw: /daemon/var/log/upstart/replicated-ui.log
    config:
      file_path: /var/log/upstart/replicated-ui.log
  - builtin: core.read-file
    raw: /daemon/var/log/upstart/replicated-operator.log
    config:
      file_path: /var/log/upstart/replicated-operator.log
  - builtin: core.read-file
    raw: /daemon/var/log/upstart/docker.log
    config:
      file_path: /var/log/upstart/docker.log
  - id: All <container_id> tasks TODO 
  - id: All <node_id> tasks TODO 
  - id: scheduler/nodes.txt TODO Internal to Replicated
`

var centosYaml = `
specs:
- builtin: core.read-file
  raw: /daemon/etc/centos-release
  config:
    file_path: /etc/centos-release
`

func DefaultSpecs() ([]types.Spec, error) {
	unconditionalSpecs, err := spec.Parse([]byte(unconditionalYaml))
	defaultSpecs := unconditionalSpecs
	if err != nil {
		return defaultSpecs, err
	}

	// if running on centos system
	if true {
		centosSpecs, err := spec.Parse([]byte(centosYaml))
		if err != nil {
			return defaultSpecs, err
		}

		defaultSpecs = append(defaultSpecs, centosSpecs...)
	}

	return defaultSpecs, err
}
