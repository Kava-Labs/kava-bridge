import { expect } from "chai";
import { ethers } from "hardhat";
import { Contract } from "ethers";

describe("Bridge", function () {
  let bridge: Contract;

  beforeEach(async function () {
    const Bridge = await ethers.getContractFactory("Bridge");
    bridge = await Bridge.deploy();
  });

  it("should reject eth transfers", async function () {
    await bridge.deployed();

    const contractAddress = bridge.address;
    const addr1 = (await ethers.getSigners())[1];

    const tx = addr1.sendTransaction({
      to: contractAddress,
      value: ethers.utils.parseEther("1.0"),
    });

    await expect(tx).to.be.revertedWith("no fallback nor receive function");
  });
});
