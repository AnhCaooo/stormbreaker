// AnhCao 2024
package models

// Config represents the configuration structure for the application.
// It includes settings for the server, database, Supabase, and message broker.
type Config struct {
	Server        Server   `yaml:"server"`
	Database      Database `yaml:"database"`
	Supabase      Supabase `yaml:"supabase"`
	MessageBroker Broker   `yaml:"message_broker"`
}

// Server represents the configuration settings for the server.
// It includes the port and host information required to run the server.
// The fields are annotated for YAML parsing.
type Server struct {
	Port string `yaml:"port"`
	Host string `yaml:"host"`
}

// Broker represents the configuration settings for connecting to a broker.
// It includes the username, password, port, and host information required
// for establishing a connection.
type Broker struct {
	// The username for authentication with the broker.
	Username string `yaml:"user"`
	// The password for authentication with the broker.
	Password string `yaml:"pass"`
	// The port number on which the broker is running.
	Port string `yaml:"port"`
	// The hostname or IP address of the broker.
	Host string `yaml:"host"`
}

// Database represents the configuration settings for connecting to a database.
// It includes fields for the username, password, port, host, database name, and collection.
type Database struct {
	// The username for database authentication.
	Username string `yaml:"user"`
	// The password for database authentication.
	Password string `yaml:"pass"`
	// The port number on which the database server is listening.
	Port string `yaml:"port"`
	// The hostname or IP address of the database server.
	Host string `yaml:"host"`
	// The name of the database (e.g., MongoDB, PostgreSQL).
	Name string `yaml:"name"`
	// The name of the collection within the database.
	Collection string `yaml:"collection"`
}

// Supabase represents the configuration settings for connecting to Supabase.
type Supabase struct {
	Auth auth `yaml:"auth"`
}

type auth struct {
	JwtSecret string `yaml:"jwt_secret"`
}

// todo: validate configuration
func (c *Config) Validate() error {
	return nil
}
