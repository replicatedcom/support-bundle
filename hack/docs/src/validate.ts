import * as fs from "fs";
import * as yaml from "js-yaml";
import * as util from "util";
import * as _ from "lodash";
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

export function shouldSkipKey(schemaKey: string) {
  return schemaKey === "output_dir" ||
    schemaKey === "description" ||
    schemaKey === "meta" ||
    schemaKey === "scrub" ||
    schemaKey === "timeout_seconds" ||
    schemaKey === "meta.customer";
}

export function validate(schemaType: any, path: string, maxDepth: number, schema: any) {
  const schemaKey: string = _.toPath(path).slice(-1)[0];
  console.log(`VALIDATING ${path}`);
  if (!schemaType.description) {
    if (!shouldSkipKey(schemaKey)) {
      throw new Error(`missing ${chalk.yellow("description")} at ${chalk.green(path)}; Children: ${chalk.green(`${Object.keys(schemaType.items || schemaType.properties || {})}`)}`);
    }
  }

  if (maxDepth === 1) {
    if (shouldSkipKey(schemaKey)) {
      return
    }
    if (schemaType.type !== "object") {
      return;
    }

    if (!schemaType.examples || !schemaType.examples.length) {
      throw new Error(`missing ${chalk.yellow("examples")} at ${chalk.green(path)}; Children: ${chalk.green(`${Object.keys(schemaType.items || schemaType.properties || {})}`)}`);
    }

    if (!schemaType._ext_outputs || !schemaType._ext_outputs.length) {
      if (!shouldSkipKey(schemaKey)) {
        throw new Error(`missing ${chalk.yellow("_ext_outputs")} at ${chalk.green(path)}; Children: ${chalk.green(`${Object.keys(schemaType.items || schemaType.properties || {})}`)}`);
      }
    }
    let i = 0;
    for (const example of schemaType.examples) {
      i += 1;
      let exampleToValidate = {
        specs: [
          {
            [schemaKey]: example,
          },
        ],
      };
      console.log(chalk.blue(yaml.safeDump(exampleToValidate)));
      const res = tv4.validateMultiple(exampleToValidate, schema, false, true);
      if (!res.valid) {
        console.log(util.inspect(exampleToValidate, false, 100, true));
        throw new Error(`invalid example ${example} at ${i} ${chalk.green(path)}; Error: at \n${chalk.red(`${res.errors.map((e) => "\t" + e.dataPath + " " + e.message).join("\n")}`)}`);
      }
    }

  }

  if (maxDepth === 0) {
    return;
  }

  if (schemaType.items) {
    validate(schemaType.items, path + ".items", maxDepth - 1, schema);
  }
  if (schemaType.properties) {
    for (const key of Object.keys(schemaType.properties)) {
      validate(schemaType.properties[key], path + ".properties[\"" + key + "\"]", maxDepth - 1, schema)
    }

  }
}

