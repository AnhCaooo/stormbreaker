package models

type Supabase struct {
	Auth auth `yaml:"auth"`
}

type auth struct {
	JwtSecret string `yaml:"jwt_secret"`
}
