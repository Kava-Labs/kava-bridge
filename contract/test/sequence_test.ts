import { expect } from "chai";
import { ethers } from "hardhat";
import { Contract, ContractFactory } from "ethers";

describe("Sequence", function () {
  let sequence: Contract;
  let sequenceFactory: ContractFactory;

  beforeEach(async function () {
    sequenceFactory = await ethers.getContractFactory("SequenceMock");
  });

  describe("increment", function () {
    it("should increment by 1 for each increment() call", async function () {
      sequence = await sequenceFactory.deploy(0);
      await sequence.deployed();

      let curSequence = 1n;

      for (let i = 0; i < 10; i++) {
        await sequence.increment();
        expect(await sequence.get()).to.eq(curSequence++);
      }
    });

    it("should wrap around to 0", async function () {
      const maxInt = 2n ** 256n - 1n;
      sequence = await sequenceFactory.deploy(maxInt);
      await sequence.deployed();

      // First check if the current value is max
      expect(await sequence.get()).to.eq(maxInt);
      await sequence.increment();
      expect(await sequence.get()).to.eq(0n);
    });
  });
});
