// Package config for read config from env and flags.
package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/k0st1a/metrics/internal/pkg/netaddr"
)

const (
	defaultPollInterval   = 2
	defaultReportInterval = 10
	defaultServerAddr     = "localhost:8080"
	defaultHashKey        = ""
	defaultCryptoKey      = ""
	defaultRateLimit      = 1
	defaultConfig         = ""
	defaultAPIType        = "http"
	defaultLogLevel       = "info"
)

// Config - структура с конфигурационными параметрами агента.
type Config struct {
	// ServerAddr - адрес эндпоинта HTTP-сервера (по умолчанию `localhost:8080`).
	// Задается через флаг `-a=<ЗНАЧЕНИЕ>` или переменную окружения `ADDRESS=<ЗНАЧЕНИЕ>`
	ServerAddr string
	// HashKey - ключ для подписи передаваемых данных по алгоритму SHA256 (по умолчанию пустая строка).
	// Задается через флаг `-k=<ЗНАЧЕНИЕ>` или переменную окружения `KEY=<ЗНАЧЕНИЕ>`
	HashKey string
	// CryptoKey - путь до файла с открытым ключом (по умолчанию пустая строка). Если путь задан, то
	// с помощью открытого ключа будут шифровываться сообщения, отправляемые агентом.
	// Задается через флаг `-crypto-key=<ЗНАЧЕНИЕ>` или переменную окружения `CRYPTO_KEY=<ЗНАЧЕНИЕ>`
	CryptoKey string
	// Config - путь до файла конфигурации сервера (по умолчанию пустая строка).
	// Задается через флаг `-c=<ЗНАЧЕНИЕ>` или переменную окружения `CONFIG=<ЗНАЧЕНИЕ>`
	Config string
	// Тип запускаемого сервером API (по умолчанию http). Возможные значения http и grpc.
	// "Задается через флаг `-api-type=<ЗНАЧЕНИЕ>` или переменную окружения `API_TYPE=<ЗНАЧЕНИЕ>")
	APIType string
	// LogLevel - уровень логирования. Возможные значения: debug, info, warn, error (по умолчанию info).
	// Задается через флаг `-log-level=<ЗНАЧЕНИЕ>` или переменную окружения `LOG_LEVEL=<ЗНАЧЕНИЕ>`
	LogLevel string
	// PollInterval - частота опроса метрик из пакета `runtime` (по умолчанию 2 секунды).
	// Задается через флаг `-p=<ЗНАЧЕНИЕ>` или переменную окружения `POLL_INTERVAL=<ЗНАЧЕНИЕ>`
	PollInterval int
	// ReportInterval — частоту отправки метрик на сервер (по умолчанию 10 секунд).
	// Задается через флаг `-r=<ЗНАЧЕНИЕ>` или переменную окружения `REPORT_INTERVAL=<ЗНАЧЕНИЕ>`
	ReportInterval int
	// RateLimit - количество одновременно исходящих запросов на сервер (по умолчанию `1`).
	// Задается через флаг `-l=<ЗНАЧЕНИЕ>` или переменную окружения `RATE_LIMIT=<ЗНАЧЕНИЕ>`
	RateLimit int
}

// New - создать конфигурацию агента из файла конфигурации, аргументов командой строки и переменных окружения.
func New() (*Config, error) {
	var path string

	flag.StringVar(&path, "c", defaultConfig,
		"Путь до файла конфигурации агента (по умолчанию пустая строка).\n"+
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
		ServerAddr:     defaultServerAddr,
		HashKey:        defaultHashKey,
		CryptoKey:      defaultCryptoKey,
		PollInterval:   defaultPollInterval,
		ReportInterval: defaultReportInterval,
		RateLimit:      defaultRateLimit,
		APIType:        defaultAPIType,
		LogLevel:       defaultLogLevel,
	}
}

