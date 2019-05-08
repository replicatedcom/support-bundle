package main

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/alecthomas/jsonschema"
	analyzeapi "github.com/replicatedcom/support-bundle/pkg/analyze/api"
	collectapi "github.com/replicatedcom/support-bundle/pkg/collect/types"
)

type multiDoc struct {
	Collect   collectapi.Collect         `json:"collect,omitempty"`
	Analyze   analyzeapi.Analyze         `json:"analyze,omitempty"`
	Lifecycle []collectapi.LifecycleTask `json:"lifecycle,omitempty"`
}

func main() {
	r := &jsonschema.Reflector{
		ExpandedStruct: true,
		PkgPathInName:  true,
		StripPrefix:    "github.com/replicatedcom/support-bundle",
	}
	s := r.ReflectFromType(reflect.TypeOf(&multiDoc{}))
	b, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
}
