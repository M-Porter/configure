package configure

import (
	"errors"
	"path/filepath"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

type Configure struct {
	// The config file name.
	ConfigName string
	// The config type. If not supplied, it will be assumed from the ConfigName extension.
	ConfigType string
	// The absolute path to the directory the config is stored in.
	ConfigDir string
	// Defines a prefix that env variables will use.
	EnvPrefix string
	// Default values to apply
	Defaults interface{}
	// If no config exists at the given config file path, should one be written.
	WriteIfNotExists bool
}

func New() Configure {
	return Configure{
		WriteIfNotExists: false,
	}
}

func (configure *Configure) Get(conf any) error {
	vpr := viper.New()

	vpr.SetConfigName(configure.ConfigName)

	// If config type is not set, attempt to assume it from the ConfigName
	if configure.ConfigType == "" {
		ext := filepath.Ext(configure.ConfigName)
		if len(ext) > 1 {
			vpr.SetConfigType(ext[1:])
		}
	} else {
		vpr.SetConfigType(configure.ConfigType)
	}

	vpr.AddConfigPath(configure.ConfigDir)

	vpr.SetEnvPrefix(configure.EnvPrefix)
	vpr.AutomaticEnv()

	if configure.Defaults != nil {
		var defaults map[string]interface{}
		if err := mapstructure.Decode(configure.Defaults, &defaults); err != nil {
			return err
		}
		for k, v := range defaults {
			vpr.SetDefault(k, v)
		}
	}

	if configure.WriteIfNotExists {
		saveErr := vpr.SafeWriteConfig()
		if _, ok := errors.AsType[viper.ConfigFileAlreadyExistsError](saveErr); ok {
			// bc the option is "save if not exists", we can ignore this error
		} else {
			return saveErr
		}
	}

	readError := vpr.ReadInConfig()
	var fileNotFoundErr viper.ConfigFileNotFoundError
	if readError != nil && errors.As(readError, &fileNotFoundErr) && readError.Error() != fileNotFoundErr.Error() {
		// return the error if the error is not a ConfigFileNotFoundError
		return readError
	}

	unmarshalErr := vpr.Unmarshal(&conf)
	if unmarshalErr != nil {
		return unmarshalErr
	}

	return nil
}

// Save persists the given conf if the provided Configure refer to a file
//func Save(conf any, options *Configure) error {
//	actualOptions := mergeOptions(options)
//
//	return nil
//}
