// SPDX-License-Identifier: MIT

pragma solidity ^0.8.9;

import "../Sequence.sol";
import "../ERC20MintableBurnable.sol";

contract Withdrawer {
    function withdrawFor(
        address contractAddr,
        address account,
        address toAddr,
        uint256 amount
    ) external {
        ERC20MintableBurnable(contractAddr).withdrawFrom(
            account,
            toAddr,
            amount
        );
    }
}