func (c *Config) applyFromArgsAndEnv() error {
	addr := &netaddr.NetAddress{}
	err := addr.Set(c.ServerAddr)
	if err != nil {
		return fmt.Errorf("addr set error:%w", err)
	}

	flag.Var(addr, "a", "server network address")

	flag.IntVar(&c.PollInterval, "p", c.PollInterval, "metrics polling rate to the server")
	flag.IntVar(&c.ReportInterval, "r", c.ReportInterval, "frequency of sending metrics to the server")
	flag.StringVar(&c.HashKey, "k", c.HashKey,
		"Hash key with which the request body will be encoded "+
			"HTTP Header HashSHA256 will be added to the HTTP request")
	flag.StringVar(&(c.CryptoKey), "crypto-key", c.CryptoKey,
		"Путь до файла с открытым ключом (по умолчанию пустая строка). Если путь задан, то "+
			"с помощью открытого ключа будут шифровываться сообщения, отправляемые агентом.")
	flag.IntVar(&(c.RateLimit), "l", c.RateLimit, "number of simultaneously outgoing requests to the server")
	flag.StringVar(&c.APIType, "api-type", c.APIType,
		"Тип запускаемого сервером API (по умолчанию http). Возможные значения http и grpc.\n"+
			"Задается через флаг `-api-type=<ЗНАЧЕНИЕ>` или переменную окружения `API_TYPE=<ЗНАЧЕНИЕ>")
	flag.StringVar(&c.LogLevel, "log-level", c.LogLevel,
		"Уровень логирования. Задается через флаг `-log-level=<ЗНАЧЕНИЕ>` или переменную окружения "+
			"`LOG_LEVEL=<ЗНАЧЕНИЕ>.\nВозможные значения: debug, info, warn, error.")

	flag.Parse()

	if len(flag.Args()) != 0 {
		return fmt.Errorf("unknown args:%v", flag.Args())
	}

	c.ServerAddr = addr.String()

	sa, ok := os.LookupEnv("ADDRESS")
	if ok {
		c.ServerAddr = sa
	}

	k, ok := os.LookupEnv("KEY")
	if ok {
		c.HashKey = k
	}

	ck, ok := os.LookupEnv("CRYPTO_KEY")
	if ok {
		c.CryptoKey = ck
	}

	pi, ok := os.LookupEnv("POLL_INTERVAL")
	if ok {
		piInt, err := strconv.Atoi(pi)
		if err != nil {
			return fmt.Errorf("POLL_INTERVAL parse error:%w", err)
		}

		c.PollInterval = piInt
	}

	ri, ok := os.LookupEnv("REPORT_INTERVAL")
	if ok {
		riInt, err := strconv.Atoi(ri)
		if err != nil {
			return fmt.Errorf("REPORT_INTERVAL parse error:%w", err)
		}

		c.ReportInterval = riInt
	}

	rl, ok := os.LookupEnv("RATE_LIMIT")
	if ok {
		rlInt, err := strconv.Atoi(rl)
		if err != nil {
			return fmt.Errorf("RATE_LIMIT parse error:%w", err)
		}

		c.RateLimit = rlInt
	}

	at, ok := os.LookupEnv("API_TYPE")
	if ok {
		c.APIType = at
	}

	ll, ok := os.LookupEnv("LOG_LEVEL")
	if ok {
		c.LogLevel = ll
	}

	return nil
}

// JSONConfig - промежуточная структура с конфигурационными параметрами сервера.
// Использользуется для Unmarshal-инга файла в формате JSON в данную структуру.
// Далее данные данной структуры будут использованы для формирования структуры Config.
type JSONConfig struct {
	Address        string `json:"address"`
	ReportInterval string `json:"report_interval"`
	PollInterval   string `json:"poll_interval"`
	CryptoKey      string `json:"crypto_key"`
	APIType        string `json:"api_type"`
	LogLevel       string `json:"log_level"`
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

	if cfg.ReportInterval != "" {
		ri, err := time.ParseDuration(cfg.ReportInterval)
		if err != nil {
			return fmt.Errorf("report interval parse error:%w", err)
		}

		c.ReportInterval = int(ri.Seconds())
	}

	if cfg.PollInterval != "" {
		pi, err := time.ParseDuration(cfg.PollInterval)
		if err != nil {
			return fmt.Errorf("report interval parse error:%w", err)
		}

		c.PollInterval = int(pi.Seconds())
	}

	if cfg.CryptoKey != "" {
		c.CryptoKey = cfg.CryptoKey
	}

	if cfg.APIType != "" {
		c.APIType = cfg.APIType
	}

	if cfg.LogLevel != "" {
		c.LogLevel = cfg.LogLevel
	}

	return nil
}
