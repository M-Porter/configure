package configure_test

import (
	"os"
	"testing"

	"github.com/m-porter/configure/v3"
)

func TestDefaults(t *testing.T) {
	type TestingConfig struct {
		HostName string `mapstructure:"host_name"`
	}

	expected := "my.test.host"

	testingConfig := &TestingConfig{}

	conf := configure.New()
	conf.Defaults = TestingConfig{HostName: expected}

	err := conf.Get(testingConfig)
	if err != nil {
		t.Fatalf("error setting up config: %v", err)
	}

	if testingConfig.HostName != expected {
		t.Errorf("expected %s, got %s", expected, testingConfig.HostName)
	}
}

func TestFromEnvUsingDefaultPrefix(t *testing.T) {
	type TestingConfig struct {
		Secret string `mapstructure:"app_secret"`
	}

	key := "APP_SECRET"
	expected := "987xyz"
	defer setEnvValue(t, key, expected)()

	testingConfig := &TestingConfig{}

	conf := configure.New()
	conf.Defaults = TestingConfig{Secret: "abc123"}

	err := conf.Get(testingConfig)

	if err != nil {
		t.Fatalf("error setting up config: %v", err)
	}

	if testingConfig.Secret != expected {
		t.Errorf("expected %s, got %s", expected, testingConfig.Secret)
	}

	_ = os.Unsetenv(key)
}

func TestFromEnvWithPrefix(t *testing.T) {
	type TestingConfig struct {
		Secret string `mapstructure:"secret"`
	}

	key := "FOO_SECRET"
	expected := "987xyz"
	defer setEnvValue(t, key, expected)()

	testingConfig := &TestingConfig{}

	conf := configure.New()
	conf.Defaults = TestingConfig{Secret: "abc123"}
	conf.EnvPrefix = "foo"

	err := conf.Get(testingConfig)

	if err != nil {
		t.Fatalf("error setting up config: %v", err)
	}

	if testingConfig.Secret != expected {
		t.Errorf("expected %s, got %s", expected, testingConfig.Secret)
	}
}

func TestFromFile(t *testing.T) {
	type TestingConfig struct {
		SomeValue string `mapstructure:"some_value"`
	}

	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Unable to get working dir: %v", err)
	}

	testingConfig := &TestingConfig{}

	conf := configure.New()
	conf.ConfigName = "test_config"
	conf.ConfigType = "yaml"
	conf.ConfigDir = dir

	err = conf.Get(testingConfig)

	if err != nil {
		t.Fatalf("error setting up config: %v", err)
	}

	expected := "foo foo foo"
	if testingConfig.SomeValue != expected {
		t.Errorf("expected %s, got %s", expected, testingConfig.SomeValue)
	}
}

func TestFromFileUsingConfigNameOnly(t *testing.T) {
	type TestingConfig struct {
		SomeValue string `mapstructure:"some_value"`
	}

	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Unable to get working dir: %v", err)
	}

	testingConfig := &TestingConfig{}

	conf := configure.New()
	conf.ConfigName = "test_config.yaml"
	conf.ConfigDir = dir

	err = conf.Get(testingConfig)

	if err != nil {
		t.Fatalf("error setting up config: %v", err)
	}

	expected := "foo foo foo"
	if testingConfig.SomeValue != expected {
		t.Errorf("expected %s, got %s", expected, testingConfig.SomeValue)
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
