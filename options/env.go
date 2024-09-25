package options

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func ReadEnv[T string | int | bool](prefix string, key string, def T) T {
	val, ok := os.LookupEnv(strings.ToUpper(fmt.Sprintf("%s_%s", prefix, key)))

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
