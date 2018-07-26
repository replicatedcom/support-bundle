import * as _ from "lodash";

export function genObjectFromPathAndExample(path: string, example: any): any {
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
