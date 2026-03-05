package configure

import (
	"errors"
	"path/filepath"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

type setupOptions struct {
	configName    string
	configType    string
	configAbsPath string
	envPrefix     string
	// Default values to apply
	defaults interface{}
	// If no config exists at the given config file path, should one be written.
	writeIfNotExists bool
}

type SetupOption func(o *setupOptions)

func WithConfigName(name string) SetupOption {
	return func(o *setupOptions) {
		o.configName = name
	}
}

func WithConfigDir(dir string) SetupOption {
	return func(o *setupOptions) {
		o.configAbsPath = dir
	}
}

func WithConfigType(name string) SetupOption {
	return func(o *setupOptions) {
		o.configType = name
	}
}

func WithEnvPrefix(prefix string) SetupOption {
	return func(o *setupOptions) {
		o.envPrefix = strings.ToUpper(prefix)
	}
}

func WithWriteIfNotExists() SetupOption {
	return func(o *setupOptions) {
		o.writeIfNotExists = true
	}
}

func WithDefaultConfig(defaults any) SetupOption {
	return func(o *setupOptions) {
		o.defaults = defaults
	}
}

func Setup(conf any, options ...SetupOption) error {
	opts := setupOptions{
		writeIfNotExists: false,
	}
	for _, option := range options {
		option(&opts)
	}

	vpr := viper.New()

	vpr.SetConfigName(opts.configName)

	// If config type is not set, attempt to assume it from the configName
	if opts.configType == "" {
		ext := filepath.Ext(opts.configName)
		if len(ext) > 1 {
			vpr.SetConfigType(ext[1:])
		}
	} else {
		vpr.SetConfigType(opts.configType)
	}

	vpr.AddConfigPath(opts.configAbsPath)

	vpr.SetEnvPrefix(opts.envPrefix)
	vpr.AutomaticEnv()

	if opts.defaults != nil {
		var defaults map[string]interface{}
		if err := mapstructure.Decode(opts.defaults, &defaults); err != nil {
			return err
		}
		for k, v := range defaults {
			vpr.SetDefault(k, v)
		}
	}

	if opts.writeIfNotExists {
		saveErr := vpr.SafeWriteConfig()
		var configFileAlreadyExistsError viper.ConfigFileAlreadyExistsError
		if errors.As(saveErr, &configFileAlreadyExistsError) {
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
