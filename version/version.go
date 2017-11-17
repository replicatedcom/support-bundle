package version

import (
	"time"
)

var (
	build Build
)

type Build struct {
	Type      string
	Version   string
	GitSHA    string
	BuildTime time.Time
}

func Init() {
	build.Type = "external"
	build.Version = version
	build.GitSHA = gitSHA
	if buildTime != "" {
		build.BuildTime, _ = time.Parse(time.RFC3339, buildTime)
	} else if goReleaserTime != "" {
		build.BuildTime, _ = time.Parse("2006-01-02_15:04:05", buildTime)
	}
}

func GetBuild() Build {
	return build
}

func Version() string {
	return build.Version
}

func GitSHA() string {
	return build.GitSHA
}

func BuildTime() time.Time {
	return build.BuildTime
}
