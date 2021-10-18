pragma solidity ^0.8.9;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";

contract Bridge {
    using SafeERC20 for IERC20;

    address private _relayer;

    event Lock(address token, address sender, bytes32 toAddr, uint256 amount);

    event Unlock(address token, address toAddr, uint256 amount);

    constructor(address relayer_) {
        _relayer = relayer_;
    }

    function lock(
        address token,
        bytes32 toAddr,
        uint256 amount
    ) public {
        IERC20(token).safeTransferFrom(msg.sender, address(this), amount);

        emit Lock(token, msg.sender, toAddr, amount);
    }

    function unlock(
        address token,
        address toAddr,
        uint256 amount
    ) public {
        require(
            msg.sender == _relayer,
            "Bridge: untrusted address"
        );

        IERC20(token).safeTransfer(toAddr, amount);

        emit Unlock(token, toAddr, amount);
    }
}
