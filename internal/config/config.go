package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Filter allows route some messages to another webhook (goons post all shit to directorbot)
type Filter struct {
	Propagate bool     `yaml:"propagate"`
	Patterns  []string `yaml:"patterns"`
	Webhooks  []string `yaml:"webhooks"`
}

// Config comment
type Config struct {
	Admhooks []string          `yaml:"admhooks"`
	Webhooks []string          `yaml:"webhooks"`
	Logs     string            `yaml:"logs"`
	IconUrl  string            `yaml:"icon_url"`
	BotName  string            `yaml:"botname"`
	Filters  map[string]Filter `yaml:"filters"`
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
