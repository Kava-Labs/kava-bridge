// SPDX-License-Identifier: MIT

pragma solidity ^0.8.9;

/// @title A contract for WKAVA, a wrapped ERC20 version of KAVA
/// @author Kava Labs, LLC
/// @custom:security-contact security@kava.io
contract WKAVA {
    string public name = "Wrapped Kava";
    string public symbol = "WKAVA";
    uint8 public decimals = 6;

    event Approval(address indexed src, address indexed account, uint256 wad);
    event Transfer(address indexed src, address indexed dst, uint256 wad);
    event Deposit(address indexed src, uint256 wad);
    event Withdrawal(address indexed src, uint256 wad);

    mapping(address => uint256) public balanceOf;
    mapping(address => mapping(address => uint256)) public allowance;

    /// @notice Acts as a fallback function that accepts Ethereum directly
    receive() external payable {
        deposit();
    }

    /// @notice Convert  KAVA  function that accepts Ethereum directly
    function deposit() public payable {
        balanceOf[msg.sender] += msg.value;
        emit Deposit(msg.sender, msg.value);
    }

    /// @notice Convert wkava back to KAVA
    /// @param wad The amount of wkava to withdraw
    function withdraw(uint256 wad) public {
        require(balanceOf[msg.sender] >= wad);
        balanceOf[msg.sender] -= wad;
        payable(msg.sender).transfer(wad);
        emit Withdrawal(msg.sender, wad);
    }

    /// @notice Returns the total supply of wkava
    function totalSupply() public view returns (uint256) {
        return address(this).balance;
    }

    /// @notice Convert wkava back to KAVA
    /// @param account The account address to grant approval to
    /// @param amount The amount of wkava to grant approval of
    function approve(address account, uint256 amount) public returns (bool) {
        allowance[msg.sender][account] = amount;
        emit Approval(msg.sender, account, amount);
        return true;
    }

    /// @notice Transfer wkava to an address
    /// @param dst The destination address
    /// @param wad The amount of wkava to transfer
    function transfer(address dst, uint256 wad) public returns (bool) {
        return transferFrom(msg.sender, dst, wad);
    }

    /// @notice Transfer approved wkava to an address
    /// @param src The source address
    /// @param dst The destination address
    /// @param wad The amount of wkava to transfer
    function transferFrom(
        address src,
        address dst,
        uint256 wad
    ) public returns (bool) {
        require(balanceOf[src] >= wad);

        if (
            src != msg.sender && allowance[src][msg.sender] != type(uint256).max
        ) {
            require(allowance[src][msg.sender] >= wad);
            allowance[src][msg.sender] -= wad;
        }

        balanceOf[src] -= wad;
        balanceOf[dst] += wad;

        emit Transfer(src, dst, wad);

        return true;
    }
}
