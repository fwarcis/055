// Package config...
package config

import "os"

func EnvWithDefault(key string, byDefault string) string {
	value, found := os.LookupEnv(key)
	if !found {
		return byDefault
	}
	return value
}
