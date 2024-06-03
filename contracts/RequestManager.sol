// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract RequestManager {
    struct Request {
        bytes4 requestId;
        address requesterAddress;
        address targetAddress;
        uint256 targetChainDomain;
        uint256 requestCreationBlock;
        uint256 requestCreationTime;
        bytes message;
    }

    mapping(bytes4 => Request) public requests;

    function createRequest(bytes4 requestId, address targetAddress, uint256 targetChainDomain, bytes memory message) public {
        requests[requestId] = Request({
            requestId: requestId,
            requesterAddress: msg.sender,
            targetAddress: targetAddress,
            targetChainDomain: targetChainDomain,
            requestCreationBlock: block.number,
            requestCreationTime: block.timestamp,
            message: message
        });
    }

    function getRequest(bytes4 requestId) public view returns (Request memory) {
        return requests[requestId];
    }

    function generateRequestId(uint64 i) public pure returns (bytes4) {
        return bytes4(keccak256(abi.encodePacked(i)));
    }
}
