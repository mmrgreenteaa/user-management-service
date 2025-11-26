package config

import (
	"fmt"
	"log"

	auth "github.com/mmrgreenteaa/user-management-service/internal/auth/handlers"
	"github.com/spf13/viper"
)

func GetAuth() (*auth.AuthConfig, error) {
	type Config struct {
		AuthCfg auth.AuthConfig `mapstructure:"auth"`
		СfgEnd  ServEndPoints   `mapstructure:"services_endpoints"`
	}

	cfg := Config{
		AuthCfg: auth.AuthConfig{},
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("../../../config")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, fmt.Errorf("failed configuration not found %w", err)
		} else {
			return nil, fmt.Errorf("failed configuration load: %w", err)
		}
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed unmarshal : %w", err)
	}
	viper.AutomaticEnv()
	viper.BindEnv("secretKey", "AUTH_SECRET_KEY")
	viper.BindEnv("passDb", "AUTH_PASSWORD_DB")
	cfg.AuthCfg.SecretKey = viper.GetString("secretKey")
	cfg.AuthCfg.DbConfig.Pass = viper.GetString("passDb")
	

	cfg.AuthCfg.UserMengIp = cfg.СfgEnd.UserManager

	log.Println(cfg)
	return &cfg.AuthCfg, nil

}
