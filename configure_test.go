package configure_test

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path"
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
	conf.SetDefaults(TestingConfig{HostName: expected})

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
	conf.SetDefaults(TestingConfig{Secret: "abc123"})

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
	conf.SetDefaults(TestingConfig{Secret: "abc123"})
	conf.SetEnvPrefix("foo")

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
	conf.SetConfigName("test_config")
	conf.SetConfigType("yaml")
	conf.SetConfigDir(dir)

	err = conf.Get(testingConfig)
	if err != nil {
		t.Fatalf("error setting up config: %v", err)
	}

	expected := "foo foo foo"
	if testingConfig.SomeValue != expected {
		t.Errorf("expected %s, got %s", expected, testingConfig.SomeValue)
	}
}

func TestWriteIfNotExists(t *testing.T) {
	type TestingConfig struct {
		SomeValue string `mapstructure:"some_value"`
	}

	testingConfig := &TestingConfig{}

	tempDir := os.TempDir()
	tempFile := fmt.Sprintf("configure_tests_%08x", rand.Uint32())
	fullTempFilePath := path.Join(tempDir, tempFile+".yaml")
	defer os.Remove(fullTempFilePath)

	conf := configure.New()
	conf.SetConfigName(tempFile)
	conf.SetConfigType("yaml")
	conf.SetConfigDir(tempDir)
	conf.SetDefaults(TestingConfig{SomeValue: "abc123"})
	conf.SetWriteIfNotExists(true)

	if err := conf.Get(testingConfig); err != nil {
		t.Fatalf("error setting up config: %v", err)
	}

	if _, err := os.Stat(fullTempFilePath); os.IsNotExist(err) {
		t.Errorf("file %s not created", fullTempFilePath)
	}
}

func TestCanSaveAConfig(t *testing.T) {
	type TestingConfig struct {
		SomeValue string `mapstructure:"some_value"`
	}

	testingConfig := &TestingConfig{}

	tempDir := os.TempDir()
	tempFile := fmt.Sprintf("configure_tests_%08x", rand.Uint32())
	fullTempFilePath := path.Join(tempDir, tempFile+".yaml")
	defer os.Remove(fullTempFilePath)

	conf := configure.New()
	conf.SetConfigName(tempFile)
	conf.SetConfigType("yaml")
	conf.SetConfigDir(tempDir)
	conf.SetDefaults(TestingConfig{SomeValue: fmt.Sprintf("%08x", rand.Uint32())})
	conf.SetWriteIfNotExists(true)

	if err := conf.Get(testingConfig); err != nil {
		t.Fatalf("error setting up config: %v", err)
	}

	expected := fmt.Sprintf("%08x", rand.Uint32())
	testingConfig.SomeValue = expected

	if err := conf.Save(testingConfig); err != nil {
		t.Fatalf("error saving config: %v", err)
	}

	// next try to read in the file again with a new instance. it should now be the expected value
	testingConfig2 := &TestingConfig{}

	conf2 := configure.New()
	conf2.SetConfigName(tempFile)
	conf2.SetConfigType("yaml")
	conf2.SetConfigDir(tempDir)

	if err := conf2.Get(testingConfig2); err != nil {
		t.Fatalf("error setting up config: %v", err)
	}

	if testingConfig2.SomeValue != expected {
		t.Errorf("expected %s, got %s", expected, testingConfig2.SomeValue)
	}
}

func TestFrozen(t *testing.T) {
	type TestingConfig struct {
		SomeValue string `mapstructure:"some_value"`
	}

	conf := configure.New()

	if err := conf.Get(&TestingConfig{}); err != nil {
		t.Fatalf("error initializing config: %v", err)
	}

	setters := []struct {
		name string
		fn   func() error
	}{
		{"SetConfigName", func() error { return conf.SetConfigName("config") }},
		{"SetConfigType", func() error { return conf.SetConfigType("yaml") }},
		{"SetConfigDir", func() error { return conf.SetConfigDir("/etc/myapp") }},
		{"SetEnvPrefix", func() error { return conf.SetEnvPrefix("myapp") }},
		{"SetDefaults", func() error { return conf.SetDefaults(TestingConfig{}) }},
		{"SetWriteIfNotExists", func() error { return conf.SetWriteIfNotExists(true) }},
	}

	for _, s := range setters {
		err := s.fn()
		if err == nil {
			t.Errorf("%s: expected ConfigurationFrozenError, got nil", s.name)
			continue
		}
		if _, ok := errors.AsType[*configure.ConfigurationFrozenError](err); !ok {
			t.Errorf("%s: expected ConfigurationFrozenError, got %T: %v", s.name, err, err)
		}
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
