package agent

type Config struct {
	ServerAddr     string `env:"ADDRESS"`
	PollInterval   int    `env:"POLL_INTERVAL"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
}

func NewConfig() Config {
	return Config{
		ServerAddr:     "localhost:8080",
		PollInterval:   2,
		ReportInterval: 10,
	}
}
