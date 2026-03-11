package configure

import (
	"errors"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

type ConfigurationFrozenError string

func (err ConfigurationFrozenError) Error() string {
	return "Configuration already frozen."
}

type Configure struct {
	// The config file name. Does not include extension.
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
	// The combination of configName + configType will be used to determine the written file name.
	writeIfNotExists bool

	viper *viper.Viper
}

// SetConfigName sets the config name (does not include extension) or returns a ConfigurationFrozenError
func (c *Configure) SetConfigName(configName string) error {
	return c.checkFrozen(func() {
		c.configName = configName
	})
}

// SetConfigType sets the config type or returns a ConfigurationFrozenError
func (c *Configure) SetConfigType(ConfigType string) error {
	return c.checkFrozen(func() {
		c.configType = ConfigType
	})
}

// SetConfigDir sets the config directory or returns a ConfigurationFrozenError
func (c *Configure) SetConfigDir(ConfigDir string) error {
	return c.checkFrozen(func() {
		c.configDir = ConfigDir
	})
}

// SetEnvPrefix sets the prefix for env variable to use or returns a ConfigurationFrozenError
func (c *Configure) SetEnvPrefix(EnvPrefix string) error {
	return c.checkFrozen(func() {
		c.envPrefix = EnvPrefix
	})
}

// SetDefaults sets config defaults or returns a ConfigurationFrozenError
func (c *Configure) SetDefaults(defaults any) error {
	return c.checkFrozen(func() {
		c.defaults = defaults
	})
}

// SetWriteIfNotExists sets whether a config should be written on write or returns a ConfigurationFrozenError
func (c *Configure) SetWriteIfNotExists(WriteIfNotExists bool) error {
	return c.checkFrozen(func() {
		c.writeIfNotExists = WriteIfNotExists
	})
}

func (c *Configure) checkFrozen(cb func()) error {
	if c.viper != nil {
		return new(ConfigurationFrozenError)
	}
	cb()
	return nil
}

// New returns a new instance of Configure with defaults set.
func New() Configure {
	return Configure{
		writeIfNotExists: false,
	}
}

func (c *Configure) Get(dest any) error {
	c.setupViper()

	if c.defaults != nil {
		var defaults map[string]interface{}
		if err := mapstructure.Decode(c.defaults, &defaults); err != nil {
			return err
		}
		for k, v := range defaults {
			c.viper.SetDefault(k, v)
		}
	}

	if c.writeIfNotExists {
		saveErr := c.viper.SafeWriteConfig()
		if _, ok := errors.AsType[viper.ConfigFileAlreadyExistsError](saveErr); ok {
			// bc the option is "save if not exists", we can ignore this error
		} else {
			return saveErr
		}
	}

	readError := c.viper.ReadInConfig()
	var fileNotFoundErr viper.ConfigFileNotFoundError
	if readError != nil && errors.As(readError, &fileNotFoundErr) && readError.Error() != fileNotFoundErr.Error() {
		// return the error if the error is not a ConfigFileNotFoundError
		return readError
	}

	unmarshalErr := c.viper.Unmarshal(&dest)
	if unmarshalErr != nil {
		return unmarshalErr
	}

	return nil
}

// Save persists the given conf if the provided Configure refer to a file
func (c *Configure) Save(source any) error {
	c.setupViper()

	var dest map[string]interface{}
	if err := mapstructure.Decode(source, &dest); err != nil {
		return err
	}

	if err := c.viper.MergeConfigMap(dest); err != nil {
		return err
	}

	return c.viper.WriteConfig()
}

func (c *Configure) setupViper() {
	vpr := viper.New()

	vpr.SetConfigName(c.configName)
	vpr.SetConfigType(c.configType)
	vpr.AddConfigPath(c.configDir)
	vpr.SetEnvPrefix(c.envPrefix)
	vpr.AutomaticEnv()

	c.viper = vpr
}
