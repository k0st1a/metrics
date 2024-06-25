// Package server for read config from env and flags.
package server

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/k0st1a/metrics/internal/pkg/netaddr"
)

// Config - структура с конфигурационными параметрами сервера.
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
	// CryptoKey - путь до файла с приватным ключом (по умолчанию пустая строка). Если путь задан, то
	// с помощью приватного ключа будут дешифровываться сообщения, получаемые сервером.
	// Задается через флаг `-crypto-key=<ЗНАЧЕНИЕ>` или переменную окружения `CRYPTO_KEY=<ЗНАЧЕНИЕ>`
	CryptoKey string
	// PprofServerAddr - адрес эндпоинта HTTP-сервера профилировщика pprof (по умолчанию `localhost:8086`).
	// Задается через флаг `-p=<ЗНАЧЕНИЕ>` или переменную окружения `PPROF_ADDRESS=<ЗНАЧЕНИЕ>`
	PprofServerAddr string
	// Config - путь до файла конфигурации сервера (по умолчанию пустая строка).
	// Задается через флаг `-c=<ЗНАЧЕНИЕ>` или переменную окружения `CONFIG=<ЗНАЧЕНИЕ>`
	Config string
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
	defaultCryptoKey       = ""
	defaultPprofServerAddr = "localhost:8086"
	defaultConfig          = ""
)

// NewConfig - создать конфигурацию сервера из файла конфигурации, аргументов командой строки и переменных окружения.
func NewConfig() (*Config, error) {
	var path string

	flag.StringVar(&path, "c", defaultConfig,
		"Путь до файла конфигурации сервера (по умолчанию пустая строка).\n"+
			"Задается через флаг `-c=<ЗНАЧЕНИЕ>` или переменную окружения `CONFIG=<ЗНАЧЕНИЕ>`")

	c, ok := os.LookupEnv("CONFIG")
	if ok {
		path = c
	}

	cfg := newDefaultConfig()
	cfg.Config = path

	if path != "" {
		err := cfg.applyFromFile(path)
		if err != nil {
			return nil, fmt.Errorf("apply config from file(%v) error:%w", path, err)
		}
	}

	err := cfg.applyFromArgsAndEnv()
	if err != nil {
		return nil, fmt.Errorf("apply config from args and env:%w", err)
	}

	return cfg, nil
}

func newDefaultConfig() *Config {
	return &Config{
		DatabaseDSN:     defaultDatabaseDSN,
		ServerAddr:      defaultServerAddr,
		FileStoragePath: defaultFileStoragePath,
		HashKey:         defaultHashKey,
		CryptoKey:       defaultCryptoKey,
		PprofServerAddr: defaultPprofServerAddr,
		Config:          defaultConfig,
		StoreInterval:   defaultStoreInterval,
		Restore:         defaultRestore,
	}
}

