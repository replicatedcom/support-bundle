package main

import (
	"fmt"

	jsonschema "github.com/mcuadros/go-jsonschema-generator"
	collectapi "github.com/replicatedcom/support-bundle/pkg/collect/types"
)

func main() {
	s := &jsonschema.Document{}
	multiDoc := struct {
		Collect collectapi.Collect `json:"collect,omitempty"`
		// Analyze   analyzeapi.Analyze         `json:"analyze,omitempty"`
		Lifecycle []collectapi.LifecycleTask `json:"lifecycle,omitempty"`
	}{}
	s.Read(multiDoc)
	fmt.Println(s)
}
