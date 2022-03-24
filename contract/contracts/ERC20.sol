// SPDX-License-Identifier: Apache-2.0

pragma solidity ^0.8.9;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";

// Erc20 is used to generate bindings only and not is used for deployment
contract Erc20 is ERC20 {
    constructor() ERC20("ABIToken", "ABIToken") {}
}