func (c *Config) applyFromArgsAndEnv() error {
	addr := &netaddr.NetAddress{}
	err := addr.Set(c.ServerAddr)
	if err != nil {
		return fmt.Errorf("addr set error:%w", err)
	}

	flag.Var(addr, "a", "server network address")
	flag.IntVar(&c.StoreInterval, "i", c.StoreInterval,
		"Интервал времени в секундах, по истечении которого текущие показания сервера сохраняются на диск "+
			"(значение 0 делает запись синхронной).\nСоответствует переменной окружения STORE_INTERVAL")
	flag.StringVar(&c.FileStoragePath, "f", c.FileStoragePath,
		"Полное имя файла, куда сохраняются текущие значения (пустое значение отключает функцию записи на диск).\n"+
			"Соответствует переменной окружения FILE_STORAGE_PATH")
	flag.BoolVar(&c.Restore, "r", c.Restore,
		"Загружать или нет ранее сохранённые значения из указанного файла при старте сервера."+
			"Соответствует переменной окружения RESTORE")
	flag.StringVar(&c.DatabaseDSN, "d", c.DatabaseDSN,
		"Адрес подключения к БД. Соответствует переменной окружения DATABASE_DSN")
	flag.StringVar(&c.HashKey, "k", c.HashKey,
		"При наличии ключа во время обработки запроса сервер проверяет соответие полученного и "+
			"вычесленного(от всего тела запроса) хеша.\nПри несовпадении сервер отбрасывает данные и отвечает 400.\n"+
			"При наличии ключа на этапе формирования ответа сервер вычисляет хеш и передает его в HTTP-заголовке"+
			"ответа с именем HashSHA256.")
	flag.StringVar(&c.CryptoKey, "crypto-key", c.CryptoKey,
		"Путь до файла с приватным ключом (по умолчанию пустая строка).\nЕсли путь задан, то "+
			"с помощью приватного ключа будут дешифровываться сообщения, получаемые сервером.")
	flag.StringVar(&c.PprofServerAddr, "p", c.PprofServerAddr, "pprof server address")

	flag.Parse()

	if len(flag.Args()) != 0 {
		return fmt.Errorf("unknown args:%v", flag.Args())
	}

	c.ServerAddr = addr.String()

	dbdsn, ok := os.LookupEnv("DATABASE_DSN")
	if ok {
		c.DatabaseDSN = dbdsn
	}

	sa, ok := os.LookupEnv("ADDRESS")
	if ok {
		c.ServerAddr = sa
	}

	fsp, ok := os.LookupEnv("FILE_STORAGE_PATH")
	if ok {
		c.FileStoragePath = fsp
	}

	k, ok := os.LookupEnv("KEY")
	if ok {
		c.HashKey = k
	}

	ck, ok := os.LookupEnv("CRYPTO_KEY")
	if ok {
		c.CryptoKey = ck
	}

	si, ok := os.LookupEnv("STORE_INTERVAL")
	if ok {
		siInt, err := strconv.Atoi(si)
		if err != nil {
			return fmt.Errorf("STORE_INTERVAL parse error:%w", err)
		}

		c.StoreInterval = siInt
	}

	rs, ok := os.LookupEnv("RESTORE")
	if ok {
		rsBool, err := strconv.ParseBool(rs)
		if err != nil {
			return fmt.Errorf("RESTORE parse error:%w", err)
		}
		c.Restore = rsBool
	}

	ppa, ok := os.LookupEnv("PPROF_ADDRESS")
	if ok {
		c.PprofServerAddr = ppa
	}

	return nil
}

// JSONConfig - промежуточная структура с конфигурационными параметрами сервера.
// Использользуется для Unmarshal-инга файла в формате JSON в данную структуру.
// Далее данные данной структуры будут использованы для формирования структуры Config.
type JSONConfig struct {
	Address         string `json:"address"`
	DatabaseDSN     string `json:"database_dsn"`
	FileStoragePath string `json:"file_storage_path"`
	CryptoKey       string `json:"crypto_key"`
	StoreInterval   string `json:"store_interval"`
	Restore         bool   `json:"restore"`
}

func (c *Config) applyFromFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("os file(%s) read error:%w", path, err)
	}

	var cfg JSONConfig
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return fmt.Errorf("json unmarshal error:%w", err)
	}

	if cfg.Address != "" {
		c.ServerAddr = cfg.Address
	}

	c.Restore = cfg.Restore

	if cfg.StoreInterval != "" {
		i, err := time.ParseDuration(cfg.StoreInterval)
		if err != nil {
			return fmt.Errorf("store interval parse error:%w", err)
		}

		c.StoreInterval = int(i.Seconds())
	}

	if cfg.FileStoragePath != "" {
		c.FileStoragePath = cfg.FileStoragePath
	}

	if cfg.DatabaseDSN != "" {
		c.DatabaseDSN = cfg.DatabaseDSN
	}

	if cfg.CryptoKey != "" {
		c.CryptoKey = cfg.CryptoKey
	}

	return nil
}
