import { BigNumberish } from "ethers";
import { ethers } from "hardhat";
import {
  Bridge,
  ERC20MintableBurnable,
  WETH9,
  Multicall,
  Multicall2,
} from "../typechain-types";

const testUserAddress = "0x7Bbf300890857b8c241b219C6a489431669b3aFA";
const testRelayerAddress = "0xa2F728F997f62F47D4262a70947F6c36885dF9fa";

export async function main(): Promise<void> {
  // Hardhat always runs the compile task when running scripts with its command
  // line interface.
  //
  // If this script is run directly using `node` you may want to call compile
  // manually to make sure everything is compiled
  // await hre.run('compile');
  //
  const [signer] = await ethers.getSigners();
  const userAddr = ethers.utils.getAddress(testUserAddress);
  const relayerAddr = ethers.utils.getAddress(testRelayerAddress);

  console.log("Signing with account %s", signer.address);

  const bridge = await deployBridge(relayerAddr);
  console.log(
    "Bridge deployed:\n\tAddress: %s\n\tRelayer: %s",
    bridge.address,
    relayerAddr
  );

  const weth = await deployWETH();
  console.log("WETH deployed:\n\tAddress: %s", weth.address);

  const meowAmounts = new Map<string, BigNumberish>([
    [signer.address, 100_000_000_000_000_000_000n],
    [userAddr, 100_000_000_000_000_000_000n],
  ]);
  const erc20MEOW = await deployERC20WithAmounts(
    "Cat Token",
    "MEOW",
    18,
    meowAmounts
  );
  console.log(
    "ERC20 deployed:\n\tName:\t%s\n\tAddress: %s",
    await erc20MEOW.name(),
    erc20MEOW.address
  );

  const usdcAmounts = new Map<string, BigNumberish>([
    [signer.address, 100_000_000_000n],
    [userAddr, 100_000_000_000n],
  ]);
  const erc20USDC = await deployERC20WithAmounts(
    "USD Coin",
    "USDC",
    6,
    usdcAmounts
  );
  console.log(
    "ERC20 deployed:\n\tName:\t%s\n\tAddress: %s",
    await erc20USDC.name(),
    erc20USDC.address
  );

  const multicall = await deployMulticall();
  console.log("Multicall deployed:\n\tAddress: %s", multicall.address);

  const multicall2 = await deployMulticall2();
  console.log("Multicall2 deployed:\n\tAddress: %s", multicall2.address);

  console.log("Completed contracts deployment");

  const userFundTx = await signer.sendTransaction({
    to: userAddr,
    value: ethers.utils.parseEther("100.0"),
  });
  console.log("User funded in tx %s\n", userFundTx.hash);

  const relayerFundTx = await signer.sendTransaction({
    to: relayerAddr,
    value: ethers.utils.parseEther("100.0"),
  });
  console.log("Relayer funded in tx %s\n", relayerFundTx.hash);

  console.log("Completed funding accounts with eth");

  const depositTx = await weth.deposit({
    value: ethers.utils.parseEther("1000.0"),
  });
  await depositTx.wait();

  const transferTx = await weth.transfer(
    userAddr,
    ethers.utils.parseEther("100.0")
  );
  await transferTx.wait();

  console.log("Completed funding WETH for user");
}

export async function deployBridge(relayer: string): Promise<Bridge> {
  const bridgeFactory = await ethers.getContractFactory("Bridge");
  const bridge = await bridgeFactory.deploy(relayer);

  await bridge.deployed();

  return bridge;
}

export async function deployWETH(): Promise<WETH9> {
  const wethFactory = await ethers.getContractFactory("WETH9");
  const weth = await wethFactory.deploy();

  await weth.deployed();

  return weth;
}

export async function deployERC20WithAmounts(
  name: string,
  symbol: string,
  decimals: number,
  amounts: Map<string, BigNumberish>
): Promise<ERC20MintableBurnable> {
  const erc20Factory = await ethers.getContractFactory("ERC20MintableBurnable");
  const erc20 = await erc20Factory.deploy(name, symbol, decimals);

  await erc20.deployed();

  for (const [to, amount] of amounts) {
    const tx = await erc20.mint(to, amount);
    await tx.wait();
    console.log("minted %s %s to %s", symbol, amount, to);
  }

  return erc20;
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
