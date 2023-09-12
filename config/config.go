package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	TestConfig TestConfig `yaml:"testConfig"`
}

type TestConfig struct {
	Components map[string]string `yaml:"components"`
	Aliases    map[string]string `yaml:"aliases"`
}

var conf *Config

func GetComponentPathOrDefault(keyword string) string {
	cfg := GetConfig()

	component, ok := cfg.TestConfig.Components[keyword]
	if ok {
		return component
	}

	// Try the aliases
	component, ok = cfg.TestConfig.Aliases[keyword]
	if ok {
		return component
	}

	root, ok := cfg.TestConfig.Components["root"]
	if !ok {
		log.Fatalln("Could not find root config in lester.yaml, please make sure that's a thing.")
	}
	return root
}

func GetConfig() *Config {
	if conf != nil {
		return conf
	}

	err := viper.Unmarshal(&conf)
	if err != nil {
		log.Fatalln(err)
	}

	return conf
}
