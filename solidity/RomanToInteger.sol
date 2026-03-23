
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract RomanToInteger {
    function romanToInt(string memory s) public pure returns (uint256) {
        bytes memory input = bytes(s);
        uint256 result = 0;
        uint256 i = 0;

        while (i < input.length) {
            if (i + 1 < input.length) {
                // 检查是否为特殊组合（如IV, IX, XL等）
                if (input[i] == 'I' && input[i+1] == 'V') {
                    result += 4;
                    i += 2;
                    continue;
                }
                if (input[i] == 'I' && input[i+1] == 'X') {
                    result += 9;
                    i += 2;
                    continue;
                }
                if (input[i] == 'X' && input[i+1] == 'L') {
                    result += 40;
                    i += 2;
                    continue;
                }
                if (input[i] == 'X' && input[i+1] == 'C') {
                    result += 90;
                    i += 2;
                    continue;
                }
                if (input[i] == 'C' && input[i+1] == 'D') {
                    result += 400;
                    i += 2;
                    continue;
                }
                if (input[i] == 'C' && input[i+1] == 'M') {
                    result += 900;
                    i += 2;
                    continue;
                }
            }

            // 处理普通字符
            if (input[i] == 'I') {
                result += 1;
            } else if (input[i] == 'V') {
                result += 5;
            } else if (input[i] == 'X') {
                result += 10;
            } else if (input[i] == 'L') {
                result += 50;
            } else if (input[i] == 'C') {
                result += 100;
            } else if (input[i] == 'D') {
                result += 500;
            } else if (input[i] == 'M') {
                result += 1000;
            }
            i += 1;
        }

        return result;
    }
}
