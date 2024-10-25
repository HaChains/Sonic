// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract Coinbase {
    uint256 balance;
    function touchAddress(address addr) public {
        balance = addr.balance;
    }

    event LogAddress(address addr);
    function logCoinBaseAddress() public {
        emit LogAddress(block.coinbase);
    }
}
