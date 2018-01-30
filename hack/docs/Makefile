.PHONY: build

SCHEMA_JSON := ./schema.json
SCHEMA_MD := ./schema.md
MUTATIONS := ./mutations.json
HACK_DOCS := ./bin/hack-docs

deps:
	yarn

prebuild:
	rm -rf build
	mkdir -p build

build:
	mkdir -p bin
	go build  -o $(HACK_DOCS) ../../cmd/hack-docs

gen-schema:
	$(HACK_DOCS) > $(SCHEMA_JSON)

merge-mutations:
	`yarn bin`/ts-node src/index.ts merge-mutations -f $(SCHEMA_JSON) -m $(MUTATIONS)  > $(SCHEMA_JSON).tmp
	mv $(SCHEMA_JSON).tmp $(SCHEMA_JSON)

validate:
	`yarn bin`/ts-node src/index.ts validate -f $(SCHEMA_JSON)

integration:
	`yarn bin`/ts-node src/index.ts integration -f $(SCHEMA_JSON)

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

gen-markdown:
	`yarn bin`/jsonschema-md $(SCHEMA_JSON) > $(SCHEMA_MD)

clean-markdown:
	sed -i.bak 's/:undefined//g' $(SCHEMA_MD)
	sed -i.bak 's/undefined//g' $(SCHEMA_MD)
	rm schema.md.bak

gen-markdown-2:
	mkdir -p build
	`yarn bin`/ts-node src/index.ts markdown -f $(SCHEMA_JSON) -o ./build


pipeline-strict: prebuild gen-schema merge-mutations validate integration gen-markdown clean-markdown gen-markdown-2
pipeline: prebuild gen-schema merge-mutations gen-markdown clean-markdown gen-markdown-2


copy-to-help-center:
	: To Copy to Help center, run the following command:
	:
	:    make pipeline
	:    cp build/* [path/to/help/center]/api/support-bundle-yaml-specs/

copy-replicated-lint:
	: To Copy to replicated lint, run the following command:
	:
	:    make pipeline
	:    cp schema.json  [path/to/replicated-lint]/projects/replicated-supportbundle/
