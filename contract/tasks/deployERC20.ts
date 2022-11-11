// We require the Hardhat Runtime Environment explicitly here. This is optional
// but useful for running the script in a standalone fashion through `node <script>`.
//
// When running the script with `npx hardhat run <script>` you'll find the Hardhat
// Runtime Environment's members available in the global scope.

import { Contract } from "ethers";
import { task } from "hardhat/config";
import { HardhatRuntimeEnvironment } from "hardhat/types";

task("deploy-erc20", "Deploys an ERC20")
  .addPositionalParam("name")
  .addPositionalParam("symbol")
  .addPositionalParam("decimals")
  .setAction(async (taskArgs, hre) => {
    await main(taskArgs, hre);
  });

export async function main(
  args: { [key: string]: string },
  hre: HardhatRuntimeEnvironment
): Promise<void> {
  // Hardhat always runs the compile task when running scripts with its command
  // line interface.
  //
  // If this script is run directly using `node` you may want to call compile
  // manually to make sure everything is compiled
  // await hre.run('compile');
  //

  const [signer] = await hre.ethers.getSigners();

  console.log("Signing with account %s", signer.address);

  console.log("args", args);

  const erc20 = await deployErc20(
    hre,
    args.name,
    args.symbol,
    parseInt(args.decimals)
  );

  console.log("deployed erc20: ", erc20.address);
}

export async function deployErc20(
  hre: HardhatRuntimeEnvironment,
  name: string,
  symbol: string,
  decimals: number
): Promise<Contract> {
  const erc20Factory = await hre.ethers.getContractFactory(
    "ERC20MintableBurnable"
  );
  const erc20 = await erc20Factory.deploy(name, symbol, decimals);

  await erc20.deployed();

  return erc20;
}
