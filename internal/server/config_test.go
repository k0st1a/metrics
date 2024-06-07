package server

import (
	"flag"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
				"STORE_INTERVAL":    "100",
				"RESTORE":           "true",
				"PPROF_ADDRESS":     "localhost:9090",
			},
			cfg: Config{
				DatabaseDSN:     "DATABASE_DSN_FROM_ENV",
				ServerAddr:      "localhost:8080",
				FileStoragePath: "FILE_STORAGE_PATH_FROM_ENV",
				HashKey:         "KEY_FROM_ENV",
				StoreInterval:   100,
				Restore:         true,
				PprofServerAddr: "localhost:9090",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			for k, v := range test.env {
				t.Setenv(k, v)
			}

			cfg, err := NewConfig()
			assert.NoError(t, err)
			assert.Equal(t, test.cfg, *cfg)
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
				"-i", "200",
				"-r=false",
				"-p", "localhost:9091",
			},
			cfg: Config{
				DatabaseDSN:     "DATABASE_DSN_FROM_FLAG",
				ServerAddr:      "localhost:8081",
				FileStoragePath: "FILE_STORAGE_PATH_FROM_FLAG",
				HashKey:         "KEY_FROM_FLAG",
				StoreInterval:   200,
				Restore:         false,
				PprofServerAddr: "localhost:9091",
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
				"STORE_INTERVAL":    "300",
				"RESTORE":           "true",
				"PPROF_ADDRESS":     "localhost:9090",
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
			},
			cfg: Config{
				DatabaseDSN:     "DATABASE_DSN_FROM_ENV",
				ServerAddr:      "localhost:8080",
				FileStoragePath: "FILE_STORAGE_PATH_FROM_ENV",
				HashKey:         "KEY_FROM_ENV",
				StoreInterval:   300,
				Restore:         true,
				PprofServerAddr: "localhost:9090",
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
