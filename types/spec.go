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
		FilePath    string   `yaml:"file_path"`
		Args        []string `yaml:"args"`
		ContainerID string   `yaml:"container_id"`
		Command     string   `yaml:"command"`
	}
}
