// SPDX-License-Identifier: MIT
// Derived from https://github.com/OpenZeppelin/openzeppelin-contracts/blob/v4.3.2/contracts/mocks/SafeERC20Helper.sol

pragma solidity ^0.8.9;

contract ERC20ReturnTrueMock {
    mapping(address => uint256) private _allowances;

    // IERC20's functions are not pure, but these mock implementations are: to prevent Solidity from issuing warnings,
    // we write to a dummy state variable.
    uint256 private _dummy;

    function transfer(address, uint256) public returns (bool) {
        _dummy = 0;
        return true;
    }

    function transferFrom(
        address,
        address,
        uint256
    ) public returns (bool) {
        _dummy = 0;
        return true;
    }

    function approve(address, uint256) public returns (bool) {
        _dummy = 0;
        return true;
    }

    function setAllowance(uint256 allowance_) public {
        _allowances[msg.sender] = allowance_;
    }

    function allowance(address owner, address) public view returns (uint256) {
        return _allowances[owner];
    }
}
