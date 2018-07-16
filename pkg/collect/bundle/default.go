package bundle

import (
	"github.com/replicatedcom/support-bundle/pkg/collect/types"
)

func DefaultSpecs() []types.Spec {
	return []types.Spec{
		{
			CoreRunCommand: &types.CoreRunCommandOptions{
				SpecShared: types.SpecShared{
					OutputDir: "default/commands/date",
				},
				Name: "date",
			},
		},
		{
			CoreRunCommand: &types.CoreRunCommandOptions{
				SpecShared: types.SpecShared{
					OutputDir: "default/commands/df",
				},
				Name: "df",
				Args: []string{"-al"},
			},
		},
		{
			CoreRunCommand: &types.CoreRunCommandOptions{
				SpecShared: types.SpecShared{
					OutputDir: "default/commands/df_inodes",
				},
				Name: "df",
				Args: []string{"-ali"},
			},
		},
		{
			CoreRunCommand: &types.CoreRunCommandOptions{
				SpecShared: types.SpecShared{
					OutputDir: "default/commands/dmesg",
				},
				Name: "dmesg",
			},
		},
		{
			CoreRunCommand: &types.CoreRunCommandOptions{
				SpecShared: types.SpecShared{
					OutputDir: "default/commands/free",
				},
				Name: "free",
				Args: []string{"-m"},
			},
		},
		{
			CoreHostname: &types.CoreHostnameOptions{
				SpecShared: types.SpecShared{
					OutputDir: "default/commands/hostname",
				},
			},
		},
		{
			CoreRunCommand: &types.CoreRunCommandOptions{
				SpecShared: types.SpecShared{
					OutputDir: "default/commands/ip_addr_show",
				},
				Name: "ip",
				Args: []string{"-o", "addr", "show"},
			},
		},
		{
			CoreRunCommand: &types.CoreRunCommandOptions{
				SpecShared: types.SpecShared{
					OutputDir: "default/commands/ip_link_show",
				},
				Name: "ip",
				Args: []string{"-o", "link", "show"},
			},
		},
		{
			CoreRunCommand: &types.CoreRunCommandOptions{
				SpecShared: types.SpecShared{
					OutputDir: "default/commands/ip_route_show",
				},
				Name: "ip",
				Args: []string{"-o", "route", "show"},
			},
		},
		{
			CoreRunCommand: &types.CoreRunCommandOptions{
				SpecShared: types.SpecShared{
					OutputDir: "default/commands/ps",
				},
				Name: "ps",
				Args: []string{"fauxwww"},
			},
		},
		{
			CoreUptime: &types.CoreUptimeOptions{
				SpecShared: types.SpecShared{
					OutputDir: "default/commands/uptime",
				},
			},
		},
		{
			DockerContainerLs: &types.DockerContainerLsOptions{
				SpecShared: types.SpecShared{
					OutputDir: "default/docker",
				},
			},
		},
		{
			DockerInfo: &types.DockerInfoOptions{
				SpecShared: types.SpecShared{
					OutputDir: "default/docker",
				},
			},
		},
		{
			DockerVersion: &types.DockerVersionOptions{
				SpecShared: types.SpecShared{
					OutputDir: "default/docker",
				},
			},
		},
		{
			DockerImageLs: &types.DockerImageLsOptions{
				SpecShared: types.SpecShared{
					OutputDir: "default/docker",
				},
			},
		},
		{
			CoreReadFile: &types.CoreReadFileOptions{
				SpecShared: types.SpecShared{
					OutputDir: "default/etc/default",
				},
				Filepath: "/etc/default/docker",
			},
		},
		{
			CoreReadFile: &types.CoreReadFileOptions{
				SpecShared: types.SpecShared{
					OutputDir: "default/etc",
				},
				Filepath: "/etc/fstab",
			},
		},
		{
			CoreReadFile: &types.CoreReadFileOptions{
				SpecShared: types.SpecShared{
					OutputDir: "default/etc",
				},
				Filepath: "/etc/hostname",
			},
		},
		{
			CoreReadFile: &types.CoreReadFileOptions{
				SpecShared: types.SpecShared{
					OutputDir: "default/etc",
				},
				Filepath: "/etc/hosts",
			},
		},
		{
			CoreReadFile: &types.CoreReadFileOptions{
				SpecShared: types.SpecShared{
					OutputDir: "default/etc",
				},
				Filepath: "/etc/os-release",
			},
		},
		{
			CoreReadFile: &types.CoreReadFileOptions{
				SpecShared: types.SpecShared{
					OutputDir: "default/etc/sysconfig",
				},
				Filepath: "/etc/sysconfig/docker",
			},
		},
		{
			CoreReadFile: &types.CoreReadFileOptions{
				SpecShared: types.SpecShared{
					OutputDir: "default/etc",
				},
				Filepath: "/etc/system-release",
			},
		},
		{
			CoreReadFile: &types.CoreReadFileOptions{
				SpecShared: types.SpecShared{
					OutputDir: "default/etc/systemd/system/docker.service.d",
				},
				Filepath: "/etc/systemd/system/docker.service.d/http-proxy.conf",
			},
		},
		{
			CoreReadFile: &types.CoreReadFileOptions{
				SpecShared: types.SpecShared{
					OutputDir: "default/etc",
				},
				Filepath: "/etc/timezone",
			},
		},
		{
			CoreLoadavg: &types.CoreLoadavgOptions{
				SpecShared: types.SpecShared{
					OutputDir: "default/proc",
				},
			},
		},
		{
			CoreReadFile: &types.CoreReadFileOptions{
				SpecShared: types.SpecShared{
					OutputDir: "default/proc",
				},
				Filepath: "/proc/cpuinfo",
			},
		},
		{
			CoreReadFile: &types.CoreReadFileOptions{
				SpecShared: types.SpecShared{
					OutputDir: "default/proc",
				},
				Filepath: "/proc/meminfo",
			},
		},
		{
			CoreReadFile: &types.CoreReadFileOptions{
				SpecShared: types.SpecShared{
					OutputDir: "default/proc",
				},
				Filepath: "/proc/mounts",
			},
		},
		{
			CoreReadFile: &types.CoreReadFileOptions{
				SpecShared: types.SpecShared{
					OutputDir: "default/proc",
				},
				Filepath: "/proc/uptime",
			},
		},
		{
			CoreReadFile: &types.CoreReadFileOptions{
				SpecShared: types.SpecShared{
					OutputDir: "default/proc",
				},
				Filepath: "/proc/version",
			},
		},
		{
			CoreReadFile: &types.CoreReadFileOptions{
				SpecShared: types.SpecShared{
					OutputDir: "default/proc",
				},
				Filepath: "/proc/vmstat",
			},
		},
		{
			CoreReadFile: &types.CoreReadFileOptions{
				SpecShared: types.SpecShared{
					OutputDir: "default/var/log/upstart",
				},
				Filepath: "/var/log/upstart/docker.log",
			},
		},
		{
			CoreReadFile: &types.CoreReadFileOptions{
				SpecShared: types.SpecShared{
					OutputDir: "default/etc",
				},
				Filepath: "/etc/centos-release",
			},
		},
		{
			JournaldLogs: &types.JournaldLogsOptions{
				SpecShared: types.SpecShared{
					OutputDir: "default/journald/docker",
				},
				Unit:  "docker",
				Since: "-14 days",
			},
		},
	}
}

func SupportBundleVersionSpec() types.Spec {
	return types.Spec{
		SpecShared: types.SpecShared{
			Description: "Support bundle version information",
			OutputDir:   "/",
		},
		SupportBundleVersion: &types.SupportBundleVersionOptions{},
	}
}

func CustomerJSONSpec(customerID string) types.Spec {
	return types.Spec{
		SpecShared: types.SpecShared{
			Description: "Support Bundle Customer Metadata",
			OutputDir:   "/",
		},
		CustomerMeta: &types.CustomerMetaOptions{
			CustomerID: customerID,
		},
	}
}
