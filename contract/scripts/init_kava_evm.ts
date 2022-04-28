import { ethers } from "hardhat";
import { Multicall, Multicall2, WKAVA } from "../typechain-types";

export async function main(): Promise<void> {
  // Hardhat always runs the compile task when running scripts with its command
  // line interface.
  //
  // If this script is run directly using `node` you may want to call compile
  // manually to make sure everything is compiled
  // await hre.run('compile');
  //
  const [signer] = await ethers.getSigners();

  console.log("Signing with account %s", signer.address);

  const multicall = await deployMulticall();
  console.log("Multicall deployed:\n\tAddress: %s", multicall.address);

  const multicall2 = await deployMulticall2();
  console.log("Multicall2 deployed:\n\tAddress: %s", multicall2.address);

  console.log("Completed contracts deployment");

  const wkava = await deployWKAVA();
  console.log("WKAVA deployed:\n\tAddress: %s", wkava.address);
}

export async function deployWKAVA(): Promise<WKAVA> {
  const wkavaFactory = await ethers.getContractFactory("WKAVA");
  const wkava = await wkavaFactory.deploy();

  await wkava.deployed();

  return wkava;
}

export async function deployMulticall(): Promise<Multicall> {
  const multicallFactory = await ethers.getContractFactory("Multicall");
  const multicall = await multicallFactory.deploy();

  await multicall.deployed();

  return multicall;
}

export async function deployMulticall2(): Promise<Multicall2> {
  const multicall2Factory = await ethers.getContractFactory("Multicall2");
  const multicall2 = await multicall2Factory.deploy();

  await multicall2.deployed();

  return multicall2;
}

// Only run main() when the script is run directly
if (require.main === module) {
  // We recommend this pattern to be able to use async/await everywhere
  // and properly handle errors.
  main().catch((error) => {
    console.error(error);
    process.exitCode = 1;
  });
}
