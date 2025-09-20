package config

import "github.com/spf13/viper"

type Env struct {
	LogLevel         string `mapstructure:"LOG_LEVEL"`
	IdentitiesBucket string `mapstructure:"IDENTITITES_BUCKET"`
}

func MustLoadEnv() Env {
	viper.SetConfigFile(".env")

	if err := viper.ReadInConfig(); err != nil {
		// Panics because there is not logger yet
		panic(err)
	}

	var env Env

	if err := viper.Unmarshal(&env); err != nil {
		// Panics because there is not logger yet
		panic(err)
	}

	return env
}
