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
	build.BuildTime, _ = time.Parse(time.RFC3339, buildTime)
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
