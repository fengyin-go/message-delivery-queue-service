package utils

import (
	"strconv"
	"strings"
	"unicode"
)

// Trim 去除字符串两端空白字符。
func Trim(s string) string {
	return strings.TrimSpace(s)
}

// IsEmpty 判断字符串是否为空或仅包含空白字符。
func IsEmpty(s string) bool {
	return strings.TrimSpace(s) == ""
}

// IsNotEmpty 判断字符串是否非空。
func IsNotEmpty(s string) bool {
	return !IsEmpty(s)
}

// DefaultIfEmpty 如果字符串为空则返回默认值。
func DefaultIfEmpty(s, def string) string {
	if IsEmpty(s) {
		return def
	}
	return s
}

// ContainsIgnoreCase 判断字符串是否包含子串（忽略大小写）。
func ContainsIgnoreCase(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}

// ToInt 将字符串转为整数，失败返回 0。
func ToInt(s string) int {
	n, _ := strconv.Atoi(s)
	return n
}

// ToInt64 将字符串转为 int64，失败返回 0。
func ToInt64(s string) int64 {
	n, _ := strconv.ParseInt(s, 10, 64)
	return n
}

// ToBool 将字符串转为布尔值。
func ToBool(s string) bool {
	b, _ := strconv.ParseBool(s)
	return b
}

// SnakeToCamel 将 snake_case 转为 camelCase。
func SnakeToCamel(s string) string {
	parts := strings.Split(s, "_")
	for i := range parts {
		if len(parts[i]) > 0 {
			parts[i] = strings.ToUpper(parts[i][:1]) + parts[i][1:]
		}
	}
	return strings.Join(parts, "")
}

// Truncate 截断字符串到指定长度，超出部分用省略号代替。
func Truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}

// CountWords 统计字符串中的单词数（简单空白分隔）。
func CountWords(s string) int {
	fields := strings.Fields(s)
	return len(fields)
}

// RemoveNonPrintable 去除不可打印字符。
func RemoveNonPrintable(s string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsPrint(r) {
			return r
		}
		return -1
	}, s)
}

// PadLeft 在字符串左侧填充字符至指定长度。
func PadLeft(s string, length int, pad rune) string {
	if len(s) >= length {
		return s
	}
	padding := strings.Repeat(string(pad), length-len(s))
	return padding + s
}

// PadRight 在字符串右侧填充字符至指定长度。
func PadRight(s string, length int, pad rune) string {
	if len(s) >= length {
		return s
	}
	padding := strings.Repeat(string(pad), length-len(s))
	return s + padding
}

// Reverse 反转字符串。
func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
