// SPDX-License-Identifier: Apache-2.0
// Derived from https://github.com/OpenZeppelin/openzeppelin-contracts/blob/v4.3.2/contracts/mocks/SafeERC20Helper.sol
//
// The MIT License (MIT)
//
// Copyright (c) 2016-2022 zOS Global Limited and contributors
//
// Permission is hereby granted, free of charge, to any person obtaining
// a copy of this software and associated documentation files (the
// "Software"), to deal in the Software without restriction, including
// without limitation the rights to use, copy, modify, merge, publish,
// distribute, sublicense, and/or sell copies of the Software, and to
// permit persons to whom the Software is furnished to do so, subject to
// the following conditions:
//
// The above copyright notice and this permission notice shall be included
// in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS
// OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY
// CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT,
// TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE
// SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

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
