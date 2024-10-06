package utils

import (
	"encoding/json"
	"net/url"
)

func QueryToMap(source url.Values) map[string]string {
	values := make(map[string]string)

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

func StructToMap(obj interface{}) (map[string]interface{}, error) {
	var result map[string]interface{}

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
