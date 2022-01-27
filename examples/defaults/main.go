package main

import (
	"fmt"
	"github.com/m-porter/configure"
)

type Config struct {
	Env string `mapstructure:"env"`
}

func main() {
	options := configure.Options{
		Defaults: Config{Env: "production"},
	}

	if err := configure.Setup(options, Config{}); err != nil {
		panic(err)
	}

	conf := configure.Config().(Config)

	fmt.Printf("should output production: %v\n", conf.Env)
}
