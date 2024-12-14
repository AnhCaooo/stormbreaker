// AnhCao 2024
package models

type Config struct {
	Server        Server   `yaml:"server"`
	Database      Database `yaml:"database"`
	Supabase      Supabase `yaml:"supabase"`
	MessageBroker Broker   `yaml:"message_broker"`
}

// todo: validate configuration
func (c *Config) Validate() error {
	return nil
}
