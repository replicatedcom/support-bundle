package types

type Doc struct {
	Specs []Spec
}

type Spec struct {
	ID             string
	Builtin        string
	TimeoutSeconds int
	// paths
	Raw   string
	JSON  string
	Human string

	// Plan-specific config
	Config Config
}

type Config struct {
	FilePath      string   `yaml:"file_path"`
	Args          []string `yaml:"args"`
	Image         string   `yaml:"image"`
	EnablePull    bool     `yaml:"enable_pull"`
	ContainerID   string   `yaml:"container_id"`
	ContainerName string   `yaml:"container_name"`
	Command       string   `yaml:"command"`
	Scrub         Scrub    `yaml:"scrub"`
	Env           []string `yaml:"env"`
	Binds         []string `yaml:"binds"`
}

type Scrub struct {
	Regex   string `yaml:"regex"`
	Replace string `yaml:"replace"`
}
