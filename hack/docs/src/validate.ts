import * as fs from "fs";
import * as chalk from "chalk";
import * as process from "process";

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
  validate(schema.properties.lifecycle, "properties.lifecycle");
  validate(schema.properties.specs, "properties.specs");
};

export function validate(schemaType: any, path: string) {
  if (!schemaType.description) {
    throw new Error(`missing description at ${chalk.green(path)}`);
  }

  if (schemaType.items) {
    validate(schemaType.items, path + ".items");
  }
  if (schemaType.properties) {
    for (const key of Object.keys(schemaType.properties)) {
      validate(schemaType.properties[key], path + ".properties[\"" + key + "\"]")
    }

  }
}

