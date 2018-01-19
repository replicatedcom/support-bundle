import * as _ from "lodash";
import * as fs from "fs";
import * as yaml from "js-yaml";
import * as util from "util";
import * as process from "process";

export const name = "markdown";
export const describe = "Build markdown for examples in top-level elements";
export const builder = {
  infile: {
    alias: "f",
    describe: "the schema file",
    default: "./schema.json",
  },
};

export const handler = (argv) => {
  const schema = JSON.parse(fs.readFileSync(argv.infile).toString());


  const specs = schema.properties.specs.items.properties;
  for (const property of Object.keys(specs)) {
    if (property.indexOf(".") === -1) {
      continue
    }

    // todo unjank this with handlebars or something
    console.log(`
## ${property}

${specs[property].description || "TODO add description"}`);

    if (!specs[property].examples) {
      continue;
    }

    for (const example of specs[property].examples)
    console.log(`
${"```yaml"}
${yaml.safeDump(example)}
${"```"}
    `)
  }
};

