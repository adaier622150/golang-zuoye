
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract RomanNumeral {
    function intToRoman(uint256 num) public pure returns (string memory) {
        if (num == 0) {
            return "N";
        }

        uint256[13] memory values = [
            uint256(1000), 900, 500, 400,
            100, 90, 50, 40,
            10, 9, 5, 4,
            1
        ];

        string[13] memory symbols = [
            "M", "CM", "D", "CD",
            "C", "XC", "L", "XL",
            "X", "IX", "V", "IV",
            "I"
        ];

        string memory result = "";
        for (uint256 i = 0; i < 13; i++) {
            while (num >= values[i]) {
                //result = string(abi.encodePacked(result, symbols[i]));
                result = string.concat(result, symbols[i]);
                num -= values[i];
            }
        }

        return result;
    }

}
