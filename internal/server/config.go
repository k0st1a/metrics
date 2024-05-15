package server

import (
	"flag"
	"fmt"

	"github.com/caarlos0/env/v10"
	"github.com/k0st1a/metrics/internal/pkg/netaddr"
	"github.com/rs/zerolog/log"
)

type Config struct {
	DatabaseDSN     string `env:"DATABASE_DSN"` // Database Data Source Name
	ServerAddr      string `env:"ADDRESS"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	HashKey         string `env:"KEY"`
	StoreInterval   int    `env:"STORE_INTERVAL"`
	Restore         bool   `env:"RESTORE"`
}

const (
	defaultServerAddr      = "localhost:8080"
	defaultStoreInterval   = 300
	defaultFileStoragePath = "/tmp/metrics-db.json"
	defaultRestore         = true
	defaultDatabaseDSN     = ""
	defaultHashKey         = ""
)

func newConfig() *Config {
	return &Config{
		ServerAddr:      defaultServerAddr,
		StoreInterval:   defaultStoreInterval,
		FileStoragePath: defaultFileStoragePath,
		Restore:         defaultRestore,
		DatabaseDSN:     defaultDatabaseDSN,
		HashKey:         defaultHashKey,
	}
}

func parseEnv(cfg *Config) error {
	err := env.Parse(cfg)
	if err != nil {
		return fmt.Errorf("env parse error:%w", err)
	}
	return nil
}

func parseFlags(cfg *Config) error {
	addr := &netaddr.NetAddress{}
	err := addr.Set(cfg.ServerAddr)
	if err != nil {
		return fmt.Errorf("addr set error:%w", err)
	}

	// если интерфейс не реализован,
	// здесь будет ошибка компиляции
	_ = flag.Value(addr)
	flag.Var(addr, "a", "server network address")
	flag.IntVar(&cfg.StoreInterval, "i", cfg.StoreInterval,
		"Интервал времени в секундах, по истечении которого текущие показания сервера сохраняются на диск "+
			"(значение 0 делает запись синхронной). Соответствует переменной окружения STORE_INTERVAL")
	flag.StringVar(&cfg.FileStoragePath, "f", cfg.FileStoragePath,
		"Полное имя файла, куда сохраняются текущие значения (пустое значение отключает функцию записи на диск). "+
			"Соответствует переменной окружения FILE_STORAGE_PATH")
	flag.BoolVar(&cfg.Restore, "r", cfg.Restore,
		"Загружать или нет ранее сохранённые значения из указанного файла при старте сервера."+
			"Соответствует переменной окружения RESTORE")
	flag.StringVar(&cfg.DatabaseDSN, "d", cfg.DatabaseDSN,
		"Адрес подключения к БД. Соответствует переменной окружения DATABASE_DSN")
	flag.StringVar(&(cfg.HashKey), "k", cfg.HashKey,
		"При наличии ключа во время обработки запроса сервер проверяет соответие полученного и "+
			"вычесленного(от всего тела запроса) хеша.\nПри несовпадении сервер отбрасывает данные и отвечает 400.\n"+
			"При наличии ключа на этапе формирования ответа сервер вычисляет хеш и передает его в HTTP-заголовке"+
			"ответа с именем HashSHA256.")

	flag.Parse()

	if len(flag.Args()) != 0 {
		return fmt.Errorf("unknown args")
	}

	cfg.ServerAddr = addr.String()

	return nil
}

func collectConfig() (cfg *Config, err error) {
	cfg = newConfig()

	err = parseFlags(cfg)
	if err != nil {
		return nil, err
	}

	err = parseEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func printConfig(cfg *Config) {
	log.Debug().
		Str("cfg.DatabaseDSN", cfg.DatabaseDSN).
		Str("cfg.ServerAddr", cfg.ServerAddr).
		Int("cfg.StoreInterval", cfg.StoreInterval).
		Str("cfg.FileStoragePath", cfg.FileStoragePath).
		Bool("cfg.Restore", cfg.Restore).
		Msg("printConfig")
}
