// SPDX-License-Identifier: MIT
// Derived from https://github.com/OpenZeppelin/openzeppelin-contracts/blob/v4.3.2/contracts/mocks/SafeERC20Helper.sol

pragma solidity ^0.8.9;

import "../Bridge.sol";

contract ERC20EvilMock {
    uint256 private _reentered;
    bytes32 private _toKavaAddr;

    constructor() {
        _reentered = 0;
    }

    // Lock attack requires a separate function to set thei kava address
    function attackLock(
        address target,
        bytes32 toAddr,
        uint256 amount
    ) external {
        _toKavaAddr = toAddr;
        Bridge(target).lock(address(this), toAddr, amount);
    }

    // Bridge lock
    function transferFrom(
        address,
        address,
        uint256 amount
    ) public returns (bool) {
        if (_reentered < 3) {
            _reentered = _reentered + 1;
            Bridge(msg.sender).lock(address(this), _toKavaAddr, amount);
        }

        return true;
    }

    // Bridge unlock
    function transfer(address to, uint256 amount) public returns (bool) {
        if (_reentered < 3) {
            _reentered = _reentered + 1;
            Bridge(msg.sender).unlock(address(this), to, amount, 1);
        }

        return true;
    }
}
