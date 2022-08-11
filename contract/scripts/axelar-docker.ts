// NOTE: this script doesn't work. :)
// it's included here because it shows you how to import and call axelar contracts.
// it also may contain some info that would be helpful to building a real axelar integration.

import { getDefaultProvider, Contract, Wallet } from "ethers";
import IBurnableMintableCappedERC20 from "@axelar-network/axelar-cgp-solidity/interfaces/IBurnableMintableCappedERC20.sol/IBurnableMintableCappedERC20.json";
import IAxelarGateway from "@axelar-network/axelar-cgp-solidity/interfaces/IAxelarGateway.sol/IAxelarGateway.json";

const ethInfo = {
  name: "Ethereum",
  rpc: "http://localhost:8555",
  usdcAddress: "0x17564a31b5B67441fcF785476D3dEd9E9F2c4937",
  gatewayAddress: "0x6A71747150A8b7F58bfd5222c85122Ac502927F8",
  gasReceiver: "0x5229271b9bE9Fa0CCa9Ee74680D3c4d47C146d3f",
};

const kavaInfo = {
  name: "Kava",
  rpc: "http://localhost:8545",
  usdcAddress: "0xf401F9538AdF6D500D111aec7E54e590b42d8D21",
  gatewayAddress: "0x60D5BE29a0ceb5888F15035d8CcdeACCD5Fd837F",
  gasReceiver: "0xf20684e77C17bc0209f29F7c20F83EB55F5EfA0a",
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

  const usdcEthContract = new Contract(
    ethInfo.usdcAddress,
    IBurnableMintableCappedERC20.abi,
    ethProvider
  );
  const usdcKavaContract = new Contract(
    kavaInfo.usdcAddress,
    IBurnableMintableCappedERC20.abi,
    kavaProvider
  );

  console.log("BALANCES BEFORE:");
  console.log(
    (await usdcEthContract.balanceOf(ethWallet.address)) / 1e6,
    "aUSDC on Ethereum",
    ethWallet.address
  );
  console.log(
    (await usdcEthContract.balanceOf(ethInfo.gatewayAddress)) / 1e6,
    "aUSDC in Ethereum Gateway",
    ethInfo.gatewayAddress
  );
  const beforeBalanceOnKava =
    (await usdcKavaContract.balanceOf(kavaWallet.address)) / 1e6;
  console.log(beforeBalanceOnKava, "aUSDC on Kava EVM", kavaWallet.address);

  const ethGatewayContract = new Contract(
    ethInfo.usdcAddress,
    IAxelarGateway.abi,
    ethProvider
  );
  // const kavaGatewayContract = new Contract(
  //   kavaInfo.usdcAddress,
  //   IAxelarGateway.abi,
  //   kavaProvider
  // );

  // approve gateway to use token on source chain
  console.log("approving eth gateway to transfer 100 aUSDC");
  const ethApproveTx = await usdcEthContract
    .connect(ethWallet)
    .approve(ethInfo.gatewayAddress, 100e6);
  await ethApproveTx.wait();
  console.log("approved.");

  console.log("requesting transfer");
  // ask gateway on source chain to send tokens to destination chain
  const ethGatewayTx = await ethGatewayContract
    .connect(ethWallet)
    .sendToken(kavaInfo.name, kavaWallet.address, "aUSDC", 100e6, {
      gasLimit: 7e6,
    });
  await ethGatewayTx.wait();
  console.log("transferred.");

  // NOTE form pirtle: the script has never made it this far.
  // wait for relay
  console.log("waiting for axelar to relay the tx");
  await withTimeout(
    (async () => {
      let currentBalance = beforeBalanceOnKava;
      while (
        currentBalance ===
        (await usdcKavaContract.balanceOf(kavaWallet.address)) / 1e6
      ) {
        await new Promise((resolve) => setTimeout(resolve, 200));
        currentBalance =
          (await usdcKavaContract.balanceOf(kavaWallet.address)) / 1e6;
        console.log("checked balance:", currentBalance);
      }
      console.log("SUCCESS!");
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
    })(),
    30
  );

  // console.log("relaying transfer!");
  // await relay();
  // console.log(
  //   (await usdcEthContract.balanceOf(ethWallet.address)) / 1e6,
  //   "aUSDC on Ethereum",
  //   ethWallet.address
  // );
  // console.log(
  //   (await usdcKavaContract.balanceOf(kavaWallet.address)) / 1e6,
  //   "aUSDC on Kava EVM",
  //   kavaWallet.address
  // );
}

/** resolves as the promise or an error if longer than timeoutSec seconds have passed. */
async function withTimeout<T>(p: Promise<T>, timeoutSec: number) {
  const timeout = new Promise((resolve, reject) => {
    setTimeout(
      () => reject(new Error("operation timed out")),
      timeoutSec * 1000
    );
  });
  return Promise.race([p, timeout]);
}

if (require.main === module) {
  // setupNetworkContracts().catch((error) => {
  //   console.error(error);
  //   process.exitCode = 1;
  // });
  main().catch((error) => {
    console.error(error);
    process.exitCode = 1;
  });
}

// Setting up Ethereum on a network with a chainId of 88881...
// Deploying the ConstAddressDeployer for Ethereum...
// Deployed at 0xb44F2eF7fF2ba62Fe06D4c5818AfA972e2B90Dcb
// Deploying the Axelar Gateway for Ethereum...
// Deployed at 0x6A71747150A8b7F58bfd5222c85122Ac502927F8
// Deploying the Axelar Gas Receiver for Ethereum...
// Deployed at 0x5229271b9bE9Fa0CCa9Ee74680D3c4d47C146d3f
// Deploying Axelar Wrapped USDC for Ethereum...
// Deployed at 0x17564a31b5B67441fcF785476D3dEd9E9F2c4937
// Funding 0x21E360e198Cde35740e88572B59f2CAdE421E6b1 with 1000 USDC

// Deploying the ConstAddressDeployer for Kava...
// Deployed at 0x69aeB7Dc4f2A86873Dae8D753DE89326Cf90a77a
// Deploying the Axelar Gateway for Kava...
// Deployed at 0x60D5BE29a0ceb5888F15035d8CcdeACCD5Fd837F
// Deploying the Axelar Gas Receiver for Kava...
// Deployed at 0xf20684e77C17bc0209f29F7c20F83EB55F5EfA0a
// Deploying USDC for Kava...
// Deployed at 0xf401F9538AdF6D500D111aec7E54e590b42d8D21
