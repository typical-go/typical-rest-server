package strkit

import "strconv"

// ToInt64 convert string to int64
func ToInt64(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}
