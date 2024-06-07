package agent

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
				"ADDRESS":         "ADDRESS_FROM_ENV",
				"KEY":             "KEY_FROM_ENV",
				"POLL_INTERVAL":   "100",
				"REPORT_INTERVAL": "200",
				"RATE_LIMIT":      "300",
			},
			cfg: Config{
				ServerAddr:     "ADDRESS_FROM_ENV",
				HashKey:        "KEY_FROM_ENV",
				PollInterval:   100,
				ReportInterval: 200,
				RateLimit:      300,
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
				"-a", "localhost:8081",
				"-p", "100",
				"-r", "200",
				"-k", "KEY_FROM_FLAG",
				"-l", "300",
			},
			cfg: Config{
				ServerAddr:     "localhost:8081",
				PollInterval:   100,
				ReportInterval: 200,
				HashKey:        "KEY_FROM_FLAG",
				RateLimit:      300,
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
				"ADDRESS":         "ADDRESS_FROM_ENV",
				"KEY":             "KEY_FROM_ENV",
				"POLL_INTERVAL":   "100",
				"REPORT_INTERVAL": "200",
				"RATE_LIMIT":      "300",
			},
			args: []string{
				"cmd",
				"-a", "localhost:8081",
				"-p", "100",
				"-r", "200",
				"-k", "KEY_FROM_FLAG",
				"-l", "300",
			},
			cfg: Config{
				ServerAddr:     "ADDRESS_FROM_ENV",
				HashKey:        "KEY_FROM_ENV",
				PollInterval:   100,
				ReportInterval: 200,
				RateLimit:      300,
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
