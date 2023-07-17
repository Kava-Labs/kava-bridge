import * as dotenv from "dotenv";

import { HardhatUserConfig, task } from "hardhat/config";
import "@nomiclabs/hardhat-etherscan";
import "@nomiclabs/hardhat-ethers";
import "@nomiclabs/hardhat-waffle";
import "@typechain/hardhat";
import "hardhat-gas-reporter";
import "hardhat-watcher";
import "solidity-coverage";

import "./tasks/deployERC20";
import "./tasks/mintERC20";

dotenv.config();

// This is a sample Hardhat task. To learn how to create your own go to
// https://hardhat.org/guides/create-task.html
task("accounts", "Prints the list of accounts", async (taskArgs, hre) => {
  const accounts = await hre.ethers.getSigners();

  for (const account of accounts) {
    console.log(account.address);
  }
});

// You need to export an object to set up your config
// Go to https://hardhat.org/config/ to learn more

const config: HardhatUserConfig = {
  solidity: "0.8.9",
  networks: {
    ropsten: {
      url: process.env.ROPSTEN_URL || "",
      accounts:
        process.env.PRIVATE_KEY !== undefined ? [process.env.PRIVATE_KEY] : [],
    },
    localhost: {
      url: "http://127.0.0.1:8555",
      accounts: "remote",
    },
    kava17000: {
      url: "https://evm.app.internal.testnet-17000.us-east.production.kava.io:443",
      accounts: [
        // testnet user key
      ],
    },
    demonet: {
      url: "https://evm.data.demonet.us-east.production.kava.io:443",
      accounts: [
        "247069F0BC3A5914CB2FD41E4133BBDAA6DBED9F47A01B9F110B5602C6E4CDD9",
      ],
    },
    protonet: {
      url: "https://evm.app.protonet.us-east.production.kava.io:443",
      accounts: [
        "247069F0BC3A5914CB2FD41E4133BBDAA6DBED9F47A01B9F110B5602C6E4CDD9",
      ],
    },
    internal_testnet: {
      url: "https://evm.app.internal.testnet.us-east.production.kava.io:443",
      accounts: [
        "247069F0BC3A5914CB2FD41E4133BBDAA6DBED9F47A01B9F110B5602C6E4CDD9",
      ],
    },
    kava: {
      url: "http://127.0.0.1:8545",
      accounts: [
        // kava keys unsafe-export-eth-key user --keyring-backend test
        "9549F115B0A21E5071A8AEC1B74AC093190E18DD83D019AC6497B0ADFBEFF26D",
      ],
    },
  },
  gasReporter: {
    enabled: process.env.REPORT_GAS !== undefined,
    currency: "USD",
  },
  etherscan: {
    apiKey: process.env.ETHERSCAN_API_KEY,
  },
  watcher: {
    dev: {
      tasks: ["compile", "test"],
      files: ["./contracts/**/*.sol", "./test/**/*", "./scripts/**/*"],
      verbose: true,
    },
    test: {
      tasks: [{ command: "test", params: { testFiles: ["{path}"] } }],
      files: ["./test/**/*"],
      verbose: true,
    },
  },
  typechain: {
    outDir: "typechain-types",
    target: "ethers-v5",
    alwaysGenerateOverloads: false,
    externalArtifacts: ["artifacts/@openzepplin/*.json"],
  },
};

export default config;
