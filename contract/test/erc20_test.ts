import { expect } from "chai";
import { Signer } from "ethers";
import { ethers } from "hardhat";
import {
  ERC20MintableBurnable,
  ERC20MintableBurnable__factory as ERC20MintableBurnableFactory,
} from "../typechain-types";
import { kavaAddrToBytes32, testKavaAddrs } from "./utils";

describe("ERC20MintableBurnable", function () {
  let erc20: ERC20MintableBurnable;
  let erc20Factory: ERC20MintableBurnableFactory;
  let owner: Signer;
  let sender: Signer;
  let ethAddr: Signer;

  beforeEach(async function () {
    erc20Factory = await ethers.getContractFactory("ERC20MintableBurnable");
    erc20 = await erc20Factory.deploy("Wrapped Kava", "WKAVA", 6n);
    [owner, sender, ethAddr] = await ethers.getSigners();
  });

  describe("decimals", function () {
    it("should be the same as deployed", async function () {
      expect(await erc20.decimals()).to.be.equal(6);
    });
  });

  describe("mint", function () {
    it("should be callable by owner", async function () {
      const amount = 10n;

      const tx = erc20.connect(owner).mint(await sender.getAddress(), amount);
      await expect(tx).to.not.be.reverted;

      const bal = await erc20.balanceOf(await sender.getAddress());
      expect(bal).to.equal(amount);
    });

    it("should reject non-owner", async function () {
      const tx = erc20.connect(sender).mint(await sender.getAddress(), 10n);

      await expect(tx).to.be.revertedWith("Ownable: caller is not the owner");
    });
  });

  describe("withdraw", function () {
    let amount: bigint;

    beforeEach(async function () {
      const tx = await erc20
        .connect(owner)
        .mint(await sender.getAddress(), 100n);

      await tx.wait();
      erc20 = erc20.connect(sender);
      amount = 10n;
    });

    it("should emit a Withdraw event with (sender, toAddr, amount)", async function () {
      const withdrawTx = erc20.withdraw(await ethAddr.getAddress(), amount);

      await expect(withdrawTx)
        .to.emit(erc20, "Withdraw")
        .withArgs(
          await sender.getAddress(),
          await ethAddr.getAddress(),
          amount
        );
    });

    it("should index sender, toAddr in the Withdraw event", async function () {
      const event = erc20.interface.events["Withdraw(address,address,uint256)"];

      const [tokenParam, toAddrParam, amountParam] = event.inputs;

      expect(tokenParam.name).to.equal("sender");
      expect(tokenParam.indexed).to.equal(true);

      expect(toAddrParam.name).to.equal("toAddr");
      expect(toAddrParam.indexed).to.equal(true);

      expect(amountParam.name).to.equal("amount");
      expect(amountParam.indexed).to.equal(false);
    });

    it("should burn the account token amount from contract", async function () {
      const toAddr = await ethAddr.getAddress();

      await expect(() => erc20.withdraw(toAddr, amount)).to.changeTokenBalance(
        erc20,
        sender,
        -1n * amount
      );
    });

    it("should fail when ERC20 withdraw amount exceeds balance", async function () {
      const withdrawTx = erc20.withdraw(
        await ethAddr.getAddress(),
        amount * 100n
      );
      await expect(withdrawTx).to.be.revertedWith(
        "ERC20: burn amount exceeds balance"
      );
    });
  });

  describe("withdrawFrom", function () {
    let amount: bigint;

    beforeEach(async function () {
      const tx = await erc20
        .connect(owner)
        .mint(await sender.getAddress(), 100n);

      await tx.wait();
      erc20 = erc20.connect(sender);
      amount = 10n;
    });

    it("should permit withdrawing allowance", async function () {
      const Withdrawer = await ethers.getContractFactory("Withdrawer");
      const withdrawer = await Withdrawer.deploy();
      await withdrawer.deployed();

      // Approve from sender
      await erc20.approve(withdrawer.address, amount);

      // Make tx from withdrawer contract
      const withdrawTx = withdrawer.withdrawFor(
        erc20.address,
        await sender.getAddress(),
        await ethAddr.getAddress(),
        amount
      );

      await expect(withdrawTx).to.not.be.reverted;
    });

    it("should not exceed allowance", async function () {
      const Withdrawer = await ethers.getContractFactory("Withdrawer");
      const withdrawer = await Withdrawer.deploy();
      await withdrawer.deployed();

      await erc20.approve(withdrawer.address, amount);

      const withdrawTx = erc20.withdrawFrom(
        await sender.getAddress(),
        await ethAddr.getAddress(),
        amount + 1n
      );
      await expect(withdrawTx).to.be.revertedWith(
        "ERC20: insufficient allowance"
      );
    });
  });

  describe("convertToCoin", function () {
    let amount: bigint;
    let toKavaAddr: string;

    beforeEach(async function () {
      const tx = await erc20
        .connect(owner)
        .mint(await sender.getAddress(), 100n);

      await tx.wait();
      erc20 = erc20.connect(sender);
      amount = 10n;

      toKavaAddr = ethers.utils.hexlify(kavaAddrToBytes32(testKavaAddrs[0]));
    });

    it("should emit a ConvertToCoin event with (sender, toKavaAddr, amount)", async function () {
      const withdrawTx = erc20.convertToCoin(toKavaAddr, amount);

      await expect(withdrawTx)
        .to.emit(erc20, "ConvertToCoin")
        .withArgs(await sender.getAddress(), toKavaAddr, amount);
    });

    it("should index sender, toAddr in the ConvertToCoin event", async function () {
      const event =
        erc20.interface.events["ConvertToCoin(address,bytes32,uint256)"];

      expect(event.inputs).to.be.length(3);

      const [tokenParam, toAddrParam, amountParam] = event.inputs;

      expect(tokenParam.name).to.equal("sender");
      expect(tokenParam.indexed).to.equal(true);

      expect(toAddrParam.name).to.equal("toKavaAddr");
      expect(toAddrParam.indexed).to.equal(true);

      expect(amountParam.name).to.equal("amount");
      expect(amountParam.indexed).to.equal(false);
    });

    it("should transfer amount to owner address", async function () {
      await expect(() =>
        erc20.convertToCoin(toKavaAddr, amount)
      ).to.changeTokenBalances(erc20, [sender, owner], [-1n * amount, amount]);
    });

    it("should fail when ERC20 withdraw amount exceeds balance", async function () {
      const withdrawTx = erc20.convertToCoin(toKavaAddr, amount * 100n);
      await expect(withdrawTx).to.be.revertedWith(
        "ERC20: transfer amount exceeds balance"
      );
    });
  });
});
