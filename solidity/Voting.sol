BeggingContract3.sol// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract Voting {
    string[] private  candidateArray;

    mapping(string => uint256) private candidateVoteMap;

    function getList() public view returns( uint256 len) {

         len = candidateArray.length;

    }
    function getCandidateList() public view returns(string memory) {
        string memory votes = "";
        uint len = candidateArray.length;
         for(uint i = 0; i < len ; i++) {
            votes = string.concat(votes, " ", candidateArray[i]);
        }
        return votes;

    }

    function vote(string calldata _candidate) public{
        uint256 candidateVoteCount = candidateVoteMap[_candidate];
        if (candidateVoteCount == 0){
            candidateArray.push(_candidate);
        }

            candidateVoteMap[_candidate] = candidateVoteCount + 1;
    }
    function getVotes(string calldata _candidate) public view returns( uint256 candidateVoteCount) {
          candidateVoteCount = candidateVoteMap[_candidate];

    }
    function resetVotes() public {
        uint len = candidateArray.length;
        for(uint i = 0; i < len ; i++) {
            candidateVoteMap[candidateArray[i]] = 0;
        }
        delete candidateArray;

    }


}
