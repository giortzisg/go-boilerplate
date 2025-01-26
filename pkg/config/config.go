package config

import (
	"fmt"

	"github.com/spf13/viper"
)

func NewConfig(path string) (*viper.Viper, error) {
	conf := viper.New()
	conf.SetConfigFile(path)

	if err := conf.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %v\n and path: %s", err, path)
	}

	return conf, nil
}
