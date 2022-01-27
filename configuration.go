package configure

import (
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

var instance *interface{}

type Options struct {
	ConfigName    string
	ConfigType    string
	ConfigAbsPath string
	EnvPrefix     string
	Defaults      interface{}
}

func Setup(options Options, conf interface{}) error {
	viper.SetConfigName(options.ConfigName)
	viper.SetConfigType(options.ConfigType)
	viper.AddConfigPath(options.ConfigAbsPath)

	viper.SetEnvPrefix(options.EnvPrefix)
	viper.AutomaticEnv()

	if options.Defaults != nil {
		var defaults map[string]interface{}
		if err := mapstructure.Decode(options.Defaults, &defaults); err != nil {
			return err
		}
		for k, v := range defaults {
			viper.SetDefault(k, v)
		}
	}

	if err := viper.ReadInConfig(); err != nil {
		fileNotFoundErr, ok := err.(viper.ConfigFileNotFoundError)
		if ok {
			if err.Error() != fileNotFoundErr.Error() {
				return err
			}
		}
	}

	if err := viper.Unmarshal(&conf); err != nil {
		return err
	}

	instance = &conf

	return nil
}

func Config() interface{} {
	if instance == nil {
		return nil
	}
	return *instance
}
