// SPDX-License-Identifier: MIT
// Derived from https://github.com/OpenZeppelin/openzeppelin-contracts/blob/v4.3.2/contracts/mocks/SafeERC20Helper.sol

pragma solidity ^0.8.9;

import "../Sequence.sol";

contract SequenceMock is Sequence {
    constructor(uint256 initialValue) Sequence(initialValue) {}

    function increment() external {
        incrementSequence();
    }

    function get() external view returns (uint256) {
        return getSequence();
    }
}
