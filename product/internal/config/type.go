package config

type (
	Config struct {
		GRPC   GRPC   `yaml:"grpc"`
		Server Server `yaml:"server"`
	}

	GRPC struct {
		Address string `yaml:"address"`
	}

	Server struct {
		Port int `yaml:"port"`
	}
)
