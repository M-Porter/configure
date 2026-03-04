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

	conf := &Config{}

	err := configure.Get(
		conf,
		configure.WithDefaultConfig(Config{
			HostName: expected,
		}),
	)
	if err != nil {
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

	conf := &Config{}

	err := configure.Get(
		conf,
		configure.WithDefaultConfig(Config{
			Secret: "abc123",
		}),
		configure.WithEnvPrefix(""),
	)

	if err != nil {
		t.Fatalf("error setting up config: %v", err)
	}

	if conf.Secret != expected {
		t.Errorf("expected %s, got %s", expected, conf.Secret)
	}

	_ = os.Unsetenv(key)
}

func TestFromEnvWithPrefix(t *testing.T) {
	type Config struct {
		Secret string `mapstructure:"foo_secret"`
	}

	key := "FOO_SECRET"
	_ = os.Unsetenv(key)
	expected := "987xyz"
	_ = os.Setenv(key, expected)

	conf := &Config{}

	err := configure.Get(
		conf,
		configure.WithDefaultConfig(Config{
			Secret: "abc123",
		}),
		configure.WithEnvPrefix("foo"),
	)

	if err != nil {
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

	conf := &Config{}

	err := configure.Get(
		conf,
		configure.WithDefaultConfig(Config{
			AnotherSecret: "abc123",
		}),
		configure.WithEnvPrefix("APP"),
	)

	if err != nil {
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

	conf := &Config{}

	err = configure.Get(
		conf,
		configure.WithConfigName("test_config"),
		configure.WithConfigType("yaml"),
		configure.WithConfigDir(dir),
	)

	if err != nil {
		t.Fatalf("error setting up config: %v", err)
	}

	expected := "foo foo foo"
	if conf.SomeValue != expected {
		t.Errorf("expected %s, got %s", expected, conf.SomeValue)
	}
}

func TestFromFileUsingConfigNameOnly(t *testing.T) {
	type Config struct {
		SomeValue string `mapstructure:"some_value"`
	}

	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Unable to get working dir: %v", err)
	}

	conf := &Config{}

	err = configure.Get(
		conf,
		configure.WithConfigName("test_config.yaml"),
		configure.WithConfigDir(dir),
	)

	if err != nil {
		t.Fatalf("error setting up config: %v", err)
	}

	expected := "foo foo foo"
	if conf.SomeValue != expected {
		t.Errorf("expected %s, got %s", expected, conf.SomeValue)
	}
}
