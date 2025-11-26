package config


type ServEndPoints struct {
	Auth        string `mapstructure:"auth"`
	UserManager string `mapstructure:"user_manager"`
}