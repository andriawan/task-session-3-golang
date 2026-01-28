package config

import (
	"fmt"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

func Load() (*Template, error) {
	_ = godotenv.Load()
	v := viper.New()
	v.SetConfigName("application")
	v.SetConfigType("yaml")
	v.AddConfigPath("./config")
	v.Debug()

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	var template Template
	if err := v.Unmarshal(&template); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	return &template, nil
}
