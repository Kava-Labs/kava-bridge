import { expect } from "chai";
import { kavaAddrToBytes32, bytes32ToKavaAddr, tokens } from "./utils";

describe("utils", function () {
  it("kavaAddrToBytes32 converts a bech32 address to zero padded bytes", async function () {
    const addrBytes = kavaAddrToBytes32(
      "kava1wn74shl496ktcfgqsc6yf0vvenhgq0hwuw6z2a"
    );

    const expectedAddrBytes = new Uint8Array([
      0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 116, 253, 88, 95, 245, 46, 172, 188,
      37, 0, 134, 52, 68, 189, 140, 204, 238, 128, 62, 238,
    ]);

    expect(addrBytes.length).to.equal(32);
    expect(addrBytes).to.deep.equal(expectedAddrBytes);
  });

  it("bytes32ToKavaAddr converts a zero padded address to bech32 string", async function () {
    const paddedAddrBytes = new Uint8Array([
      0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 116, 253, 88, 95, 245, 46, 172, 188,
      37, 0, 134, 52, 68, 189, 140, 204, 238, 128, 62, 238,
    ]);
    const addr = bytes32ToKavaAddr(paddedAddrBytes);
    expect(addr).to.equal("kava1wn74shl496ktcfgqsc6yf0vvenhgq0hwuw6z2a");

    const paddedAddrBytes2 = new Uint8Array([
      0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 253, 88, 95, 245, 46, 172, 188, 37,
      0, 134, 52, 68, 189, 140, 204, 238, 128, 62, 238,
    ]);
    const addr2 = bytes32ToKavaAddr(paddedAddrBytes2);
    expect(addr2).to.equal("kava1qr74shl496ktcfgqsc6yf0vvenhgq0hwakvsej");
  });

  it("tokens converts values to bigint 18 decimal representation", async function () {
    expect(tokens(1)).to.equal(1000000000000000000n);
    expect(tokens(1n)).to.equal(1000000000000000000n);
    expect(tokens(10)).to.equal(10000000000000000000n);
    expect(tokens(10n)).to.equal(10000000000000000000n);
    expect(tokens(102, 17)).to.equal(10200000000000000000n);
    expect(tokens(102n, 17n)).to.equal(10200000000000000000n);
  });
});
