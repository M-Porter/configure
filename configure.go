package configure

import (
	"errors"
	"path/filepath"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

type Configure struct {
	// The config file name.
	configName string
	// The config type. If not supplied, it will be assumed from the configName extension.
	configType string
	// The absolute path to the directory the config is stored in.
	configDir string
	// Defines a prefix that env variables will use.
	envPrefix string
	// Default values to apply
	defaults interface{}
	// If no config exists at the given config file path, should one be written.
	writeIfNotExists bool
}

func (c *Configure) SetConfigName(configName string) {
	c.configName = configName
}

func (c *Configure) SetConfigType(ConfigType string) {
	c.configType = ConfigType
}

func (c *Configure) SetConfigDir(ConfigDir string) {
	c.configDir = ConfigDir
}

func (c *Configure) SetEnvPrefix(EnvPrefix string) {
	c.envPrefix = EnvPrefix
}

func (c *Configure) SetDefaults(defaults any) {
	c.defaults = defaults
}

func (c *Configure) SetWriteIfNotExists(WriteIfNotExists bool) {
	c.writeIfNotExists = WriteIfNotExists
}

// New returns a new instance of Configure with defaults set.
func New() Configure {
	return Configure{
		writeIfNotExists: false,
	}
}

func (c *Configure) Get(dest any) error {
	vpr := viper.New()

	vpr.SetConfigName(c.configName)

	// If config type is not set, attempt to assume it from the configName
	if c.configType == "" {
		ext := filepath.Ext(c.configName)
		if len(ext) > 1 {
			vpr.SetConfigType(ext[1:])
		}
	} else {
		vpr.SetConfigType(c.configType)
	}

	vpr.AddConfigPath(c.configDir)

	vpr.SetEnvPrefix(c.envPrefix)
	vpr.AutomaticEnv()

	if c.defaults != nil {
		var defaults map[string]interface{}
		if err := mapstructure.Decode(c.defaults, &defaults); err != nil {
			return err
		}
		for k, v := range defaults {
			vpr.SetDefault(k, v)
		}
	}

	if c.writeIfNotExists {
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

	unmarshalErr := vpr.Unmarshal(&dest)
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
