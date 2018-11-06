package config

import (
	"time"

	"github.com/spf13/viper"
)

// Provider defines a set of read-only methods for accessing the application
// configuration params as defined in one of the config files.
type Provider interface {
	ConfigFileUsed() string
	Get(key string) interface{}
	GetBool(key string) bool
	GetDuration(key string) time.Duration
	GetFloat64(key string) float64
	GetInt(key string) int
	GetInt64(key string) int64
	GetSizeInBytes(key string) uint
	GetString(key string) string
	GetStringMap(key string) map[string]interface{}
	GetStringMapString(key string) map[string]string
	GetStringMapStringSlice(key string) map[string][]string
	GetStringSlice(key string) []string
	GetTime(key string) time.Time
	InConfig(key string) bool
	IsSet(key string) bool
}

var defaultConfig *viper.Viper

// Config returns project config object
func Config() *viper.Viper {
	return defaultConfig
}

// LoadConfigProvider loads the config with appName
func LoadConfigProvider(appName string) Provider {
	return readViperConfig(appName)
}

func init() {
	defaultConfig = readViperConfig("FACE-RECOGNITION-BACKEND")
}

func readViperConfig(appName string) *viper.Viper {
	v := viper.New()
	v.SetEnvPrefix(appName)
	v.AutomaticEnv()

	// global defaults

	v.SetDefault("json_logs", false)
	v.SetDefault("loglevel", "debug")

	v.SetDefault("app-name", "face-recognition")

	v.SetDefault("db-addr", "mongodb://192.168.3.33")

	v.SetDefault("camera-addr", "http://192.168.0.196:8088")

	v.SetDefault("face-ai-record-addr", "http://192.168.0.196:8008/record")
	v.SetDefault("face-ai-detect-addr", "http://192.168.0.196:8008/detect")

	return v
}
