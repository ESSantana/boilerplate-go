package config

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/viper"
)

func Load() (*Config, error) {
	viper.AddConfigPath("/app/config")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

  viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	var config Config

	if err := viper.Unmarshal(&config.Server); err != nil {
		return nil, err
	}
	if err := viper.Unmarshal(&config.Database); err != nil {
		return nil, err
	}
	if err := viper.Unmarshal(&config.Redis); err != nil {
		return nil, err
	}
	if err := viper.Unmarshal(&config.Google); err != nil {
		return nil, err
	}
	if err := viper.Unmarshal(&config.MercadoPago); err != nil {
		return nil, err
	}
	if err := viper.Unmarshal(&config.JWT); err != nil {
		return nil, err
	}
	if err := viper.Unmarshal(&config.AWS); err != nil {
		return nil, err
	}
	if err := viper.Unmarshal(&config.Frontend); err != nil {
		return nil, err
	}

	configJSON, _ := json.Marshal(config)
	fmt.Println("Configuration loaded successfully:", string(configJSON))

	return &config, nil
}
