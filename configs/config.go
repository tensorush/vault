package configs

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	PostgresDSN         string        `mapstructure:"POSTGRES_DSN"`
	BotToken            string        `mapstructure:"BOT_TOKEN"`
	BotEncryptionKey    string        `mapstructure:"BOT_ENCRYPTION_KEY"`
	BotVisibilityPeriod time.Duration `mapstructure:"BOT_VISIBILITY_PERIOD"`
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigType("env")

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
