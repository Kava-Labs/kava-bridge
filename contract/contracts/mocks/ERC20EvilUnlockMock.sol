// SPDX-License-Identifier: MIT
// Derived from https://github.com/OpenZeppelin/openzeppelin-contracts/blob/v4.3.2/contracts/mocks/SafeERC20Helper.sol

pragma solidity ^0.8.9;

import "../Bridge.sol";

contract ERC20EvilUnlockMock {
    uint256 private _reentered;

    function transfer(address to, uint256 amount) public returns (bool) {
        if (_reentered == 0) {
            _reentered = 1;
            Bridge(msg.sender).unlock(address(this), to, amount);
        }

        return true;
    }
}
