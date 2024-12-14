package models

type Broker struct {
	Username string `yaml:"user"`
	Password string `yaml:"pass"`
	Port     string `yaml:"port"`
	Host     string `yaml:"host"`
}
