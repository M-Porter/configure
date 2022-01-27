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
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	if !strings.Contains(dir, "examples/from-file") {
		dir = path.Join(dir, "examples", "from-file")
	}

	options := configure.Options{
		ConfigName:    "config",
		ConfigType:    "yaml",
		ConfigAbsPath: dir,
	}

	if err := configure.Setup(options, Config{}); err != nil {
		panic(err)
	}

	conf := configure.Config().(Config)

	// should output development
	fmt.Printf("should output development: %v\n", conf.Env)
}
