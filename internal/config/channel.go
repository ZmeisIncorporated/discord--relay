package config

// Channel the discord channels we are listening to
type Channel struct {
	Name             string `yaml:"name"`
	ListenChannelID  string `yaml:"listen_channel_id"`
	ForwardChannelID string `yaml:"forward_channel_id"`
	WebhookID        string `yaml:"webhook_id"`
	IsMuted          bool   `yaml:"is_muted"`
}
