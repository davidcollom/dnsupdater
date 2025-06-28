package api

import (
	"github.com/spf13/viper"
)

type Config struct {
	AddressProviders []ProviderConfig `mapstructure:"address_providers"`
	DNSService       DNSService       `mapstructure:"dns_service"`
}

func Load(configPath string) (*Config, error) {
	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
