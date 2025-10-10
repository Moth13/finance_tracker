package util

import (
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// The vales are read by Viper from a config file or environment variables
type Config struct {
	DBDriver             string        `mapstructure:"DB_DRIVER"`
	DBSource             string        `mapstructure:"DB_SOURCE"`
	ServerAddress        string        `mapstructure:"SERVER_ADDRESS"`
	TokenSymmetricKey    string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
}

// LoadConfig reads configuration from file or environment variables
func LoadConfig(configFile string) (config Config, err error) {
	dir := filepath.Dir(configFile)
	file := filepath.Base(configFile)
	ext := filepath.Ext(configFile)
	extNoDot := strings.TrimPrefix(ext, ".")
	name := strings.TrimSuffix(file, ext)

	viper.AddConfigPath(dir)
	viper.SetConfigName(name)
	viper.SetConfigType(extNoDot)

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
