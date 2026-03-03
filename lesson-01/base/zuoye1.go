package main

import (
	"fmt"
	"strconv"
)

// 1. 只出现一次的数字
// 给定一个非空整数数组，除了某个元素只出现一次以外，其余每个元素均出现两次。找出那个只出现了一次的元素。
func SingleNumber(nums []int) int {

	map1 := make(map[int]int)
	for _, value := range nums {
		if _, ok := map1[value]; ok {
			map1[value]++
		} else {
			map1[value] = 1
		}
	}
	for key, value := range map1 {
		if value == 1 {
			return key
		}
	}
	return 0
}

// 2. 回文数
// 判断一个整数是否是回文数
func IsPalindrome(x int) bool {

	y := strconv.Itoa(x)
	for i := 0; i < len(y)/2; i++ {
		if y[i] != y[len(y)-i-1] {
			return false
		}
	}
	return true
}

// 3. 有效的括号
// 给定一个只包括 '(', ')', '{', '}', '[', ']' 的字符串，判断字符串是否有效
func IsValid(s string) bool {
	x := [3]string{}
	l := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '(' || s[i] == '{' || s[i] == '[' {
			x[l] = string(s[i])
			l++
		} else {
			if s[i] == ')' && x[l-1] == "(" {
				x[l-1] = ""
				l--
			} else if s[i] == '}' && x[l-1] == "{" {
				x[l-1] = ""
				l--
			} else if s[i] == ']' && x[l-1] == "[" {
				x[l-1] = ""
				l--
			} else {
				return false
			}
		}
	}
	return l == 0
}

// 4. 最长公共前缀
// 查找字符串数组中的最长公共前缀
func LongestCommonPrefix(strs []string) string {
	pre := ""
	prefix := ""
	for i := 0; true; i++ {
		pre = string(strs[0][i])
		for _, value := range strs {
			if pre != string(value[i]) {
				return prefix
			}
		}
		prefix = prefix + pre
	}
	return ""
}

// 5. 加一
// 给定一个由整数组成的非空数组所表示的非负整数，在该数的基础上加一
func PlusOne(digits []int) []int {
	nums := make([]int, len(digits)+1, len(digits)+1)
	a := 1
	for i := len(digits) - 1; i >= 0; i-- {
		if digits[i]+a > 9 {
			nums[i+1] = 0
			a = 1
		} else {
			nums[i+1] = digits[i] + a
			a = 0
		}
	}
	if a == 1 {
		nums[0] = 1
		return nums
	}
	return nums[1:len(nums)]
}

// 6. 删除有序数组中的重复项
// 给你一个有序数组 nums ，请你原地删除重复出现的元素，使每个元素只出现一次，返回删除后数组的新长度。
// 不要使用额外的数组空间，你必须在原地修改输入数组并在使用 O(1) 额外空间的条件下完成。
func RemoveDuplicates(nums []int) int {
	var x = -1000
	for a := 0; a < len(nums)-1; a++ {
		if nums[a] == x {
			continue
		}
		for b := a + 1; b < len(nums); b++ {
			if nums[b] == x {
				continue
			}
			if nums[a] == nums[b] {
				nums[b] = x
			}
		}
	}
	c := 0
	for a := 1; a < len(nums); a++ {
		if nums[a] == x {
			if c == 0 {
				c = a
			}
		} else if c != 0 {
			nums[c] = nums[a]
			nums[a] = x
			a = c
			c = 0
		}
	}
	return c
}

// 7. 合并区间
// 以数组 intervals 表示若干个区间的集合，其中单个区间为 intervals[i] = [starti, endi] 。
// 请你合并所有重叠的区间，并返回一个不重叠的区间数组，该数组需恰好覆盖输入中的所有区间。
func Merge(intervals [][]int) [][]int {
	r := make([][]int, 0, len(intervals))
	for i := 0; i < len(intervals); i++ {
		if i == 0 {
			r = append(r, intervals[i])
		} else {
			if intervals[i][0] <= r[len(r)-1][1] {
				if intervals[i][1] > r[len(r)-1][1] {
					r[len(r)-1][1] = intervals[i][1]
				}
			} else {
				r = append(r, intervals[i])
			}
		}
	}
	return r
}

// 8. 两数之和
// 给定一个整数数组 nums 和一个目标值 target，请你在该数组中找出和为目标值的那两个整数
func TwoSum(nums []int, target int) []int {
	m := make(map[int]int)
	for i := 0; i < len(nums); i++ {
		if value, ok := m[target-nums[i]]; ok {
			return []int{value, i}
		}
		m[nums[i]] = i
	}
	return nil
}

func main() {
	//nums := [][]int{{1, 3}, {2, 6}, {8, 10}, {15, 18}}
	//nums := []int{2, 7, 11, 15}
	//target := 22
	//fmt.Println("合并区间", nums)
	//x := TwoSum(nums, target)
	//fmt.Println("合并区间", x)
	m := make(map[int]int)
	m[1] = 1
	value, ok := m[2]
	fmt.Println("合并区间", value)
	fmt.Println("合并区间", ok)
}
