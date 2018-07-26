import { describe, it } from "mocha";
import { expect } from "chai";
import { genObjectFromPathAndExample } from "./common";

describe("genObjectFromPathAndExample", () => {
  it("works", async () => {
    const obj = genObjectFromPathAndExample("properties.collect.properties.v1.items.properties[\"docker.container-cp\"]", {
        description: "the supergoodtool www site access logs",
    });

    expect(obj).to.deep.equal({
        collect: {
            v1: [
                {
                    "docker.container-cp": {
                        description: "the supergoodtool www site access logs",
                    },
                },
            ],
        },
    });
  });
});
