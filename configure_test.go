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

func TestFromEnvUsingDefaultPrefix(t *testing.T) {
	type Config struct {
		Secret string `mapstructure:"app_secret"`
	}

	key := "APP_SECRET"
	expected := "987xyz"
	defer setEnvValue(t, key, expected)()

	conf := &Config{}

	err := configure.Get(
		conf,
		configure.WithDefaultConfig(Config{
			Secret: "abc123",
		}),
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
		Secret string `mapstructure:"secret"`
	}

	key := "FOO_SECRET"
	expected := "987xyz"
	defer setEnvValue(t, key, expected)()

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

func setEnvValue(t *testing.T, key, value string) func() {
	if err := os.Unsetenv(key); err != nil {
		t.Fatalf("error unsetting env var for testing: %v", err)
	}
	if err := os.Setenv(key, value); err != nil {
		t.Fatalf("error unsetting env var for testing: %v", err)
	}
	actual := os.Getenv(key)
	if os.Getenv(key) != value {
		t.Fatalf("error unsetting env var for testing: expected %s, got %s", value, actual)
	}

	return func() {
		if err := os.Unsetenv(key); err != nil {
			t.Fatalf("error unsetting env var for testing: %v", err)
		}
	}
}
