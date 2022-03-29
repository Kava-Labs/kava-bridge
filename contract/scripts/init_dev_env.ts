import { BigNumberish } from "ethers";
import { ethers } from "hardhat";
import { Bridge, ERC20MintableBurnable, WETH9 } from "../typechain-types";

export async function main(): Promise<void> {
  // Hardhat always runs the compile task when running scripts with its command
  // line interface.
  //
  // If this script is run directly using `node` you may want to call compile
  // manually to make sure everything is compiled
  // await hre.run('compile');
  //

  // Addresses with a clean geth instance
  // Signer: 0x21E360e198Cde35740e88572B59f2CAdE421E6b1
  // Bridge:
  //    Address: 0xb588617416D0B0A3C29618bf8Fb6aC0cAd4Ede7f
  //    Relayer: 0x21E360e198Cde35740e88572B59f2CAdE421E6b1
  // WKAVA:
  //    Address: 0x6098c27D41ec6dc280c2200A737D443b0AaA2E8F
  // ERC20:
  //    Name:    Wrapped Dev Ether
  //    Address: 0x8223259205A3E31C54469fCbfc9F7Cf83D515ff6

  const [signer] = await ethers.getSigners();
  console.log("Signing with account %s", signer.address);

  const bridge = await deployBridge(signer.address);
  console.log(
    "Bridge deployed:\n\tAddress: %s\n\tRelayer: %s",
    bridge.address,
    signer.address
  );

  const weth = await deployWETH();
  console.log("WETH deployed:\n\tAddress: %s", weth.address);

  const meowAmounts = new Map<string, BigNumberish>([
    [signer.address, 100_000_000_000_000_000_000n],
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

  console.log("Completed contracts deployment");
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

// Only run main() when the script is run directly
if (require.main === module) {
  // We recommend this pattern to be able to use async/await everywhere
  // and properly handle errors.
  main().catch((error) => {
    console.error(error);
    process.exitCode = 1;
  });
}
