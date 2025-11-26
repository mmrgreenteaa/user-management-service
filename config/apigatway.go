package config

import (
	"fmt"
	"log"

	api "github.com/mmrgreenteaa/user-management-service/internal/api_gatway/handlers"
	"github.com/spf13/viper"
)



func GetApiGatwayConfg() (*api.ApiGatwayConfig, error) {

	type Config struct {
		CfgApi api.ApiGatwayConfig `mapstructure:"api_gateway"`
		CfgEnd ServEndPoints       `mapstructure:"services_endpoints"`
	}
	cfg := Config{
		CfgApi: api.ApiGatwayConfig{},
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("../../../config")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, fmt.Errorf("failed configuration not found")
		} else {
			return nil, fmt.Errorf("failed onfiguration load: %w", err)
		}
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("falied unmarshal - %w", err)
	}

	cfg.CfgApi.AuthIp = cfg.CfgEnd.Auth
	cfg.CfgApi.UserMengIP = cfg.CfgEnd.UserManager
	log.Println(cfg)
	return &cfg.CfgApi, nil

}
