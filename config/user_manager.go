package config

import (
	"fmt"
	"log"

	usm "github.com/mmrgreenteaa/user-management-service/internal/user_management/handlers"
	"github.com/spf13/viper"
)

func GetUsrMnger() (*usm.UserMengConfig, error) {

	type Config struct {
		UserMengr usm.UserMengConfig `mapstructure:"manager"`
	}

	cfg := Config{
		UserMengr: usm.UserMengConfig{},
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
	viper.BindEnv("passDb", "USER_MANADGER_PASS_DB")
	cfg.UserMengr.DB.Pass = viper.GetString("passDb")
	log.Println(cfg)
	return &cfg.UserMengr, nil

}
