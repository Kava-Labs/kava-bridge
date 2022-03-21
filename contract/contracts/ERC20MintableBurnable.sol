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

    /// @notice Represents an ERC20 token lock emitted during a lock call
    /// @param sender The Kava address of the sender that locked the funds
    /// @param toAddr The Ethereum address to send the funds to
    /// @param amount The amount that was locked
    event Withdraw(
        address indexed sender,
        address indexed toAddr,
        uint256 amount
    );

    /// @notice Represents a conversion from ERC20 to sdk.Coin
    /// @param sender The Kava address of the sender that converted coins
    /// @param toKavaAddr The Kava address where to send the converted coins to
    /// @param amount The amount that was converted
    event ConvertToCoin(
        address indexed sender,
        bytes32 indexed toKavaAddr,
        uint256 amount
    );

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

    /// @notice Withdraws `amount` of tokens to Ethereum from the caller.
    /// @dev Destroys `amount` tokens from the caller and emits a withdraw
    ///      event for the relayer to unlock funds on Ethereum.
    /// @param toAddr The account on Ethereum to withdraw the funds to.
    /// @param amount The amount of the token to withdraw.
    function withdraw(address toAddr, uint256 amount) public virtual {
        _burn(msg.sender, amount);
        emit Withdraw(msg.sender, toAddr, amount);
    }

    /// @notice Withdraws `amount` of tokens to Ethereum from `account`.
    /// @dev Destroys `amount` tokens from the caller, deducts from the caller's
    ///      allowance, and emits a withdraw event for the relayer to unlock
    ///      funds on Ethereum.
    ///
    ///      See {ERC20-_burn} and {ERC20-allowance}.
    ///
    ///      Requirements:
    ///      - the caller must have allowance for ``accounts``'s tokens of at
    ///        least `amount`.
    /// @param toAddr The account on Ethereum to withdraw the funds to.
    /// @param amount The amount of the token to withdraw.
    function withdrawFrom(
        address account,
        address toAddr,
        uint256 amount
    ) public virtual {
        _spendAllowance(account, msg.sender, amount);
        _burn(account, amount);
        emit Withdraw(msg.sender, toAddr, amount);
    }

    /// @notice Converts an amount of tokens to Cosmos sdk.Coin
    /// @dev Transfers amount of tokens to the contract address and emits a
    ///      ConvertToCoin event.
    /// @param toKavaAddr The Kava address where to send the converted coins to.
    /// @param amount The amount of the token to convert.
    function convertToCoin(bytes32 toKavaAddr, uint256 amount) public virtual {
        _transfer(msg.sender, address(this), amount);
        emit ConvertToCoin(msg.sender, toKavaAddr, amount);
    }
}
