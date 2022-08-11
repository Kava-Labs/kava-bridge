import {
  setupNetwork,
  relay,
  getNetwork,
} from "@axelar-network/axelar-local-dev";
import { ethers } from "hardhat";
const { getDefaultProvider, Wallet } = ethers;

// update these addresses after running the setup!
const addresses = {
  eth: {
    constAddressDeployer: "0xdB59Bbd48435dd01e60775160b77baB900B8C0db",
    gateway: "0xF478dA1fFd00F94Fde70D8f2AD9849E62f1dCb03",
    gasReceiver: "0x25Bb16cDCE27d954276040aDFD463C1b6326F1ce",
    usdc: "0x7D4F9441d0Df048918D676ED640653c9324FeBE2",
  },
  kava: {
    constAddressDeployer: "0x881405825674671d81231e0bd4A39456C7b30677",
    gateway: "0x0dC81Ba8fefe0525234838CAFc5E1B7595808601",
    gasReceiver: "0xdD1Cc8dafB6b22d86FF5392B35CBe4223267Cf3e",
    usdc: "0x3F8511D40551f93B0BD498e9C18fE00b5f4e68A2",
  },
};

// miner of local geth
const key =
  "0xa19b60b9de368bffed131ff991d391a33a5adf1ca96e325e5d028008c1b18cd5";
const wallet = new Wallet(key);

async function main() {
  // setup ethereum connection
  const ethProvider = getDefaultProvider("http://localhost:8555");
  const ethWallet = wallet.connect(ethProvider);

  const kavaProvider = getDefaultProvider("http://localhost:8545");
  const kavaWallet = wallet.connect(kavaProvider);

  // setup kava connection
  const eth = await getNetwork(ethProvider, {
    name: "Ethereum",
    chainId: 88881,
    userKeys: [key],
    ownerKey: key,
    operatorKey: key,
    relayerKey: key,
    adminKeys: [key],
    threshold: 1,
    lastRelayedBlock: await ethProvider.getBlockNumber(),
    gatewayAddress: addresses.eth.gateway,
    gasReceiverAddress: addresses.eth.gasReceiver,
    constAddressDeployerAddress: addresses.eth.constAddressDeployer,
    tokens: {
      aUSDC: "aUSDC",
    },
  });

  // deploy kava token
  const kava = await getNetwork(kavaProvider, {
    name: "Kava",
    chainId: 8888,
    userKeys: [key],
    ownerKey: key,
    operatorKey: key,
    relayerKey: key,
    adminKeys: [key],
    threshold: 1,
    lastRelayedBlock: await kavaProvider.getBlockNumber(),
    gatewayAddress: addresses.kava.gateway,
    gasReceiverAddress: addresses.kava.gasReceiver,
    constAddressDeployerAddress: addresses.kava.constAddressDeployer,
    tokens: {
      aUSDC: "aUSDC",
    },
  });

  // mint tokens on source chain
  // console.log("minting 100 aUSDC on ethereum");
  // await eth.giveToken(wallet.address, "aUSDC", BigInt(100e6));

  const usdcEthContract = await eth.getTokenContract("aUSDC");
  const usdcKavaContract = await kava.getTokenContract("aUSDC");
  console.log("USDC contract addresses:");
  console.log("Ethereum:", usdcEthContract.address);
  console.log("Kava:", usdcKavaContract.address);

  console.log("BALANCES BEFORE:");
  console.log(
    (await usdcEthContract.balanceOf(ethWallet.address)) / 1e6,
    "aUSDC on Ethereum",
    ethWallet.address
  );
  console.log(
    (await usdcKavaContract.balanceOf(kavaWallet.address)) / 1e6,
    "aUSDC on Kava EVM",
    kavaWallet.address
  );

  // approve gateway to use token on source chain
  console.log("approving eth gateway to transfer 100 aUSDC");
  const ethApproveTx = await usdcEthContract
    .connect(ethWallet)
    .approve(eth.gateway.address, 100e6);
  await ethApproveTx.wait();

  console.log("requesting transfer");
  // ask gateway on source chain to send tokens to destination chain
  const ethGatewayTx = await eth.gateway
    .connect(ethWallet)
    .sendToken(kava.name, kavaWallet.address, "aUSDC", 100e6);
  await ethGatewayTx.wait();

  // relay transactions
  console.log("relaying transfer!");
  await relay();
  console.log(
    (await usdcEthContract.balanceOf(ethWallet.address)) / 1e6,
    "aUSDC on Ethereum",
    ethWallet.address
  );
  console.log(
    (await usdcKavaContract.balanceOf(kavaWallet.address)) / 1e6,
    "aUSDC on Kava EVM",
    kavaWallet.address
  );
}

export async function setupNetworkContracts() {
  // setup ethereum network
  const eth = await setupNetwork("http://localhost:8555", {
    name: "Ethereum",
    chainId: 88881,
    ownerKey: wallet,
    userKeys: [],
    adminKeys: [wallet],
    operatorKey: wallet,
    relayerKey: wallet,
    threshold: 1,
    lastRelayedBlock: 1,
  });
  await eth.deployToken("Axelar Wrapped USDC", "aUSDC", 6, BigInt(100_000e6));
  console.log(`Funding ${wallet.address} with 1000 USDC`);
  await eth.giveToken(wallet.address, "aUSDC", BigInt(10000e6));

  // setup kava network
  const kava = await setupNetwork("http://localhost:8545", {
    name: "Kava",
    chainId: 8888,
    ownerKey: wallet,
    userKeys: [],
    adminKeys: [wallet],
    operatorKey: wallet,
    relayerKey: wallet,
    threshold: 1,
    lastRelayedBlock: 1,
  });
  await kava.deployToken("USDC", "aUSDC", 6, BigInt(100_000e6));

  return { eth, kava };
}

if (require.main === module) {
  // uncomment me for the setup!
  // setupNetworkContracts().catch((error) => {
  //   console.error(error);
  //   process.exitCode = 1;
  // });
  // uncomment me for the running!
  // main().catch((error) => {
  //   console.error(error);
  //   process.exitCode = 1;
  // });
}
