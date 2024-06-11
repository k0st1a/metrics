package server

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/k0st1a/metrics/internal/pkg/netaddr"
	"github.com/rs/zerolog/log"
)

// Config - структура с конфигурационными параметрами агента.
type Config struct {
	// DatabaseDSN - cтрока с адресом подключения к БД.
	// Задается через флаг `-d=<ЗНАЧЕНИЕ>` или переменную окружения `DATABASE_DSN=<ЗНАЧЕНИЕ>`
	DatabaseDSN string
	// ServerAddr - адрес эндпоинта HTTP-сервера (по умолчанию `localhost:8080`).
	// Задается через флаг `-a=<ЗНАЧЕНИЕ>` или переменную окружения `ADDRESS=<ЗНАЧЕНИЕ>`
	ServerAddr string
	// FileStoragePath - полное имя файла, куда сохраняются текущие значения (по умолчанию `/tmp/metrics-db.json`,
	// пустое значение отключает функцию записи на диск).
	// Задается через флаг `-f=<ЗНАЧЕНИЕ>` или переменную окружения `FILE_STORAGE_PATH=<ЗНАЧЕНИЕ>`
	FileStoragePath string
	// HashKey - ключ для подписи передаваемых данных по алгоритму SHA256 (по умолчанию пустая строка).
	// Задается через флаг `-k=<ЗНАЧЕНИЕ>` или переменную окружения `KEY=<ЗНАЧЕНИЕ>`
	HashKey string
	// PprofServerAddr - адрес эндпоинта HTTP-сервера профилировщика pprof (по умолчанию `localhost:8086`).
	// Задается через флаг `-p=<ЗНАЧЕНИЕ>` или переменную окружения `PPROF_ADDRESS=<ЗНАЧЕНИЕ>`
	PprofServerAddr string
	// StoreInterval - интервал времени в секундах, по истечении которого текущие показания сервера сохраняются на
	// диск (по умолчанию 300 секунд, значение `0` делает запись синхронной).
	// Задается через флаг `-i=<ЗНАЧЕНИЕ>` или переменную окружения `STORE_INTERVAL=<ЗНАЧЕНИЕ>`
	StoreInterval int
	// Restore - булево значение (`true/false`), определяющее, загружать или нет ранее сохранённые значения из
	// указанного файла при старте сервера (по умолчанию `true`).
	// Задается через флаг `-r=<ЗНАЧЕНИЕ>` или переменную окружения `RESTORE=<ЗНАЧЕНИЕ>`
	Restore bool
}

const (
	defaultServerAddr      = "localhost:8080"
	defaultStoreInterval   = 300
	defaultFileStoragePath = "/tmp/metrics-db.json"
	defaultRestore         = true
	defaultDatabaseDSN     = ""
	defaultHashKey         = ""
	defaultPprofServerAddr = "localhost:8086"
)

func NewConfig() (*Config, error) {
	cfg := &Config{}

	addr := &netaddr.NetAddress{}
	err := addr.Set(defaultServerAddr)
	if err != nil {
		return nil, fmt.Errorf("addr set error:%w", err)
	}

	flag.Var(addr, "a", "server network address")
	flag.IntVar(&cfg.StoreInterval, "i", defaultStoreInterval,
		"Интервал времени в секундах, по истечении которого текущие показания сервера сохраняются на диск "+
			"(значение 0 делает запись синхронной). Соответствует переменной окружения STORE_INTERVAL")
	flag.StringVar(&cfg.FileStoragePath, "f", defaultFileStoragePath,
		"Полное имя файла, куда сохраняются текущие значения (пустое значение отключает функцию записи на диск). "+
			"Соответствует переменной окружения FILE_STORAGE_PATH")
	flag.BoolVar(&cfg.Restore, "r", defaultRestore,
		"Загружать или нет ранее сохранённые значения из указанного файла при старте сервера."+
			"Соответствует переменной окружения RESTORE")
	flag.StringVar(&cfg.DatabaseDSN, "d", defaultDatabaseDSN,
		"Адрес подключения к БД. Соответствует переменной окружения DATABASE_DSN")
	flag.StringVar(&(cfg.HashKey), "k", defaultHashKey,
		"При наличии ключа во время обработки запроса сервер проверяет соответие полученного и "+
			"вычесленного(от всего тела запроса) хеша.\nПри несовпадении сервер отбрасывает данные и отвечает 400.\n"+
			"При наличии ключа на этапе формирования ответа сервер вычисляет хеш и передает его в HTTP-заголовке"+
			"ответа с именем HashSHA256.")
	flag.StringVar(&cfg.PprofServerAddr, "p", defaultPprofServerAddr, "pprof server address")

	flag.Parse()

	if len(flag.Args()) != 0 {
		return nil, fmt.Errorf("unknown args:%v", flag.Args())
	}

	cfg.ServerAddr = addr.String()

	dbdsn, ok := os.LookupEnv("DATABASE_DSN")
	if ok {
		cfg.DatabaseDSN = dbdsn
	}

	sa, ok := os.LookupEnv("ADDRESS")
	if ok {
		cfg.ServerAddr = sa
	}

	fsp, ok := os.LookupEnv("FILE_STORAGE_PATH")
	if ok {
		cfg.FileStoragePath = fsp
	}

	k, ok := os.LookupEnv("KEY")
	if ok {
		cfg.HashKey = k
	}

	si, ok := os.LookupEnv("STORE_INTERVAL")
	if ok {
		siInt, err := strconv.Atoi(si)
		if err != nil {
			return nil, fmt.Errorf("STORE_INTERVAL parse error:%w", err)
		}

		cfg.StoreInterval = siInt
	}

	rs, ok := os.LookupEnv("RESTORE")
	if ok {
		rsBool, err := strconv.ParseBool(rs)
		if err != nil {
			return nil, fmt.Errorf("RESTORE parse error:%w", err)
		}
		cfg.Restore = rsBool
	}

	ppa, ok := os.LookupEnv("PPROF_ADDRESS")
	if ok {
		cfg.PprofServerAddr = ppa
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
