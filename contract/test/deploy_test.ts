import { expect } from "chai";
import { ethers } from "hardhat";
import { main, deployBridge } from "../scripts/deploy";
import { Bridge__factory } from '../typechain-types'

describe("deploy script", function () {
  describe("main", function () {
    it("sets the relayer address from the environment", async function () {
      const signers = await ethers.getSigners();
      const relayer = await signers[1].getAddress();

      process.env.KAVA_BRIDGE_RELAYER_ADDRESS = relayer;
      const bridgeAddress = await main();

      const Bridge = await ethers.getContractFactory("Bridge");
      const bridge = Bridge.attach(bridgeAddress);

      expect(await bridge.relayer()).to.eq(relayer);
    });
  });

  describe("deployBridge", function () {
    it("deploys the contract with correct relayer", async function () {
      const signers = await ethers.getSigners();
      const relayer = await signers[1].getAddress();
      const bridge = await deployBridge(relayer);

      expect(await bridge.relayer()).to.eq(relayer);
    });
  });
});
