import { expect } from "chai";
import { Signer } from "ethers";
import { ethers } from "hardhat";
import {
  ERC20MintableBurnable,
  ERC20MintableBurnable__factory,
} from "../typechain-types";

describe("ERC20MintableBurnable", function () {
  let erc20: ERC20MintableBurnable;
  let erc20Factory: ERC20MintableBurnable__factory;
  let owner: Signer;
  let addr1: Signer;

  beforeEach(async function () {
    erc20Factory = await ethers.getContractFactory("ERC20MintableBurnable");
    erc20 = await erc20Factory.deploy("Wrapped Kava", "WKAVA", 6n);
    [owner, addr1] = await ethers.getSigners();
  });

  describe("decimals", function () {
    it("should be the same as deployed", async function () {
      expect(await erc20.decimals()).to.be.equal(6);
    });
  });

  describe("mint", function () {
    it("should be callable by owner", async function () {
      const amount = 10n;

      const tx = erc20.mint(await addr1.getAddress(), amount);
      await expect(tx).to.not.be.reverted;

      const bal = await erc20.balanceOf(await addr1.getAddress())
      expect(bal).to.equal(amount);
    });

    it("should reject non-owner", async function () {
      const tx = erc20.connect(addr1).mint(await addr1.getAddress(), 10n);

      await expect(tx).to.be.revertedWith("Ownable: caller is not the owner");
    });
  });
});
