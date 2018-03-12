import * as fs from "fs";
import * as yaml from "js-yaml";

const SHARED_DOC = `---
categories:
- support-bundle-yaml-specs
date: 2018-01-17T23:51:55Z
description: Reference Documentation for defining your Support Bundle collection and contents
index: docs
title: Support Bundle YAML Specs
weight: "1"
gradient: "purpleToPink"
---

## Support Bundle Collection Specs

Support Bundle collection specs can be used to define and customize what kinds diagnostic
information you want to collect to debug your application. All Support Bundle specs support the following shared parameters:

### Required Parameters

- ${"`"}output_dir${"`"} - The directory in the bundle to store the collection results

### Optional Parameters

- ${"`"}timeout_seconds${"`"} - An amount of time to allow a collection to run before abandoning it

- ${"`"}description${"`"} - A description of the file(s) being collected

- ${"`"}scrub${"`"} - A ${"`"}regex${"`"} and ${"`"}replace${"`"} specification for removing sensitive data from files in the bundle

- ${"`"}meta${"`"} - A ${"`"}name${"`"} and ${"`"}labels${"`"} that can be used to organize and identify support bundle elements in generated bundles

### Usage

An example is shown below for the ${"`"}os.read-file${"`"} collector.

${"```"}yaml
specs:
  - os.read-file:
      # path on the host
      filepath: /etc/goodtool.conf
      # path in the bundle
      output_dir: /files/etc/goodtool-conf
      # a description that will be included in the bundle
      description: The GoodTool application configuration file
      # give up if we can't read the file in 10 seconds
      timeout_seconds: 10
      # scrub anything that might be sensitive
      scrub:
        regex: (db_password|api_secret_key) = (.*)
        replace: $1 = REDACTED
      # metadata that will be included in the bundle
      meta:
        name: goodtool_conf
        labels:
          area: "configuration"
          type: "readfile"
${"```"}

`;


export const name = "markdown";
export const describe = "Build markdown for examples in top-level elements";
export const builder = {
  infile: {
    alias: "f",
    describe: "the schema file",
    default: "./schema.json",
  },
  output: {
    alias: "o",
    describe: "output dir",
    default: "./build",
  },
};

function maybeRenderOutputs(specTypes: any, specType) { let doc = "";
  const outputs = specTypes[specType]._ext_outputs;
  if (outputs) {
    doc += `
    
### Outputs

`;
    for (const output of outputs) {
      doc += `
- ${"`" + output.path + "`"} - ${output.description}
`;
    }

  }
  return doc;
}

function maybeRenderParameters(required: any[], typeOf) {
  let doc = "";
  if (required.length !== 0) {
    doc += `
    
### ${typeOf} Parameters

`;

    for (const fieldDescr of required) {
      doc += `
- ${"`" + fieldDescr.field + "`"} - ${fieldDescr.description}

`;
    }
    return doc;
  }
  return doc;
}

function parseParameters(specTypes: any, specType) {
  const required = [] as any[];
  const optional = [] as any[];
  for (const field of Object.keys(specTypes[specType].properties)) {
    let description = specTypes[specType].properties[field].description;
    if (description) {
      console.log(`${field}: ${specType}.required: ${specTypes[specType].required}`);
      let isRequired = specTypes[specType].required.indexOf(field) !== -1;
      if (isRequired) {
        console.log(`\tREQUIRED ${field}`);
        required.push({field, description});
      } else {
        console.log(`\tOPTIONAL ${field}`);
        optional.push({field, description})
      }
    }
  }
  return {required, optional};
}

function maybeRenderExamples(specTypes: any, specType) {
  let doc = "";
  if (specTypes[specType].examples) {
    for (const example of specTypes[specType].examples) {
      doc += `
${"```yaml"}
${yaml.safeDump({specs: [{[specType]: example}]})}${"```"}
`;
    }
  }
  return doc;
}

function writeHeader(specTypes: any, specType) {
  return `---
categories:
- support-bundle-yaml-specs
date: 2018-01-17T23:51:55Z
description: ${specTypes[specType].description || ""}
index: docs
title: ${specType}
weight: "100"
gradient: "purpleToPink"
---

## ${specType}

${specTypes[specType].description || ""}

`;
}

export const handler = (argv) => {
  const schema = JSON.parse(fs.readFileSync(argv.infile).toString());


  fs.writeFileSync(`${argv.output}/shared.md`, SHARED_DOC);

  const specTypes = schema.properties.specs.items.properties;
  for (const specType of Object.keys(specTypes)) {
    console.log(`PROPERTY ${specType}`);
    if (specType.indexOf(".") === -1) {
      console.log(`SKIPPING ${specType}`);
      continue
    }
    const cleanProperty = specType.replace(/\./g, "-");

    let doc = "";
    doc += writeHeader(specTypes, specType);
    doc += maybeRenderExamples(specTypes, specType);

    const {required, optional} = parseParameters(specTypes, specType);
    doc += maybeRenderParameters(required, `Required`);
    doc += maybeRenderParameters(optional, `Optional`);
    doc += maybeRenderOutputs(specTypes, specType);

    doc += `
    
<br>
{{< note title="Shared Parameters" >}}
This spec also inherits all of the required and optional [Shared Parameters](/api/support-bundle-yaml-specs/shared/)
{{< /note >}}
    
    `;



    fs.writeFileSync(`${argv.output}/${cleanProperty}.md`, doc);
  }
};

