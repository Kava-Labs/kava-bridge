// SPDX-License-Identifier: MIT
// Derived from https://github.com/OpenZeppelin/openzeppelin-contracts/blob/v4.3.2/contracts/mocks/SafeERC20Helper.sol

pragma solidity ^0.8.9;

import "../Bridge.sol";

contract ERC20EvilLockMock {
  uint256 private _reentered;
  bytes32 private _toKavaAddr;
  Bridge private _bridge;

  constructor(address _bridgeAddress) {
    _bridge = Bridge(_bridgeAddress);
  }

  function attack(bytes32 toAddr, uint256 amount) external {
    _toKavaAddr = toAddr;
    _reentered = 0;
    _bridge.lock(address(this), toAddr, amount);
  }

  function transferFrom(
    address,
    address,
    uint256 amount
  ) public returns (bool) {
    if (_reentered < 3) {
      _reentered = _reentered + 1;
      _bridge.lock(address(this), _toKavaAddr, amount);
    }

    return true;
  }
}
