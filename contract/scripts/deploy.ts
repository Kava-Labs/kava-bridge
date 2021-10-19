// We require the Hardhat Runtime Environment explicitly here. This is optional
// but useful for running the script in a standalone fashion through `node <script>`.
//
// When running the script with `npx hardhat run <script>` you'll find the Hardhat
// Runtime Environment's members available in the global scope.
import { ethers } from "hardhat";
import { Contract } from 'ethers';

export async function main(): Promise<string> {
  // Hardhat always runs the compile task when running scripts with its command
  // line interface.
  //
  // If this script is run directly using `node` you may want to call compile
  // manually to make sure everything is compiled
  // await hre.run('compile');
  //
  const addr = process.env.KAVA_BRIDGE_RELAYER_ADDRESS;

  if (!addr) {
    throw new Error("relayer address not set");
  }

  const relayer = ethers.utils.getAddress(addr);
  const bridge = await deployBridge(relayer);

  return bridge.address;
}

export async function deployBridge(relayer: string): Promise<Contract> {
  const Bridge = await ethers.getContractFactory("Bridge");
  const bridge = await Bridge.deploy(relayer);

  await bridge.deployed();

  return bridge;
}

// Only run main() when the script is run directly
if (require.main === module) {
// We recommend this pattern to be able to use async/await everywhere
// and properly handle errors.
  main().then((bridgeAddress: string) => {
    console.log(bridgeAddress);
  }).catch((error) => {
    console.error(error);
    process.exitCode = 1;
  });
}
