//SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract MultiSigWallet {
    address[] public owners;
    uint public threshold;
    mapping(address => bool) public isOwner;

    struct Transaction {
        address to;
        uint value;
        bool executed;
        uint confirmations;
    }

    Transaction[] public transactions;

    constructor(address[] memory _owners, uint _threshold) {
        require(_owners.length > 0, "Owners required");
        require(_threshold <= _owners.length, "Invalid threshold");

        for (uint i=0; i< _owners.length; i++) {
            isOwner[_owners[i]] = true;
        }
        owners = _owners;
        threshold = _threshold;

    }

    function submitTransaction(address to, uint value) public onlyOwner {
       transactions.push(Transaction({
        to: to,
        value: value,
        executed:false,
        confirmations: 0
    }));
    }

    function confirmTransaction(uint txIndex) internal {
        Transaction storage transaction =  transactions[txIndex];
        require(!transaction.executed, "Already executed");

        transaction.executed = true;
        (bool success, ) = transaction.to.call{value: transaction.value}("");
        require(success, "Transaction failed");
    }

    modifier onlyOwner() {
        require(isOwner[msg.sender], "Not an owner");
        _;
    }

}
