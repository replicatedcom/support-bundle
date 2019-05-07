import * as _ from "lodash";

export function genObjectFromPathAndExample(path: string, example: any): any {
  if (path.startsWith("definitions.")) {
    path = path.split(".").slice(2).join(".");
  }
  let obj = example;
  const parts = _.toPath(path);
  parts.reverse().forEach((part) => {
    switch (part) {
      case "items":
        obj = [obj];
        return;
      case "properties":
        return;
      default:
        obj = {[part]: obj};
        return;
    }
  });
  return obj;
}

export function genSchemaFromPath(path: string, schema: any): any {
  if (!path.startsWith("definitions.")) {
    return schema;
  }
  const definitionName = path.split(".")[1];
  return {
    $schema: "http://json-schema.org/draft-04/schema#",
    $ref: `#/definitions/${definitionName}`,
    definitions: schema.definitions,
  };
}
