
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract BinarySearch {
    /**
     * @dev 在有序数组中使用二分查找算法查找目标值
     * @param arr 有序数组（升序排列）
     * @param target 要查找的目标值
     * @return index 目标值在数组中的索引，如果不存在则返回-1
     */
    function binarySearch(int256[] memory arr, int256 target) public pure returns (int256) {
        uint256 left = 0;
        uint256 right = arr.length - 1;

        while (left <= right) {
            uint256 mid = left + (right - left) / 2;

            if (arr[mid] == target) {
                return int256(mid);
            } else if (arr[mid] < target) {
                left = mid + 1;
            } else {
                if (mid > 0) {
                    right = mid - 1;
                } else {
                    break;
                }
            }
        }

        return -1; // 未找到目标值
    }


}
