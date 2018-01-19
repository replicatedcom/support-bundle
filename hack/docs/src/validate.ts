import * as fs from "fs";
import * as chalk from "chalk";
import * as process from "process";
import * as tv4 from "tv4";

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
    validate(schema.properties.lifecycle, "properties.lifecycle", 100, schema);
    validate(schema.properties.specs, "properties.specs", 3, schema);
  } catch (err) {
    console.log(`\n\nFAILED ${err.message}`);
    process.exit(1);
  }
};

export function validate(schemaType: any, path: string, maxDepth: number, schema: any) {
  console.log(`VALIDATING ${path}`);
  if (!schemaType.description) {
    throw new Error(`missing description at ${chalk.green(path)}; Children: ${chalk.green(`${Object.keys(schemaType.items || schemaType.properties || {})}`)}`);
  }

  if (maxDepth === 1) {
    if (schemaType.type !== "object") {
      return;
    }

    if (!schemaType.examples || ! schemaType.examples.length) {
      throw new Error(`missing examples at ${chalk.green(path)}; Children: ${chalk.green(`${Object.keys(schemaType.items || schemaType.properties || {})}`)}`);
    }
    let i = 0;
    for (const example of schemaType.examples) {
      i += 1;
      const res = tv4.validateMultiple(example, schema, false, true);
      if (!res.valid) {
        throw new Error(`invalid example ${example} at ${i} ${chalk.green(path)}; Error: at ${chalk.red(`${res.errors.map((e) => e.dataPath + " " + e.message)}`)}`);
      }
    }

  }

  if (maxDepth === 0) {
    return;
  }

  if (schemaType.items) {
    validate(schemaType.items, path + ".items", maxDepth -1, schema);
  }
  if (schemaType.properties) {
    for (const key of Object.keys(schemaType.properties)) {
      validate(schemaType.properties[key], path + ".properties[\"" + key + "\"]", maxDepth -1, schema)
    }

  }
}

