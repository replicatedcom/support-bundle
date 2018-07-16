package main

import (
	"fmt"

	"github.com/mcuadros/go-jsonschema-generator"
	"github.com/replicatedcom/support-bundle/pkg/collect/types"
)

func main() {
	s := &jsonschema.Document{}
	s.Read(&types.Doc{})
	fmt.Println(s)

}
