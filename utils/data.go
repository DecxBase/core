package utils

import (
	"fmt"

	"github.com/DecxBase/core/types"
)

func CheckContains[T comparable](s []T, e T) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

type DataList[T comparable] struct {
	values []T
}

func (l DataList[T]) AsList() []T {
	return l.values
}

func (l DataList[T]) Contains(c T) bool {
	for _, v := range l.values {
		if v == c {
			return true
		}
	}

	return false
}

func (l DataList[T]) Each(cb func(T, int)) {
	for i, v := range l.values {
		cb(v, i)
	}
}

func NewDataList[T comparable](values ...T) DataList[T] {
	return DataList[T]{values: values}
}

func MakeDataList[T comparable](data []T) DataList[T] {
	return DataList[T]{values: data}
}

type DataMap[T any] struct {
	values map[string]T
}

func (m DataMap[T]) AsMap() map[string]T {
	return m.values
}

func (m DataMap[T]) Contains(key string) bool {
	_, ok := m.values[key]

	return ok
}

func (m *DataMap[T]) Set(key string, value T) string {
	m.values[key] = value

	return key
}

func (m *DataMap[T]) Get(key string) T {
	return m.values[key]
}

func (m *DataMap[T]) Keys(key string) []string {
	keys := make([]string, len(m.values))

	i := 0
	for k := range m.values {
		keys[i] = k
		i++
	}

	return keys
}

func NewDataMap[T any]() DataMap[T] {
	return DataMap[T]{values: make(map[string]T)}
}

func MakeDataMap[T any](data map[string]T) DataMap[T] {
	return DataMap[T]{values: data}
}

func MapData[T comparable, K any](list DataList[T], cb func(T) K) []K {
	newData := make([]K, 0)

	list.Each(func(v T, _ int) {
		newData = append(newData, cb(v))
	})

	return newData
}

func ToJsonString(data any) types.JSONStringData {
	res := make(types.JSONStringData)

	switch data := data.(type) {
	case map[string]any:
	case types.JSONDumpData:
		for key, val := range data {
			switch val := val.(type) {
			case string:
				res[key] = val
			case int, int32, int64, float32, float64:
				res[key] = fmt.Sprintf("%v", val)
			}
		}
	case map[string]string:
	case types.JSONStringData:
		for key, val := range data {
			res[key] = val
		}
	}

	return res
}
