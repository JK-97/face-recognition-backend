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

func Config() *viper.Viper {
	return defaultConfig
}

func LoadConfigProvider(appName string) Provider {
	return readViperConfig(appName)
}

func init() {
	defaultConfig = readViperConfig("TF-FENCE-BACKEND")
}

func readViperConfig(appName string) *viper.Viper {
	v := viper.New()
	v.SetEnvPrefix(appName)
	v.AutomaticEnv()

	// global defaults

	v.SetDefault("json_logs", false)
	v.SetDefault("loglevel", "debug")

	v.SetDefault("app-name", "tf-fence")

	v.SetDefault("data-in-addr", "192.168.3.33:6379")
	v.SetDefault("data-in-chan", "tf-fence")

	v.SetDefault("event-out-addr", "192.168.3.33:6379")
	v.SetDefault("event-out-chan", "edge_dashboard_events")

	// fence: ymin, xmin, ymax, xmax
	// no fence by default
	v.SetDefault("fence-position", [4]float32{0, 0, 1, 1})

	return v
}
