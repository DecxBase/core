package utils

import (
	"strconv"

	"github.com/DecxBase/core/types"
)

func StringToIntBase(v string, bitSize int) any {
	val, err := strconv.ParseInt(v, 10, bitSize)

	if err == nil {
		return val
	}

	return 0
}

func StringToInt(v string) int {
	return StringToIntBase(v, 32).(int)
}

func StringToInt32(v string) int32 {
	return StringToIntBase(v, 32).(int32)
}

func StringToInt64(v string) int64 {
	return StringToIntBase(v, 64).(int64)
}

func StringToUInt32(v string) uint32 {
	return StringToIntBase(v, 32).(uint32)
}

func StringToUInt64(v string) uint64 {
	return StringToIntBase(v, 64).(uint64)
}

func IsType[T any](data any) bool {
	_, ok := data.(T)
	return ok
}

func PluckFields(data types.JSONDumpData, fields []string) types.JSONDumpData {
	result := make(types.JSONDumpData)

	for _, key := range fields {
		result[key] = data[key]
	}

	return result
}
