package models

type Database struct {
	Username   string `yaml:"user"`
	Password   string `yaml:"pass"`
	Port       string `yaml:"port"`
	Host       string `yaml:"host"`
	Name       string `yaml:"name"` // name of the database. Example: MongoDB, PostgreSQL, etc.
	Collection string `yaml:"collection"`
}
