package configs

import (
	"github.com/spf13/viper"
)

func LoadDBconfig(configFile string) (dbConfig DBconfig, err error) {
	viper.SetConfigFile(configFile)

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&dbConfig)
	return
}

func LoadAppconfig(configFile string) (appConfig ApplicationConfig, err error) {
	viper.SetConfigFile(configFile)

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&appConfig)
	return
}
