package models

type Config struct {
	Server   Server   `yaml:"server"`
	Database Database `yaml:"database"`
}
