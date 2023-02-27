package config

import "fmt"

// Listener is the user we wish to listen from
type Listener struct {
	UserToken string    `yaml:"user_token"`
	Name      string    `yaml:"name"`
	Channels  []Channel `yaml:"channels"`
}

//ListenerFromToken returns a listener based on the provided token
func (c *Config) ListenerFromToken(token string) (*Listener, error) {
	for _, listener := range c.Listeners {
		if listener.UserToken == token {
			return &listener, nil
		}
	}
	return nil, fmt.Errorf("Token Not Found")
}
