package env

import (
	"os"
	"strings"
)

const (
	AppEnv        = "APPLICATION_ENV"
	ConfigFileEnv = "CONFIG_FILE"
)

func GetEnv(key string) string {
	return os.Getenv(strings.ToUpper(key))
}

func GetEnvDefault(key, d string) string {
	if v := GetEnv(key); v != "" {
		return v
	}
	return d
}
