package config

import (
	"fmt"
	"github.com/spf13/viper"
)

func NewConfig(env string) (*viper.Viper, error) {
	conf := viper.New()
	conf.AddConfigPath("./config")
	conf.SetConfigName(env)

	if err := conf.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %v\n", err)
	}

	return conf, nil
}
