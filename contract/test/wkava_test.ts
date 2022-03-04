import { expect } from "chai";
import { Signer, ContractReceipt, ContractTransaction } from "ethers";
import { ethers } from "hardhat";
import { WKAVA, WKAVA__factory as WKAVAFactory } from "../typechain-types";

describe("WKAVA", function () {
  let wkava: WKAVA;
  let wkavaFactory: WKAVAFactory;
  let addr1: Signer;

  beforeEach(async function () {
    wkavaFactory = await ethers.getContractFactory("WKAVA");
    wkava = await wkavaFactory.deploy();
    addr1 = (await ethers.getSigners())[1];
  });

  describe("Initialization values", function () {
    it("should have correct name", async function () {
      expect(await wkava.name()).to.be.equal("Wrapped Kava");
    });

    it("should have correct symbol", async function () {
      expect(await wkava.symbol()).to.be.equal("WKAVA");
    });

    it("should have correct number of decimals", async function () {
      expect(await wkava.decimals()).to.be.equal(6);
    });
  });

  describe("deposit", function () {
    it("should allow users to deposit", async function () {
      const userBalanceBefore = await wkava.balanceOf(await addr1.getAddress());
      expect(userBalanceBefore).to.equal(0);

      const depositAmt = 5000;
      await wkava.connect(addr1).deposit({ value: depositAmt });

      const userBalanceAfter = await wkava.balanceOf(await addr1.getAddress());
      expect(userBalanceAfter).to.equal(depositAmt);
    });

    it("should accept ethereum sent directly to the contract", async function () {
      const userBalanceBefore = await wkava.balanceOf(await addr1.getAddress());
      expect(userBalanceBefore).to.equal(0);

      const depositAmt = ethers.utils.parseEther("1.0");
      await addr1.sendTransaction({
        to: wkava.address,
        value: depositAmt, // Sends exactly 1.0 ether
      });

      const userBalanceAfter = await wkava.balanceOf(await addr1.getAddress());
      expect(userBalanceAfter).to.equal(depositAmt);
    });

    it("should emit Deposit event", async function () {
      const depositAmt = 5000;
      const tx: ContractTransaction = await wkava
        .connect(addr1)
        .deposit({ value: depositAmt });
      const receipt: ContractReceipt = await tx.wait();

      const event = receipt.events?.find((x: any) => {
        return x.event === "Deposit";
      });

      expect(event?.args?.src).to.equal(await addr1.getAddress());
      expect(event?.args?.wad).to.equal(depositAmt);
    });
  });

  describe("withdraw", function () {
    const depositAmt = 5000;

    beforeEach(async function () {
      await wkava.connect(addr1).deposit({ value: depositAmt });
    });

    it("should allow users to withdraw", async function () {
      const userBalanceBefore = await wkava.balanceOf(await addr1.getAddress());
      expect(userBalanceBefore).to.equal(depositAmt);

      const withdrawAmt = depositAmt / 2;
      await wkava.connect(addr1).withdraw(withdrawAmt);

      const userBalanceAfter = await wkava.balanceOf(await addr1.getAddress());
      expect(userBalanceAfter).to.equal(depositAmt - withdrawAmt);

      await wkava.connect(addr1).withdraw(withdrawAmt);
      const userBalanceFinal = await wkava.balanceOf(await addr1.getAddress());
      expect(userBalanceFinal).to.equal(0);
    });

    it("should emit Withdrawal event", async function () {
      const tx: ContractTransaction = await wkava
        .connect(addr1)
        .withdraw(depositAmt);
      const receipt: ContractReceipt = await tx.wait();

      const event = receipt.events?.find((x: any) => {
        return x.event === "Withdrawal";
      });

      expect(event?.args?.src).to.equal(await addr1.getAddress());
      expect(event?.args?.wad).to.equal(depositAmt);
    });
  });
});
