import * as fs from "fs";
import * as process from "process";
import * as path from "path";
import * as yaml from "js-yaml";
import * as _ from "lodash";
import * as util from "util";
import * as tmp from "tmp";
import { spawnSync } from "child_process";
import chalk from "chalk";

export const name = "integration";
export const describe = "Test support bundle binary for each example";
export const builder = {
  infile: {
    alias: "f",
    describe: "the schema file",
    default: "./schema.json",
  },
  supportbundle: {
    alias: "s",
    describe: "path to support bundle binary",
    default: "/home/dex/go/src/github.com/replicatedcom/support-bundle/bin/amd64/support-bundle",
  },
  bootstrap: {
    alias: "b",
    describe: "Don't run the tests, just bootstrap the necessary fixture containers",
    default: false,
  },
};

export const handler = (argv) => {
  if (argv.bootstrap) {
    hackhackbootstrap();
    return
  }

  const schema = JSON.parse(fs.readFileSync(argv.infile).toString());
  const specTypes = schema.properties.specs.items.properties;
  const failures = [] as any[];
  const f = () => {
    for (const specType of Object.keys(specTypes)) {
      if (specType.indexOf(".") === -1) {
        console.log(`SKIPPING ${specType}`);
        continue
      }

      console.log(chalk.blue(`TEST ${specType}`));

      if (!specTypes[specType].examples) {
        continue;
      }

      const outputs = specTypes[specType]._ext_outputs || [];
      for (const example of specTypes[specType].examples) {
        let expectedOutputs = outputs.map(o => path.join(example.output_dir, o.path));
        for (const file of  ['/VERSION.json', '/VERSION.human', '/VERSION.raw']) {
          expectedOutputs.push(file);
        }
        let tmpdir;
        try {
          tmpdir = tmp.dirSync({unsafeCleanup: true});
          validateExample(tmpdir.name, specType, example, argv.supportbundle, expectedOutputs);

        } catch (err) {
          failures.push(err);
          return;
        } finally {
          if (tmpdir) {
            tmpdir.removeCallback();
          }
        }
      }

    }
  };
  f();
  if (failures && failures.length) {
    console.log(chalk.red("FAILED"));
    for (const failure of failures) {
      console.log(util.inspect(failure));
    }
    process.exit(1);
  }
};

function compareExpectedOutputsAgainstIndexJSON(cwd, expectedOutputs: string[]) {
  let indexjson = fs.readFileSync(`${cwd}/index.json`);
  const index = JSON.parse(indexjson);
  const indexFiles = index.map(f => f.path);

  // if (indexFiles.length === 3) {
  //   throw {indexFiles, "err": "only found three, expected more"}
  // }

  const diff = _.difference(indexFiles, expectedOutputs);

  if (!_.isEmpty(diff)) {
    throw {indexFiles, outputs: expectedOutputs, compare: diff}
  }
}

function ensureNoErrorsInErrorsJSON(cwd) {
  let errorjson = fs.readFileSync(`${cwd}/error.json`);
  const error = JSON.parse(errorjson);
  if (error) {
    throw {error}
  }
}

function untarBundleIn(cwd) {
  const {status: tstatus, stdout: tstdout, stderr: tstderr} = spawnSync("tar", [
    "xvf", "supportbundle.tar.gz",
  ], {cwd});

  if (tstatus !== 0) {
    console.log(chalk.red('X tar'));
    throw {
      "step": "tar",
      tstatus,
      stdout: tstdout.toString(),
      stderr: tstderr.toString(),
    };
  }
}

function generateBundleFromBuiltGoBinary(supportbundle, exampleSpec: string, cwd) {
  const {status, stdout, stderr} = spawnSync(supportbundle, [
    "generate", "--spec", exampleSpec, "--require-journald",
  ], {cwd});
  // console.log({stdout, stderr});

  if (status !== 0) {
    console.log(chalk.red('X generate'));
    throw {
      "step": "generate",
      status,
      stdout: stdout.toString(),
      stderr: stderr.toString(),
    };
  }
}

function validateExample(cwd, specType, example, supportbundleBinaryPath, expectedOutputs: string[]) {

  const exampleSpec = JSON.stringify({specs: [{[specType]: example}]});
  console.log(chalk.blue(yaml.safeDump({specs: [{[specType]: example}]})));

  generateBundleFromBuiltGoBinary(supportbundleBinaryPath, exampleSpec, cwd);
  untarBundleIn(cwd);
  ensureNoErrorsInErrorsJSON(cwd);

  let hasExpectedOutputs = expectedOutputs &&
    expectedOutputs.length &&
    expectedOutputs[0].indexOf("*") === -1 &&
    expectedOutputs[0].indexOf("{{") === -1;

  if (hasExpectedOutputs) {
    compareExpectedOutputsAgainstIndexJSON(cwd, expectedOutputs);
  }
}

function hackhackbootstrap() {
  console.log(`nothing to do`);
}
