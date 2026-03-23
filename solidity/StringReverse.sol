
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract StringReverse {
    function reverseString(string memory input) public pure returns (string memory) {
        bytes memory inputBytes = bytes(input);
        bytes memory reversed = new bytes(inputBytes.length);

        for (uint256 i = 0; i < inputBytes.length; i++) {
            reversed[inputBytes.length - 1 - i] = inputBytes[i];
        }

        return string(reversed);
    }
}
