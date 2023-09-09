package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Config comment
type Config struct {
	Admhooks             []string   `yaml:"admhooks"`
	Webhooks             []string   `yaml:"webhooks"`
	Logs                 string     `yaml:"logs"`
}

// NewConfig Loads the config from the provided path
func NewConfig(path string) (*Config, error) {
	config := &Config{}
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return config, err
	}

	err = yaml.Unmarshal(file, config)
	return config, err
}
