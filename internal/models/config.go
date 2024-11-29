// AnhCao 2024
package models

type Config struct {
	Server   Server   `yaml:"server"`
	Database Database `yaml:"database"`
	Supabase Supabase `yaml:"supabase"`
}

// todo: validate configuration
func (c *Config) Validate() error {
	return nil
}
