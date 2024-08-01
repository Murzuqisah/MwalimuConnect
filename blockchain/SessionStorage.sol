// SPDX-License-Identifier: MIT
pragma solidity ^0.5.0;

contract SessionStorage {
    mapping(address => string) public userSessions;

    function storeSessionID(string memory sessionId) public {
        userSessions[msg.sender] = sessionId;
    }

    function getSessionID() public view returns (string memory) {
        return userSessions[msg.sender];
    }
}
