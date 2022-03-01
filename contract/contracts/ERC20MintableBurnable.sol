// SPDX-License-Identifier: MIT

pragma solidity ^0.8.9;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

/// @title A contract for an mintable and burnable ERC20 that is handled by the bridge Cosmos SDK module account
/// @author Kava Labs, LLC
/// @custom:security-contact security@kava.io
contract ERC20MintableBurnable is ERC20, Ownable {
    /// @notice The decimals places of the token.
    uint8 private immutable _decimals;

    /// @notice Registers the ERC20 token with mint and burn permissions for the
    ///         contract owner, by default the account that deploys this contract.
    /// @param name The name of the ERC20 token.
    /// @param symbol The symbol of the ERC20 token.
    /// @param decimals_ The number of decimals of the ERC20 token.
    constructor(
        string memory name,
        string memory symbol,
        uint8 decimals_
    ) ERC20(name, symbol) {
        _decimals = decimals_;
    }

    /// @notice Creates more supply for a given account and amount by contract owner.
    /// @param to The account to send the minted supply to.
    /// @param amount The amount of the token to mint.
    function mint(address to, uint256 amount) public onlyOwner {
        _mint(to, amount);
    }

    function decimals() public view override returns (uint8) {
        return _decimals;
    }
}
