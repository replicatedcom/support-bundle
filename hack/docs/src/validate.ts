import * as fs from "fs";
import * as yaml from "js-yaml";
import * as util from "util";
import * as _ from "lodash";
import * as chalk from "chalk";
import * as process from "process";
import * as tv4 from "tv4";
import { genObjectFromPathAndExample, genSchemaFromPath } from "./common";

export const name = "validate";
export const describe = "Ensure every field has a description";
export const builder = {
  infile: {
    alias: "f",
    describe: "the schema file",
    default: "./schema.json",
  },
};
export const handler = (argv) => {
  process.stderr.write("merge-mutations called\n");
  const schema = JSON.parse(fs.readFileSync(argv.infile).toString());
  try {
    if (schema.properties.lifecycle) {
      validate(schema.properties.lifecycle, "properties.lifecycle", 100, schema, true, false);
      validate(schema.properties.collect, "properties.collect", 4, schema, true, false);
    }
    if (schema.properties.analyze) {
      validate(schema.properties.analyze, "properties.analyze", 7, schema, false, false);
      validate(schema.definitions.analyzeCondition, "definitions.analyzeCondition", 4, schema, false, false);
      validate(schema.definitions.analyzeInsight, "definitions.analyzeInsight", 100, schema, false, false);
    }
  } catch (err) {
    console.log(`\n\nFAILED ${err.message}`);
    process.exit(1);
  }
};

export function shouldSkipKey(schemaKey: string) {
  return schemaKey === "output_dir" ||
    schemaKey === "description" ||
    schemaKey === "name" ||
    schemaKey === "labels" ||
    schemaKey === "meta" ||
    schemaKey === "defer" ||
    schemaKey === "scrub" ||
    schemaKey === "timeout_seconds" ||
    schemaKey === "include_empty" ||
    schemaKey === "meta.customer" ||
    schemaKey === "meta.channel" ||
    schemaKey === "meta.watch";
}

export function shouldValidateExamples(path: string) {
  return path.startsWith("properties.collect.properties.v1.items.properties") ||
    path.startsWith("properties.lifecycle.items.properties") ||
    path.startsWith("properties.analyze.properties.v1.items.properties") ||
    path.startsWith("definitions.analyzeCondition.properties");
}

export function shouldValidateExtOutputs(path: string) {
  return path.startsWith("properties.collect.properties.v1.items.properties");
}

export function validate(schemaType: any, path: string, maxDepth: number, schema: any, hasOutput: boolean, parentHasExamples: boolean) {
  const schemaKey: string = _.toPath(path).slice(-1)[0];
  console.log(`VALIDATING ${path}`);

  if (shouldSkipKey(schemaKey)) {
    console.log(`\tSKIPPING ${path}`);
    return;
  }

  if (!schemaType.description) {
    throw new Error(`missing ${chalk.yellow("description")} at ${chalk.green(path)}; Children: ${chalk.green(`${Object.keys(schemaType.items || schemaType.properties || {})}`)}`);
  }

  if (schemaType.type !== "object" && schemaType.type !== "array") {
    return;
  }

  let hasExamples = parentHasExamples;
  if (!parentHasExamples && shouldValidateExamples(path)) {
    if (!schemaType.examples || !schemaType.examples.length) {
      throw new Error(`missing ${chalk.yellow("examples")} at ${chalk.green(path)}; Children: ${chalk.green(`${Object.keys(schemaType.items || schemaType.properties || {})}`)}`);
    }
    hasExamples = true;

    if (shouldValidateExtOutputs(path)) {
      if (hasOutput && (!schemaType._ext_outputs || !schemaType._ext_outputs.length)) {
        throw new Error(`missing ${chalk.yellow("_ext_outputs")} at ${chalk.green(path)}; Children: ${chalk.green(`${Object.keys(schemaType.items || schemaType.properties || {})}`)}`);
      }
    }

    let i = 0;
    for (const example of schemaType.examples) {
      i += 1;
      const exampleToValidate = genObjectFromPathAndExample(path, example);
      const schemaForValidation = genSchemaFromPath(path, schema);
      console.log(chalk.blue(yaml.safeDump(exampleToValidate)));
      const res = tv4.validateMultiple(exampleToValidate, schemaForValidation, false, true);
      if (!res.valid) {
        console.log(util.inspect(exampleToValidate, false, 100, true));
        throw new Error(`invalid example ${JSON.stringify(example)} at ${i} ${chalk.green(path)}; Error: at \n${chalk.red(`${res.errors.map((e) => "\t" + e.dataPath + " " + e.message).join("\n")}`)}`);
      }
    }
  }

  if (maxDepth === 0) {
    return;
  }

  if (schemaType.items) {
    validate(schemaType.items, path + ".items", maxDepth - 1, schema, hasOutput, hasExamples);
  }
  if (schemaType.properties) {
    for (const key of Object.keys(schemaType.properties)) {
      validate(schemaType.properties[key], path + ".properties[\"" + key + "\"]", maxDepth - 1, schema, hasOutput, hasExamples);
    }

  }
}
