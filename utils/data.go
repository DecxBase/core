package utils

import "strconv"

func StringToIntBase(v string, bitSize int) any {
	val, err := strconv.ParseInt(v, 10, bitSize)

	if err == nil {
		return val
	}

	return 0
}

func StringToInt32(v string) int32 {
	return StringToIntBase(v, 32).(int32)
}

func StringToInt64(v string) int64 {
	return StringToIntBase(v, 64).(int64)
}
