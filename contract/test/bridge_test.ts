import { expect } from "chai";
import { ethers } from "hardhat";
import { Contract, Signer } from "ethers";
import { kavaAddrToBytes32, tokens, testKavaAddrs } from "./utils";

describe("Bridge", function () {
  // the main bridge contract
  let bridge: Contract;
  // the deployer of the contracts
  let deployer: Signer;
  // relayer provided to the bridge upon deployment
  let relayer: Signer;
  // a sender / user of the bridge
  let sender: Signer;
  // a receiver / user of the bridge
  let receiver: Signer;

  beforeEach(async function () {
    const Bridge = await ethers.getContractFactory("Bridge");

    // returns 10 signers in order of owner, addr1, addr2...
    const signers = await ethers.getSigners();

    // assign commonly used addresses
    deployer = signers[0];
    relayer = signers[1];
    sender = signers[2];
    receiver = signers[3];

    // deploy bridge with relayer set, and sender connected
    bridge = await Bridge.deploy(await relayer.getAddress());
    // wait for bridge to be deployed
    await bridge.deployed();
    // connect sender
  });

  it("should not allow native eth transfers", async function () {
    const tx = sender.sendTransaction({
      to: bridge.address,
      value: tokens(1),
    });

    await expect(tx).to.be.reverted;
  });

  describe("#lock", function () {
    let token: Contract;
    let toAddr: string;
    let amount: bigint;

    beforeEach(async function () {
      // connect sender to the bridge
      bridge = await bridge.connect(sender);

      // assign valid attribute for #lock
      toAddr = ethers.utils.hexlify(kavaAddrToBytes32(testKavaAddrs[0]));
      amount = tokens(1);

      // deploy an ERC20 token
      const Token = await ethers.getContractFactory("ERC20Mock");
      token = await Token.deploy(
        "Token A",
        "TOKENA",
        await deployer.getAddress(),
        tokens(1000)
      );
      await token.deployed();

      // fund sender account with deployed token
      await token.transfer(await sender.getAddress(), 10n * amount);

      // allow bridge to transfer erc20 tokens on users behalf
      const tokenCon = await token.connect(sender);
      await tokenCon.approve(bridge.address, amount);
    });

    it("should not be payable", async function () {
      const lockTx = bridge.lock(token.address, toAddr, amount, {
        value: tokens(1),
      });
      await expect(lockTx).to.be.reverted;
    });

    it("should emit a Lock event with (token, sender, to, amount)", async function () {
      await token.approve(bridge.address, tokens(10));
      const lockTx = bridge.lock(token.address, toAddr, amount);

      await expect(lockTx)
        .to.emit(bridge, "Lock")
        .withArgs(token.address, await sender.getAddress(), toAddr, amount);
    });

    it("should transfer the token amount to the contract", async function () {
      const lockTx = bridge.lock(token.address, toAddr, amount);

      await expect(() => lockTx).to.changeTokenBalances(
        token,
        [sender, bridge],
        [-1n * amount, amount]
      );
    });

    it("should fail when ERC20 transferFrom reverts", async function () {
      const lockTx = bridge.lock(token.address, toAddr, 2n * amount);
      await expect(lockTx).to.be.revertedWith(
        "ERC20: transfer amount exceeds allowance"
      );
    });

    it("should not revert when ERC20 transferFrom returns no value", async function () {
      const Token = await ethers.getContractFactory("ERC20NoReturnMock");
      token = await Token.deploy();
      await token.deployed();

      const lockTx = bridge.lock(token.address, toAddr, amount);
      await expect(lockTx).to.not.be.reverted;
    });

    it("should not revert when ERC20 transferFrom returns true", async function () {
      const Token = await ethers.getContractFactory("ERC20ReturnTrueMock");
      token = await Token.deploy();
      await token.deployed();

      const lockTx = bridge.lock(token.address, toAddr, amount);
      await expect(lockTx).to.not.be.reverted;
    });

    it("should revert when ERC20 transferFrom returns false", async function () {
      const Token = await ethers.getContractFactory("ERC20ReturnFalseMock");
      token = await Token.deploy();
      await token.deployed();

      const lockTx = bridge.lock(token.address, toAddr, amount);
      await expect(lockTx).to.be.revertedWith(
        "SafeERC20: ERC20 operation did not succeed"
      );
    });
  });

  describe("unlock", function () {
    let token: Contract;
    let toAddr: string;
    let amount: bigint;

    beforeEach(async function () {
      // connect relayer to the bridge
      bridge = await bridge.connect(relayer);

      // assign valid attribute for #lock
      toAddr = await receiver.getAddress();
      amount = tokens(1);

      // deploy an ERC20 token
      const Token = await ethers.getContractFactory("ERC20Mock");
      token = await Token.deploy(
        "Token A",
        "TOKENA",
        await deployer.getAddress(),
        tokens(1000)
      );
      await token.deployed();

      // fund bridge with deployed token
      await token.transfer(bridge.address, 10n * amount);
    });

    it("should not be payable", async function () {
      const unlockTx = bridge.unlock(token.address, toAddr, amount, {
        value: tokens(1),
      });
      await expect(unlockTx).to.be.reverted;
    });

    it("should emit a Unlock event with (token, to, amount)", async function () {
      await token.approve(bridge.address, tokens(10));
      const unlockTx = bridge.unlock(token.address, toAddr, amount);

      await expect(unlockTx)
        .to.emit(bridge, "Unlock")
        .withArgs(token.address, toAddr, amount);
    });

    it("should transfer the token amount to the toAddr from contract", async function () {
      const unlockTx = bridge.unlock(token.address, toAddr, amount);

      await expect(() => unlockTx).to.changeTokenBalances(
        token,
        [bridge, receiver],
        [-1n * amount, amount]
      );
    });

    it("should fail when ERC20 transferFrom reverts", async function () {
      const unlockTx = bridge.unlock(token.address, toAddr, 100n * amount);
      await expect(unlockTx).to.be.revertedWith(
        "ERC20: transfer amount exceeds balance"
      );
    });

    it("should not revert when ERC20 transfer returns no value", async function () {
      const Token = await ethers.getContractFactory("ERC20NoReturnMock");
      token = await Token.deploy();
      await token.deployed();

      const unlockTx = bridge.unlock(token.address, toAddr, amount);
      await expect(unlockTx).to.not.be.reverted;
    });

    it("should not revert when ERC20 transfer returns true", async function () {
      const Token = await ethers.getContractFactory("ERC20ReturnTrueMock");
      token = await Token.deploy();
      await token.deployed();

      const unlockTx = bridge.unlock(token.address, toAddr, amount);
      await expect(unlockTx).to.not.be.reverted;
    });

    it("should revert when ERC20 transfer returns false", async function () {
      const Token = await ethers.getContractFactory("ERC20ReturnFalseMock");
      token = await Token.deploy();
      await token.deployed();

      const unlockTx = bridge.unlock(token.address, toAddr, amount);
      await expect(unlockTx).to.be.revertedWith(
        "SafeERC20: ERC20 operation did not succeed"
      );
    });

    it("should not be callable from an untrusted address", async function () {
      bridge = await bridge.connect(sender);
      const unlockTx = bridge.unlock(token.address, toAddr, amount);
      await expect(unlockTx).to.be.revertedWith(
        "Bridge: untrusted address"
      );
    });

    it("should not be callable from a re-entrant ERC20 contract", async function () {
      const Token = await ethers.getContractFactory("ERC20EvilUnlockMock");
      token = await Token.deploy();
      await token.deployed();

      const unlockTx = bridge.unlock(token.address, toAddr, amount);
      await expect(unlockTx).to.be.revertedWith(
        "Bridge: untrusted address"
      );
    });
  });
});
