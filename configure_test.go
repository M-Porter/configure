package configure_test

import (
	"os"
	"testing"

	"github.com/m-porter/configure/v2"
)

func TestDefaults(t *testing.T) {
	type Config struct {
		HostName string `mapstructure:"host_name"`
	}

	expected := "my.test.host"

	options := configure.Options{
		Defaults: Config{
			HostName: expected,
		},
	}

	conf := &Config{}

	if err := configure.Setup(options, conf); err != nil {
		t.Fatalf("error setting up config: %v", err)
	}

	if conf.HostName != expected {
		t.Errorf("expected %s, got %s", expected, conf.HostName)
	}
}

func TestFromEnv(t *testing.T) {
	type Config struct {
		Secret string `mapstructure:"app_secret"`
	}

	key := "APP_SECRET"
	_ = os.Unsetenv(key)
	expected := "987xyz"
	_ = os.Setenv(key, expected)

	options := configure.Options{
		EnvPrefix: "",
		Defaults: Config{
			Secret: "abc123",
		},
	}

	conf := &Config{}

	if err := configure.Setup(options, conf); err != nil {
		t.Fatalf("error setting up config: %v", err)
	}

	if conf.Secret != expected {
		t.Errorf("expected %s, got %s", expected, conf.Secret)
	}

	_ = os.Unsetenv(key)
}

func TestFromEnvWithNoEnvPrefix(t *testing.T) {
	type Config struct {
		AnotherSecret string `mapstructure:"another_secret"`
	}

	key := "APP_ANOTHER_SECRET"
	_ = os.Unsetenv(key)
	expected := "987xyz"
	_ = os.Setenv(key, expected)

	options := configure.Options{
		EnvPrefix: "APP",
		Defaults: Config{
			AnotherSecret: "abc123",
		},
	}

	conf := &Config{}

	if err := configure.Setup(options, conf); err != nil {
		t.Fatalf("error setting up config: %v", err)
	}

	if conf.AnotherSecret != expected {
		t.Errorf("expected %s, got %s", expected, conf.AnotherSecret)
	}

	_ = os.Unsetenv(key)
}

func TestFromFile(t *testing.T) {
	type Config struct {
		SomeValue string `mapstructure:"some_value"`
	}

	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Unable to get working dir: %v", err)
	}

	options := configure.Options{
		ConfigName:    "test_config",
		ConfigType:    "yaml",
		ConfigAbsPath: dir,
	}

	conf := &Config{}

	if err := configure.Setup(options, conf); err != nil {
		t.Fatalf("error setting up config: %v", err)
	}

	expected := "foo foo foo"
	if conf.SomeValue != expected {
		t.Errorf("expected %s, got %s", expected, conf.SomeValue)
	}
}
