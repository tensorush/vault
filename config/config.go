package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Type string `mapstructure:"DB_TYPE"`
	Data string `mapstructure:"DB_DATA"`
	Name string `mapstructure:"DB_NAME"`
	User string `mapstructure:"DB_USER"`
	Pass string `mapstructure:"DB_PASS"`

	Token            string        `mapstructure:"BOT_TOKEN"`
	EncryptionKey    string        `mapstructure:"BOT_ENC_KEY"`
	ExpirationPeriod time.Duration `mapstructure:"BOT_EXP_PERIOD"`
}

func LoadConfig() (*Config, error) {
	var err error
	var config Config

	viper.SetConfigFile(".env")
	err = viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
