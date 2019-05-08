import * as fs from "fs";
import * as yaml from "js-yaml";
import { SHARED_DOC, SHARED_SPEC_TEMPLATE } from "./template";
import { genObjectFromPathAndExample } from "./common";

type SCHEMA_TYPE = "support-bundle-yaml-lifecycle" | "support-bundle-yaml-specs" | "analyze-yaml-variable-specs"| "analyze-yaml-condition-specs";

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
  if (specTypes[specType].properties) {
    for (const field of Object.keys(specTypes[specType].properties)) {
      const description = specTypes[specType].properties[field].description;
      if (description) {
        console.log(`\t${field}: ${specType}.required: ${specTypes[specType].required}`);
        const { required: requiredArray = [] } = specTypes[specType];
        const isRequired = requiredArray.indexOf(field) !== -1;
        if (isRequired) {
          console.log(`\t\tREQUIRED ${field}`);
          required.push({field, description});
        } else {
          console.log(`\t\tOPTIONAL ${field}`);
          optional.push({field, description});
        }
      }
    }
  }
  return {required, optional};
}

function maybeRenderExamples(specTypes: any, specType: string, schemaType: SCHEMA_TYPE) {
  const schemaTypeMap: { [K in SCHEMA_TYPE]: string} = {
    "support-bundle-yaml-specs": "properties.collect.properties.v1.items",
    "analyze-yaml-variable-specs": "properties.analyze.properties.v1.items.properties.registerVariables.items",
    "analyze-yaml-condition-specs": "definitions.analyzeCondition",
    "support-bundle-yaml-lifecycle": "properties.lifecycle.items",
  };

  let doc = "";
  if (specTypes[specType].examples) {
    console.log(`\tEXAMPLES ${specType}`);
    for (const example of specTypes[specType].examples) {
      const path = schemaTypeMap[schemaType] + ".properties[\"" + specType + "\"]";
      const obj = genObjectFromPathAndExample(path, example);
      doc += `
${"```yaml"}
${yaml.safeDump(obj)}${"```"}
`;
    }
  }
  return doc;
}

function writeHeader(specTypes: any, specType: string, schemaType: SCHEMA_TYPE) {
  return `---
categories:
- ${schemaType}
date: 2019-05-07T12:00:00Z
description: ${specTypes[specType].description || ""}
index: docs
title: ${specType}
weight: "100"
gradient: "purpleToPink"
---

## ${specType}

**type ${specTypes[specType].type || "object"}**

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
  if (schemaType === "support-bundle-yaml-specs" && specType.indexOf(".") === -1) {
    console.log(`\tSKIPPING ${specType}`);
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
  const analyzeVariableTypes = schema.properties.analyze.properties.v1.items.properties.registerVariables.items.properties;
  const analyzeConditionTypes = schema.definitions.analyzeCondition.properties;
  const lifecycleSpecTypes = schema.properties.lifecycle.items.properties;

  const generateDocSpecTypes = generateDoc(`${output}/collect`, specTypes, "support-bundle-yaml-specs");
  const generateDocAnalyzeVariableTypes = generateDoc(`${output}/analyze-variables`, analyzeVariableTypes, "analyze-yaml-variable-specs");
  const generateDocAnalyzeConditionTypes = generateDoc(`${output}/analyze-conditions`, analyzeConditionTypes, "analyze-yaml-condition-specs");
  const generateDocLifecycleSpecTypes = generateDoc(`${output}/lifecycle`, lifecycleSpecTypes, "support-bundle-yaml-lifecycle");

  Object.keys(specTypes).forEach(generateDocSpecTypes);
  Object.keys(analyzeVariableTypes).forEach(generateDocAnalyzeVariableTypes);
  Object.keys(analyzeConditionTypes).forEach(generateDocAnalyzeConditionTypes);
  Object.keys(lifecycleSpecTypes).forEach(generateDocLifecycleSpecTypes);
};
