package types

type Spec struct {
	ID             string
	Builtin        string
	TimeoutSeconds int
	// paths
	Raw   string
	JSON  string
	Human string

	// Plan-specific config
	Config struct {
		FilePath      string   `yaml:"file_path"`
		Args          []string `yaml:"args"`
		ContainerID   string   `yaml:"container_id"`
		ContainerName string   `yaml:"container_name"`
		Command       string   `yaml:"command"`
		Scrub         Scrub    `yaml:"scrub"`
	}
}

type Scrub struct {
	Regex   string `yaml:"regex"`
	Replace string `yaml:"replace"`
}
