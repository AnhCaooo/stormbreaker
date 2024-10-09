package models

type Database struct {
	Username   string `yaml:"user"`
	Password   string `yaml:"pass"`
	Port       string `yaml:"port"`
	Host       string `yaml:"host"`
	Name       string `yaml:"name"`
	Collection string `yaml:"collection"`
}
