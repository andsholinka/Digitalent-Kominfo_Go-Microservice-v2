package config

type Config struct {
	Port        string
	AuthService AuthService `mapstructure:"auth_service"`
	Database    Database
}
