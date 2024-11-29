package options

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/DecxBase/core/types"
)

func ReadEnvs(keys ...string) types.JSONStringData {
	data := make(types.JSONStringData)

	for _, key := range keys {
		data[key] = ReadEnv("", key, "")
	}

	return data
}

func ReadEnv[T string | int | bool](prefix string, key string, def T) T {
	theKey := key
	if len(prefix) > 0 {
		theKey = fmt.Sprintf("%s_%s", prefix, key)
	}

	val, ok := os.LookupEnv(strings.ToUpper(theKey))
	if !ok {
		return def
	}

	var result any
	var err error
	switch any(def).(type) {
	case int:
		result, err = strconv.Atoi(val)
		if err != nil {
			return def
		}

	case bool:
		result, err = strconv.ParseBool(val)
		if err != nil {
			return def
		}

	default:
		result = val
	}

	return result.(T)
}
