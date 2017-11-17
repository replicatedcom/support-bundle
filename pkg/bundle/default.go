package bundle

import (
	"github.com/replicatedcom/support-bundle/pkg/spec"
	"github.com/replicatedcom/support-bundle/pkg/types"
)

var unconditionalYaml = `
specs:
  - builtin: supportbundle.version
    json: VERSION.json
  - builtin: core.read-command
    raw: /default/commands/date
    config:
      command: "date"
  - builtin: core.read-command
    raw: /default/commands/df
    config:
      command: "df"
      args:
      - "-al"
  - builtin: core.read-command
    raw: /default/commands/df_inodes
    config:
      command: "df"
      args:
      - "-ali"
  - builtin: core.read-command
    raw: /default/commands/dmesg
    config:
      command: "dmesg"
  - builtin: core.read-command
    raw: /default/commands/free
    config:
      command: "free"
      args:
      - "-m"
  - builtin: core.hostname
    raw: /default/commands/hostname
  - builtin: core.read-command
    raw: /default/commands/ip_addr_show
    config:
      command: "ip"
      args:
      - "-o"
      - "addr"
      - "show"
  - builtin: core.read-command
    raw: /default/commands/ip_link_show
    config:
      command: "ip"
      args:
      - "-o"
      - "link"
      - "show"
  - builtin: core.read-command
    raw: /default/commands/ip_route_show
    config:
      command: "ip"
      args:
      - "-o"
      - "route"
      - "show"
  - builtin: core.read-command
    raw: /default/commands/ps
    config:
      command: "ps"
      args:
      - "fauxwww"
  - builtin: core.uptime
    raw: /default/commands/uptime
  - builtin: docker.daemon
    json: /default/docker/
  - builtin: docker.read-host-file
    raw: /default/etc/default
    config:
      file_path: /etc/default/docker
  - builtin: docker.read-host-file
    raw: /default/etc
    config:
      file_path: /etc/hostname
  - builtin: docker.read-host-file
    raw: /default/etc
    config:
      file_path: /etc/hosts
  - builtin: docker.read-host-file
    raw: /default/etc
    config:
      file_path: /etc/os-release
  - builtin: docker.read-host-file
    raw: /default/etc/sysconfig
    config:
      file_path: /etc/sysconfig/docker
  - builtin: docker.read-host-file
    raw: /default/etc
    config:
      file_path: /etc/system-release
  - builtin: docker.read-host-file
    raw: /default/etc/systemd/system/docker.service.d
    config:
      file_path: /etc/systemd/system/docker.service.d/http-proxy.conf
  - builtin: docker.read-host-file
    raw: /default/etc
    config:
      file_path: /etc/timezone
  - builtin: core.read-file
    raw: /default/proc
    config:
      file_path: /proc/cpuinfo
  - builtin: core.read-file
    raw: /default/proc/meminfo
    config:
      file_path: /proc/meminfo
  - builtin: core.read-file
    raw: /default/proc/mounts
    config:
      file_path: /proc/mounts
  - builtin: core.read-file
    raw: /default/proc/uptime
    config:
      file_path: /proc/uptime
  - builtin: core.read-file
    raw: /default/proc/version
    config:
      file_path: /proc/version
  - builtin: core.read-file
    raw: /default/proc/vmstat
    config:
      file_path: /proc/vmstat
  - builtin: docker.read-host-file
    raw: /default/var/log/upstart
    config:
      file_path: /var/log/upstart/docker.log
`

var centosYaml = `
specs:
- builtin: docker.read-host-file
  raw: /default/etc/centos-release
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
