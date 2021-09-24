package config

type Config struct {
	Server   Server   `yaml:"server"`
	Database Database `yaml:"database"`
	Redis    Redis    `yaml:"redis"`
}

type Server struct {
	GPort int `yaml:"g_port"`
}

type Database struct {
	Driver   string `yaml:"driver"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	DBName   string `yaml:"db_name"`
	Password string `yaml:"password"`
}

type Redis struct {
	Addr       string `yaml:"addr"`
	Password   string `yaml:"password"`
	DB         int    `yaml:"db"`
	MaxRetries int    `yaml:"max_retries"`
}

