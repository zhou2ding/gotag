package v

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

var gViper *viper.Viper

func init() {
	v := viper.New()
	gViper = v
	err1 := LoadConfig("../../config/app.toml")
	err2 := LoadConfig("config/app.toml")
	if err1 != nil && err2 != nil {
		panic(err2)
	}
}

func GetViper() *viper.Viper {
	return gViper
}

func LoadConfig(file string) error {
	if file != "" {
		GetViper().SetConfigName("app")
		GetViper().SetConfigType("toml")
		GetViper().SetConfigFile(file)
		if err := GetViper().ReadInConfig(); err != nil {
			return errors.Wrapf(err, "Failed to load config %s", file)
		}
	}
	return nil
}
