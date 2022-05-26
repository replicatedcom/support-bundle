.PHONY: build

SCHEMA_JSON := ./schema.json
SCHEMA_ANALYZE_JSON := ./schema-analyze.json
MUTATIONS := ./mutations.json
MUTATIONS_ANALYZE := ./mutations_analyze.json
HACK_DOCS := ./bin/hack-docs

HELP_CENTER_PATH ?= [path/to/help/center]
REPLICATED_LINT_PATH ?= [path/to/replicated-lint]

deps:
	yarn

test: deps
	yarn test

prebuild:
	rm -rf build
	mkdir -p build/collect
	mkdir -p build/analyze-variables
	mkdir -p build/analyze-conditions
	mkdir -p build/lifecycle

build: prebuild
	go build -o $(HACK_DOCS) .

gen-schema:
	$(HACK_DOCS) > $(SCHEMA_JSON)
	$(HACK_DOCS) analyze > $(SCHEMA_ANALYZE_JSON)

merge-mutations:
	`yarn bin`/ts-node src/index.ts merge-mutations -f $(SCHEMA_JSON) -m $(MUTATIONS)  > $(SCHEMA_JSON).tmp
	mv $(SCHEMA_JSON).tmp $(SCHEMA_JSON)
	`yarn bin`/ts-node src/index.ts merge-mutations -f $(SCHEMA_ANALYZE_JSON) -m $(MUTATIONS_ANALYZE)  > $(SCHEMA_ANALYZE_JSON).tmp
	mv $(SCHEMA_ANALYZE_JSON).tmp $(SCHEMA_ANALYZE_JSON)

validate:
	`yarn bin`/ts-node src/index.ts validate -f $(SCHEMA_JSON)
	`yarn bin`/ts-node src/index.ts validate -f $(SCHEMA_ANALYZE_JSON)

integration:
	`yarn bin`/ts-node src/index.ts integration -f $(SCHEMA_JSON)
	# `yarn bin`/ts-node src/index.ts integration -f $(SCHEMA_ANALYZE_JSON)

# this is some weird E2E testing for making sure the examples we put in the docs actually generate the files we say they do. Its jank and kinda WIP still. I don't intend for this to replace the ginkgo e2e testing
integration-bootstrap: integration-down
	:
	: Setting up docker stuff
	:
	docker run --rm -d --name=supergoodtool-www nginx
	docker pull debian
	:
	: doing some swarm things -- you may want to initialize with
	:
	:     docker swarm init --advertise-addr=172.17.0.1 --listen-addr=0.0.0.0
	:
	docker service create --detach=true --name cooltool-www nginx
	sleep 5
	docker stack deploy cooltool-core -c fixtures/fake-stack.yml
	:
	:
	: Adding systemd units, this will require sudo
	:
	:
	sudo cp fixtures/cooltool-api.service /etc/systemd/system/
	sudo systemctl daemon-reload
	sudo systemctl enable cooltool-api

integration-down:
	docker rm -f supergoodtool-www || :
	docker rmi debian || :
	docker service remove cooltool-www || :
	docker stack down cooltool-core || :
	:
	:
	: removing systemd units. this will require sudo
	:
	:
	sudo systemctl disable cooltool-api.service || :
	sudo rm -f /etc/systemd/system/cooltool-api.service || :
	sudo systemctl daemon-reload || :

gen-markdown-2:
	mkdir -p build
	`yarn bin`/ts-node src/index.ts markdown -f $(SCHEMA_JSON) -o ./build
	`yarn bin`/ts-node src/index.ts markdown -f $(SCHEMA_ANALYZE_JSON) -o ./build


pipeline-strict: test build gen-schema merge-mutations validate integration gen-markdown-2
pipeline-nointegration: test build gen-schema merge-mutations validate gen-markdown-2
pipeline: test build gen-schema merge-mutations gen-markdown-2

copy-to-help-center:
	: To Copy to replicated lint, run the following command
	@echo ""
	@echo "    make pipeline"
	@echo "    /bin/cp build/shared.md $(HELP_CENTER_PATH)/content/api/support-bundle-yaml-specs/"
	@echo "    /bin/cp build/collect/* $(HELP_CENTER_PATH)/content/api/support-bundle-yaml-specs/"
	@echo "    /bin/cp build/analyze-variables/* $(HELP_CENTER_PATH)/content/api/analyze-yaml-variable-specs/"
	@echo "    /bin/cp build/analyze-conditions/* $(HELP_CENTER_PATH)/content/api/analyze-yaml-condition-specs/"
	@echo "    /bin/cp build/lifecycle/* $(HELP_CENTER_PATH)/content/api/support-bundle-yaml-lifecycle/"

copy-replicated-lint:
	: To Copy to replicated lint, run the following command
	@echo ""
	@echo "    make pipeline"
	@echo "    /bin/cp schema.json  $(REPLICATED_LINT_PATH)/projects/replicated-supportbundle/schema.json"
	@echo "    /bin/cp schema-analyze.json  $(REPLICATED_LINT_PATH)/projects/replicated-analyze/schema.json"
