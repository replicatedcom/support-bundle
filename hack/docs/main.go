package main

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"

	"github.com/alecthomas/jsonschema"
	analyzeapi "github.com/replicatedcom/support-bundle/pkg/analyze/api"
	collectapi "github.com/replicatedcom/support-bundle/pkg/collect/types"
)

type collectDoc struct {
	Collect   collectapi.Collect         `json:"collect,omitempty"`
	Lifecycle []collectapi.LifecycleTask `json:"lifecycle,omitempty"`
}

type analyzeDoc struct {
	Analyze analyzeapi.Analyze `json:"analyze,omitempty"`
}

func main() {
	var s *jsonschema.Schema
	r := &jsonschema.Reflector{
		ExpandedStruct: true,
		PkgPathInName:  true,
		StripPrefix:    "github.com/replicatedcom/support-bundle",
	}
	if len(os.Args) > 1 && os.Args[1] == "analyze" {
		s = r.ReflectFromType(reflect.TypeOf(&analyzeDoc{}))
	} else {
		s = r.ReflectFromType(reflect.TypeOf(&collectDoc{}))
	}
	b, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
}
