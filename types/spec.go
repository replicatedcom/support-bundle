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
	Config interface{}
}

type ContainerConfig struct {
	ContainerID string `json="container_id"`
}
