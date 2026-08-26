// Package utils 提供通用工具函数。
package utils

import "time"

// FormatTime 将时间格式化为本地字符串。
func FormatTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

// FormatDate 将时间格式化为日期字符串。
func FormatDate(t time.Time) string {
	return t.Format("2006-01-02")
}

// IsExpired 判断给定时间是否已经过期（基于当前时间）。
func IsExpired(t time.Time, ttl time.Duration) bool {
	return time.Since(t) > ttl
}

// TruncateToSecond 将时间截断到秒。
func TruncateToSecond(t time.Time) time.Time {
	return t.Truncate(time.Second)
}

// StartOfDay 返回当天的开始时间。
func StartOfDay(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, t.Location())
}

// EndOfDay 返回当天的结束时间。
func EndOfDay(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 23, 59, 59, 999999999, t.Location())
}

// DaysBetween 计算两个时间之间的天数差。
func DaysBetween(a, b time.Time) int {
	if a.After(b) {
		a, b = b, a
	}
	return int(b.Sub(a).Hours() / 24)
}

// AddDays 给时间增加指定天数。
func AddDays(t time.Time, days int) time.Time {
	return t.AddDate(0, 0, days)
}

// TimeEqual 判断两个时间是否相等（精确到秒）。
func TimeEqual(a, b time.Time) bool {
	return a.Truncate(time.Second).Equal(b.Truncate(time.Second))
}
