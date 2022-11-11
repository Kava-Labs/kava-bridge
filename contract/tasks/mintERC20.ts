// We require the Hardhat Runtime Environment explicitly here. This is optional
// but useful for running the script in a standalone fashion through `node <script>`.
//
// When running the script with `npx hardhat run <script>` you'll find the Hardhat
// Runtime Environment's members available in the global scope.

import { task } from "hardhat/config";
import { HardhatRuntimeEnvironment } from "hardhat/types";

task(
  "mint-erc20",
  "Mints ERC20 tokens to an account. Signer must be owner of ERC20 contract."
)
  .addPositionalParam("contractAddress")
  .addPositionalParam("to")
  .addPositionalParam("amount")
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

  const erc20Factory = await hre.ethers.getContractFactory(
    "ERC20MintableBurnable"
  );

  const erc20 = erc20Factory.attach(args.contractAddress);

  const mintTx = await erc20.mint(args.to, args.amount);
  await mintTx.wait();

  console.log(`minted ${args.amount} tokens to ${args.to}`);
}
