package env

import (
	"os"
	"strings"
)

const (
	AppEnv        = "ENV"
	ConfigFileEnv = "CONFIG_FILE"
)

var Prefix = "KIT_"

func GetEnv(key string) string {
	return os.Getenv(Prefix + strings.ToUpper(key))
}

func GetEnvDefault(key, d string) string {
	if v := GetEnv(key); v != "" {
		return v
	}
	return d
}
