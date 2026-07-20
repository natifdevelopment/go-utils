package utils

import (
	"os"
	"strconv"
)

// GetEnv reads an environment variable and converts it to the requested type.
// If the variable is not set and a default value is provided, the default is
// returned. Supported types are string, int, and bool.
func GetEnv[T any](key string, defaultValue ...T) T {
	var zero T
	v := os.Getenv(key)
	if v == "" {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return zero
	}

	switch any(zero).(type) {
	case string:
		return any(v).(T)
	case int:
		i, _ := strconv.Atoi(v)
		return any(i).(T)
	case bool:
		b, _ := strconv.ParseBool(v)
		return any(b).(T)
	default:
		return zero
	}
}
