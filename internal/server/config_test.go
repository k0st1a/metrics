package server

import (
	"flag"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigFromFile(t *testing.T) {
	origStateFun := func() {
		func(args []string) {
			//nolint:reassign //for tests only
			os.Args = args
		}(os.Args)

		func(cl *flag.FlagSet) {
			//nolint:reassign //for tests only
			flag.CommandLine = cl
		}(flag.CommandLine)
	}
	defer origStateFun()

	tests := []struct {
		name string
		env  map[string]string
		args []string
		cfg  Config
	}{
		{
			name: "Check config from file",
			env: map[string]string{
				"CONFIG": "./config_test.json",
			},
			args: []string{
				"cmd",
				"-c", "./config_test.json",
			},
			cfg: Config{
				DatabaseDSN:     "DATABASE_DSN_FROM_FILE",
				ServerAddr:      "localhost:8090",
				FileStoragePath: "FILE_STORAGE_PATH_FROM_FILE",
				CryptoKey:       "CRYPTO_KEY_FROM_FILE",
				TrustedSubnet:   "TRUSTED_SUBNET_FROM_FILE",
				StoreInterval:   500,
				Restore:         false,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			for k, v := range test.env {
				t.Setenv(k, v)
			}

			//nolint:reassign //for tests only
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

			//nolint:reassign //for tests only
			os.Args = test.args

			cfg, err := NewConfig()
			assert.NoError(t, err)

			assert.Equal(t, test.cfg.DatabaseDSN, cfg.DatabaseDSN)
			assert.Equal(t, test.cfg.ServerAddr, cfg.ServerAddr)
			assert.Equal(t, test.cfg.FileStoragePath, cfg.FileStoragePath)
			assert.Equal(t, test.cfg.CryptoKey, cfg.CryptoKey)
			assert.Equal(t, test.cfg.StoreInterval, cfg.StoreInterval)
			assert.Equal(t, test.cfg.Restore, cfg.Restore)
			origStateFun()
		})
	}
}

func TestConfigFromEnv(t *testing.T) {
	origStateFun := func() {
		func(args []string) {
			//nolint:reassign //for tests only
			os.Args = args
		}(os.Args)

		func(cl *flag.FlagSet) {
			//nolint:reassign //for tests only
			flag.CommandLine = cl
		}(flag.CommandLine)
	}
	defer origStateFun()

	tests := []struct {
		name string
		env  map[string]string
		cfg  Config
	}{
		{
			name: "Check config from env",
			env: map[string]string{
				"DATABASE_DSN":      "DATABASE_DSN_FROM_ENV",
				"ADDRESS":           "localhost:8080",
				"FILE_STORAGE_PATH": "FILE_STORAGE_PATH_FROM_ENV",
				"KEY":               "KEY_FROM_ENV",
				"CRYPTO_KEY":        "CRYPTO_KEY_FROM_ENV",
				"STORE_INTERVAL":    "100",
				"RESTORE":           "true",
				"PPROF_ADDRESS":     "localhost:9090",
				"TRUSTED_SUBNET":    "TRUSTED_SUBNET_FROM_ENV",
			},
			cfg: Config{
				DatabaseDSN:     "DATABASE_DSN_FROM_ENV",
				ServerAddr:      "localhost:8080",
				FileStoragePath: "FILE_STORAGE_PATH_FROM_ENV",
				HashKey:         "KEY_FROM_ENV",
				CryptoKey:       "CRYPTO_KEY_FROM_ENV",
				StoreInterval:   100,
				Restore:         true,
				PprofServerAddr: "localhost:9090",
				TrustedSubnet:   "TRUSTED_SUBNET_FROM_ENV",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			for k, v := range test.env {
				t.Setenv(k, v)
			}

			//nolint:reassign //for tests only
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

			cfg, err := NewConfig()
			assert.NoError(t, err)
			assert.Equal(t, test.cfg, *cfg)
			origStateFun()
		})
	}
}

func TestConfigFromFlags(t *testing.T) {
	origStateFun := func() {
		func(args []string) {
			//nolint:reassign //for tests only
			os.Args = args
		}(os.Args)

		func(cl *flag.FlagSet) {
			//nolint:reassign //for tests only
			flag.CommandLine = cl
		}(flag.CommandLine)
	}
	defer origStateFun()

	tests := []struct {
		name string
		args []string
		cfg  Config
	}{
		{
			name: "Check config from flags",
			args: []string{
				"cmd",
				"-d", "DATABASE_DSN_FROM_FLAG",
				"-a", "localhost:8081",
				"-f", "FILE_STORAGE_PATH_FROM_FLAG",
				"-k", "KEY_FROM_FLAG",
				"-crypto-key", "CRYPTO_KEY_FROM_FLAG",
				"-i", "200",
				"-r=false",
				"-p", "localhost:9091",
				"-t", "TRUSTED_SUBNET_FROM_FLAG",
			},
			cfg: Config{
				DatabaseDSN:     "DATABASE_DSN_FROM_FLAG",
				ServerAddr:      "localhost:8081",
				FileStoragePath: "FILE_STORAGE_PATH_FROM_FLAG",
				HashKey:         "KEY_FROM_FLAG",
				CryptoKey:       "CRYPTO_KEY_FROM_FLAG",
				StoreInterval:   200,
				Restore:         false,
				PprofServerAddr: "localhost:9091",
				TrustedSubnet:   "TRUSTED_SUBNET_FROM_FLAG",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			//nolint:reassign //for tests only
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

			//nolint:reassign //for tests only
			os.Args = test.args

			cfg, err := NewConfig()
			assert.NoError(t, err)
			assert.Equal(t, test.cfg, *cfg)
			origStateFun()
		})
	}
}

func TestConfig(t *testing.T) {
	origStateFun := func() {
		func(args []string) {
			//nolint:reassign //for tests only
			os.Args = args
		}(os.Args)

		func(cl *flag.FlagSet) {
			//nolint:reassign //for tests only
			flag.CommandLine = cl
		}(flag.CommandLine)
	}
	defer origStateFun()

	tests := []struct {
		name string
		env  map[string]string
		args []string
		cfg  Config
	}{
		{
			name: "Check config from args and env",
			env: map[string]string{
				"DATABASE_DSN":      "DATABASE_DSN_FROM_ENV",
				"ADDRESS":           "localhost:8080",
				"FILE_STORAGE_PATH": "FILE_STORAGE_PATH_FROM_ENV",
				"KEY":               "KEY_FROM_ENV",
				"CRYPTO_KEY":        "CRYPTO_KEY_FROM_ENV",
				"STORE_INTERVAL":    "300",
				"RESTORE":           "true",
				"PPROF_ADDRESS":     "localhost:9090",
				"TRUSTED_SUBNET":    "TRUSTED_SUBNET_FROM_ENV",
			},
			args: []string{
				"cmd",
				"-d", "DATABASE_DSN_FROM_FLAG",
				"-a", "localhost:8081",
				"-f", "FILE_STORAGE_PATH_FROM_FLAG",
				"-k", "KEY_FROM_FLAG",
				"-i", "400",
				"-r=false",
				"-p", "localhost:9091",
				"-t", "TRUSTED_SUBNET_FROM_FLAG",
			},
			cfg: Config{
				DatabaseDSN:     "DATABASE_DSN_FROM_ENV",
				ServerAddr:      "localhost:8080",
				FileStoragePath: "FILE_STORAGE_PATH_FROM_ENV",
				HashKey:         "KEY_FROM_ENV",
				CryptoKey:       "CRYPTO_KEY_FROM_ENV",
				StoreInterval:   300,
				Restore:         true,
				PprofServerAddr: "localhost:9090",
				TrustedSubnet:   "TRUSTED_SUBNET_FROM_ENV",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			//nolint:reassign //for tests only
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

			for k, v := range test.env {
				t.Setenv(k, v)
			}
			//nolint:reassign //for tests only
			os.Args = test.args

			cfg, err := NewConfig()
			assert.NoError(t, err)
			assert.Equal(t, test.cfg, *cfg)
			origStateFun()
		})
	}
}
