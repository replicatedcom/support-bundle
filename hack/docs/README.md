hack/docs
========

Hack to generate support bundle docs from the golang types.

Currently has several steps:

- generate jsonschema from go source, dump to `schema.json`
- modify generated jsonschema to add descriptions and examples from `mutations.json`
- (optional) validate the modified jsonschema to ensure all spec keys have examples and descriptions
- Turn the modified jsoncschema into schema.md, uses jsonschema-md to generate reference docs
- clean up the schema.md file with sed (because `jsonschema-md` is a crappy library and I might say we should get rid of it someday)
- generate schema2.md with a janky `src/markdown.ts`, this will probably become the real docs

The contents of schema.md and schema2.md can then be pasted into the help-center hugo site or wherever.

## Contributing

`mutations.json` is the only file that should be modified to extend the generated schema. The following files are committed, but autogenerated:

```
schema.json
schema.md
schema2.md
```


### Make commands

Two make commands for e2e, one with validation `make pipeline-strict` and one without `make pipeline`:




