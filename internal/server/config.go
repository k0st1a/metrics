package server

type Config struct {
	ServerAddr string `env:"ADDRESS"`
	User       string `env:"USER"`
}

func NewConfig() Config {
	return Config{
		ServerAddr: "localhost:8080",
	}
}
