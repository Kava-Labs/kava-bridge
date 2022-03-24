import { expect } from "chai";
import { tokens } from "./utils";

describe("utils", function () {
  it("tokens converts values to bigint 18 decimal representation", async function () {
    expect(tokens(1)).to.equal(1000000000000000000n);
    expect(tokens(1n)).to.equal(1000000000000000000n);
    expect(tokens(10)).to.equal(10000000000000000000n);
    expect(tokens(10n)).to.equal(10000000000000000000n);
    expect(tokens(102, 17)).to.equal(10200000000000000000n);
    expect(tokens(102n, 17n)).to.equal(10200000000000000000n);
  });
});
