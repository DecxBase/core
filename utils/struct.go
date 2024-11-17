package utils

import (
	"encoding/json"
	"net/url"

	"github.com/DecxBase/core/types"
)

func QueryToMap(source url.Values) types.JSONStringData {
	values := make(types.JSONStringData)

	for key := range source {
		values[key] = source.Get(key)
	}

	return values
}

func DataToBytes(source any) ([]byte, error) {
	return json.Marshal(source)
}

func BytesToData(bytes []byte, obj any) error {
	return json.Unmarshal(bytes, obj)
}

func MapToStruct(source map[string]any, obj any) error {
	bytes, err := DataToBytes(source)
	if err != nil {
		return err
	}

	err = BytesToData(bytes, obj)
	if err != nil {
		return err
	}

	return nil
}

func StructToMap(obj interface{}) (types.JSONDumpData, error) {
	var result types.JSONDumpData

	jsonBytes, err := DataToBytes(obj)
	if err != nil {
		return nil, err
	}

	err = BytesToData(jsonBytes, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
