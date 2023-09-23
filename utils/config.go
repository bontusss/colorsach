package utils

import "github.com/spf13/viper"

type Config struct {
	PEXELAPIKEY string `mapstructure:"PEXELS_API_KEY"`
	PEXELSAPI string `mapstructure:"PEXELS_API"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}