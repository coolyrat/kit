package env

import (
	"os"
	"strings"
)

const Prefix = "ENV_PREFIX"

// env with prefix
const (
	AppEnv        = "ENV"
	ConfigFileEnv = "CONFIG_FILE"
)

var prefix = prefixEnv()

func GetEnv(key string) string {
	return os.Getenv(prefix + strings.ToUpper(key))
}

func GetEnvDefault(key, d string) string {
	if v := GetEnv(key); v != "" {
		return v
	}
	return d
}

func prefixEnv() string {
	if p := os.Getenv(Prefix); p != "" {
		return p
	}

	return "KIT_"
}
