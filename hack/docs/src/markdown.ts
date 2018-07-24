import * as fs from "fs";
import * as yaml from "js-yaml";
import { SHARED_DOC, SHARED_SPEC_TEMPLATE } from "./template";

type SCHEMA_TYPE = "support-bundle-yaml-lifecycle" | "support-bundle-yaml-specs" | "analyze-yaml-specs";

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

function maybeRenderOutputs(specTypes: any, specType) {
  let doc = "";
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
    const description = specTypes[specType].properties[field].description;
    if (description) {
      console.log(`${field}: ${specType}.required: ${specTypes[specType].required}`);
      const { required: requiredArray = [] } = specTypes[specType];
      const isRequired = requiredArray.indexOf(field) !== -1;
      if (isRequired) {
        console.log(`\tREQUIRED ${field}`);
        required.push({field, description});
      } else {
        console.log(`\tOPTIONAL ${field}`);
        optional.push({field, description});
      }
    }
  }
  return {required, optional};
}

function maybeRenderExamples(specTypes: any, specType: string, schemaType: SCHEMA_TYPE) {
  let doc = "";
  if (specTypes[specType].examples) {
    for (const example of specTypes[specType].examples) {
      doc += `
${"```yaml"}
${yaml.safeDump(example)}${"```"}
`;
    }
  }
  return doc;
}

function writeHeader(specTypes: any, specType: string, schemaType: SCHEMA_TYPE) {
  return `---
categories:
- ${schemaType}
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

const renderEnd = (schemaType: SCHEMA_TYPE) => {
  if (schemaType === "support-bundle-yaml-specs") {
    return SHARED_SPEC_TEMPLATE;
  }
  return "";
};

const generateDoc: (outputLocation: string, specTypes: any, schemaType: SCHEMA_TYPE) => (specType: string) => void = (outputLocation, specTypes, schemaType) => (specType) => {
  console.log(`PROPERTY ${specType}`);
  // all specs follow the pattern "<plugin>.<spec_name>", if it doesnt follow the pattern its probably a shared property
  if ((schemaType === "support-bundle-yaml-specs" || schemaType === "analyze-yaml-specs") && specType.indexOf(".") === -1) {
    console.log(`SKIPPING ${specType}`);
    return;
  }
  const cleanProperty = specType.replace(/\./g, "-");

  let doc = "";
  doc += writeHeader(specTypes, specType, schemaType);
  doc += maybeRenderExamples(specTypes, specType, schemaType);

  const {required, optional} = parseParameters(specTypes, specType);
  doc += maybeRenderParameters(required, `Required`);
  doc += maybeRenderParameters(optional, `Optional`);
  doc += maybeRenderOutputs(specTypes, specType);
  doc += renderEnd(schemaType);

  fs.writeFileSync(`${outputLocation}/${cleanProperty}.md`, doc);
};

export const handler = (argv) => {
  const { infile, output } = argv;
  const schema = JSON.parse(fs.readFileSync(infile).toString());

  fs.writeFileSync(`${output}/shared.md`, SHARED_DOC);

  const specTypes = schema.properties.collect.properties.v1.items.properties;
  const analyzeTypes = schema.properties.analyze.properties.v1alpha1.items.properties;
  const lifecycleSpecTypes = schema.properties.lifecycle.items.properties;

  const generateDocSpecTypes = generateDoc(`${output}/collect`, specTypes, "support-bundle-yaml-specs");
  const generateDocAnalyzeTypes = generateDoc(`${output}/analyze`, analyzeTypes, "analyze-yaml-specs");
  const generateDocLifecycleSpecTypes = generateDoc(`${output}/lifecycle`, lifecycleSpecTypes, "support-bundle-yaml-lifecycle");

  Object.keys(specTypes).forEach(generateDocSpecTypes);
  Object.keys(analyzeTypes).forEach(generateDocAnalyzeTypes);
  Object.keys(lifecycleSpecTypes).forEach(generateDocLifecycleSpecTypes);
};
