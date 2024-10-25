// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract Coinbase {
    uint256 balance;
    function touchAddress(address addr) public {
        balance = addr.balance;
    }

    function touchCoinbase() public {
        touchAddress(block.coinbase);
    }

    event LogAddress(address addr);
    function logCoinBaseAddress() public {
        emit LogAddress(block.coinbase);
    }
}
