package utils

import (
	"math/rand"
	"sort"
	"time"
)

// ContainsInt 判断整数切片是否包含指定值。
func ContainsInt(slice []int, val int) bool {
	for _, v := range slice {
		if v == val {
			return true
		}
	}
	return false
}

// ContainsString 判断字符串切片是否包含指定值。
func ContainsString(slice []string, val string) bool {
	for _, v := range slice {
		if v == val {
			return true
		}
	}
	return false
}

// UniqueStrings 去重字符串切片。
func UniqueStrings(slice []string) []string {
	seen := make(map[string]bool, len(slice))
	result := make([]string, 0, len(slice))
	for _, s := range slice {
		if !seen[s] {
			seen[s] = true
			result = append(result, s)
		}
	}
	return result
}

// UniqueInts 去重整数切片。
func UniqueInts(slice []int) []int {
	seen := make(map[int]bool, len(slice))
	result := make([]int, 0, len(slice))
	for _, n := range slice {
		if !seen[n] {
			seen[n] = true
			result = append(result, n)
		}
	}
	return result
}

// ShuffleStrings 随机打乱字符串切片。
func ShuffleStrings(slice []string) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Shuffle(len(slice), func(i, j int) {
		slice[i], slice[j] = slice[j], slice[i]
	})
}

// SortStringsDesc 字符串切片降序排序。
func SortStringsDesc(slice []string) {
	sort.Sort(sort.Reverse(sort.StringSlice(slice)))
}

// ChunkStrings 将字符串切片分块。
func ChunkStrings(slice []string, size int) [][]string {
	if size <= 0 {
		return [][]string{slice}
	}
	var chunks [][]string
	for i := 0; i < len(slice); i += size {
		end := i + size
		if end > len(slice) {
			end = len(slice)
		}
		chunks = append(chunks, slice[i:end])
	}
	return chunks
}

// FilterStrings 过滤字符串切片。
func FilterStrings(slice []string, fn func(string) bool) []string {
	result := make([]string, 0, len(slice))
	for _, s := range slice {
		if fn(s) {
			result = append(result, s)
		}
	}
	return result
}

// MapStrings 映射字符串切片。
func MapStrings(slice []string, fn func(string) string) []string {
	result := make([]string, len(slice))
	for i, s := range slice {
		result[i] = fn(s)
	}
	return result
}

// ReverseStrings 反转字符串切片。
func ReverseStrings(slice []string) []string {
	result := make([]string, len(slice))
	for i, s := range slice {
		result[len(slice)-1-i] = s
	}
	return result
}

// IntersectionStrings 求两个字符串切片的交集。
func IntersectionStrings(a, b []string) []string {
	set := make(map[string]bool, len(a))
	for _, s := range a {
		set[s] = true
	}
	result := make([]string, 0)
	for _, s := range b {
		if set[s] {
			result = append(result, s)
			delete(set, s)
		}
	}
	return result
}

// DifferenceStrings 求两个字符串切片的差集（在 a 中但不在 b 中）。
func DifferenceStrings(a, b []string) []string {
	set := make(map[string]bool, len(b))
	for _, s := range b {
		set[s] = true
	}
	result := make([]string, 0)
	for _, s := range a {
		if !set[s] {
			result = append(result, s)
		}
	}
	return result
}
