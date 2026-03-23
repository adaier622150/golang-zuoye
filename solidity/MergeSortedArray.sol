
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract MergeSortedArray {
    // 合并两个有序数组到第一个数组中
    function merge(
        uint256[] memory nums1,
        uint256 m,
        uint256[] memory nums2,
        uint256 n
    ) public pure returns (uint256[] memory) {
        // 创建结果数组
        uint256[] memory result = new uint256[](m + n);

        // 双指针遍历
        uint256 i = 0; // nums1的指针
        uint256 j = 0; // nums2的指针
        uint256 k = 0; // 结果数组的指针

        // 比较两个数组的元素，将较小的放入结果数组
        while (i < m && j < n) {
            if (nums1[i] <= nums2[j]) {
                result[k] = nums1[i];
                i++;
            } else {
                result[k] = nums2[j];
                j++;
            }
            k++;
        }

        // 将nums1中剩余元素复制到结果数组
        while (i < m) {
            result[k] = nums1[i];
            i++;
            k++;
        }

        // 将nums2中剩余元素复制到结果数组
        while (j < n) {
            result[k] = nums2[j];
            j++;
            k++;
        }

        return result;
    }
}
