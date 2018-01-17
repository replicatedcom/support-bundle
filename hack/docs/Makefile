
SCHEMA_JSON := ./schema.json
SCHEMA_MD := ./schema.md
MUTATIONS := ./mutations.json
HACK_DOCS := ./bin/hack-docs

deps:
	yarn

build:
	mkdir -p bin
	go build  -o ./bin/hack-docs ../../cmd/hack-docs

gen-schema:
	$(HACK_DOCS) > $(SCHEMA_JSON)

merge-mutations:
	`yarn bin`/ts-node src/index.ts merge-mutations -f $(SCHEMA_JSON) -m $(MUTATIONS)  > $(SCHEMA_JSON).tmp
	mv $(SCHEMA_JSON).tmp $(SCHEMA_JSON)

validate:
	`yarn bin`/ts-node src/index.ts validate -f $(SCHEMA_JSON)

gen-markdown:
	`yarn bin`/jsonschema-md $(SCHEMA_JSON) > $(SCHEMA_MD)

pipeline-strict: deps gen-schema merge-mutations validate gen-markdown
pipeline: deps gen-schema merge-mutations gen-markdown

