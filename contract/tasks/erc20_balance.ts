import { task } from "hardhat/config";
import "@nomiclabs/hardhat-ethers";

task("erc20:balanceOf", "Prints an account's ERC20 balance")
  .addParam("contract", "The contract address")
  .addParam("account", "The account's address")
  .setAction(async ({ contract, account }, { ethers }) => {
    const factory = await ethers.getContractFactory("ERC20");
    const erc20 = factory.attach(contract);

    const bal = await erc20.balanceOf(account);
    console.log(bal.toString());
  });
