package main

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type ServerConfig struct {
	Port uint `yaml:"port"`
}

type Config struct {
	Server ServerConfig `yaml:"server"`
}

func NewConfig(filename string) (cfg *Config, err error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}
	cfg = &Config{}
	if err = yaml.Unmarshal(data, cfg); err != nil {
		return
	}

	// validators --
	if cfg.Server.Port <= 0 {
		err = fmt.Errorf("Invalid server port %q", cfg.Server.Port)
		return
	}
	return
}
