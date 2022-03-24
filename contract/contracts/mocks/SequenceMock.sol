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

import "../Sequence.sol";

contract SequenceMock is Sequence {
    // solhint-disable-next-line no-empty-blocks
    constructor(uint256 initialValue) Sequence(initialValue) {}

    function increment() external {
        incrementSequence();
    }

    function get() external view returns (uint256) {
        return getSequence();
    }
}
