package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Config comment
type Config struct {
	IsWebhook            bool       `yaml:"is_webhook"`
	ForwarderToken       string     `yaml:"forwarder_token"`
	ErrorLogChannelID    string     `yaml:"error_log_channel_id"`
	DiscordRusWebhookURL string     `yaml:"discord_rus_webhook_url"`
	Listeners            []Listener `yaml:"discord_listeners"`
	ConfigPath           string     `yaml:"-"`
	MetricsEnabled       bool       `yaml:"mectrics_enabled"`
	MetricsEndpoint      string     `yaml:"metrics_endpoint"`
}

// NewConfig Loads the config from the provided path
func NewConfig(path string) (*Config, error) {
	config := &Config{
		ConfigPath: path,
		Listeners:  make([]Listener, 0),
	}
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return config, err
	}

	err = yaml.Unmarshal(file, config)
	return config, err
}

// SaveConfig Saves the currently loaded config to file
func (c *Config) SaveConfig() error {
	file, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(c.ConfigPath, file, 0644)

}

// TODO Why does this exist?
// ListenerNames Returns a list of listernames
func (c *Config) ListenerNames() []string {
	var names = []string{}

	for _, l := range c.Listeners {
		names = append(names, l.Name)
	}

	names = append(names, "Back")
	return names
}

//ListenerTokens returns all listener tokens
func (c *Config) ListenerTokens() []string {
	var toks = []string{}
	for _, l := range c.Listeners {
		toks = append(toks, l.UserToken)
	}
	return toks
}

func (c *Config) ChannelMap() map[string]string {
	cm := make(map[string]string)
	for _, l := range c.Listeners {
		for _, c := range l.Channels {
			cm[c.ListenChannelID] = c.ForwardChannelID
		}
	}
	return cm
}
