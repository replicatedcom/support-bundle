
SCHEMA_JSON := ./schema.json
SCHEMA_MD := ./schema.md
SCHEMA_MD_2 := ./schema2.md
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

clean-markdown:
	sed -i.bak 's/:undefined//g' $(SCHEMA_MD)
	sed -i.bak 's/undefined//g' $(SCHEMA_MD)
	rm schema.md.bak

gen-markdown-2:
	`yarn bin`/ts-node src/index.ts markdown -f $(SCHEMA_JSON) > $(SCHEMA_MD_2)


pipeline-strict: deps gen-schema merge-mutations validate gen-markdown clean-markdown gen-markdown-2
pipeline: deps gen-schema merge-mutations gen-markdown clean-markdown gen-markdown-2

