// SPDX-License-Identifier: MIT

pragma solidity ^0.8.9;

/// @title A contract for an incrementing sequence
/// @author Kava Labs, LLC
abstract contract Sequence {
    /// @notice The incrementing sequence. This *can* overflow and is expected behavior.
    uint256 private _sequence;

    /// @notice Initialize with a start sequence value
    /// @param startValue The start value for the sequence
    constructor(uint256 startValue) {
        _sequence = startValue;
    }

    /// @notice Increment the sequence, *allowing* overflows.
    function incrementSequence() internal {
        unchecked {
            _sequence = _sequence + 1;
        }
    }

    /// @notice Get the current sequence
    /// @return The current sequence
    function getSequence() internal view returns (uint256) {
        return _sequence;
    }
}
