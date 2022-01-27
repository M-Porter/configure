package main

import (
	"fmt"
	"github.com/m-porter/configure"
	"os"
	"path"
	"strings"
)

type Config struct {
	Env string `mapstructure:"env"`
}

func main() {
	err := os.Setenv("APP_ENV", "staging")
	if err != nil {
		panic(err)
	}

	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	if !strings.Contains(dir, "examples/from-env") {
		dir = path.Join(dir, "examples", "from-env")
	}

	options := configure.Options{
		EnvPrefix: "APP",
		Defaults:  Config{Env: "production"},
	}

	if err := configure.Setup(options, Config{}); err != nil {
		panic(err)
	}

	conf := configure.Config().(Config)

	fmt.Printf("should output staging: %v\n", conf.Env)
}

func init() {
	_ = os.Setenv("APP_ENV", "staging")
}
