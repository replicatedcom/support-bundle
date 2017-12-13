package version

var (
	build Build
)

type Build struct {
	Type         string
	Version      string
	GitSHA       string
	TimeFallback string `json:"time_fallback,omitempty"`
}

func Init() {
	build.Type = "external"
	build.Version = version
	build.GitSHA = gitSHA[:7]
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
