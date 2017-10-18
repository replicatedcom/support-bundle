package main

import (
	"io/ioutil"
	"log"
	"time"

	"github.com/replicatedcom/support-bundle/bundle"
	"github.com/replicatedcom/support-bundle/plugins/core"
	"github.com/replicatedcom/support-bundle/plugins/docker"
	"github.com/replicatedcom/support-bundle/spec"
	"github.com/replicatedcom/support-bundle/types"
)

func main() {
	yml, err := ioutil.ReadFile("./spec.yml")
	if err != nil {
		log.Fatal(err)
	}
	specs, err := spec.Parse(yml)
	if err != nil {
		log.Fatal(err)
	}

	d, err := docker.New()
	if err != nil {
		log.Fatal(err)
	}
	planner := bundle.Planner{
		Plugins: map[string]types.Plugin{
			"core":   core.New(),
			"docker": d,
		},
	}
	tasks := planner.Plan(specs)
	if err := bundle.Generate(tasks, time.Minute, "/tmp/bundle.tar.gz"); err != nil {
		log.Fatal(err)
	}
}
