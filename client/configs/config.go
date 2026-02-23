package configs

import (
	"fmt"
	"os"
	"github.com/spf13/pflag"
)

type Config struct {
	Host string
	Port int
}

func Load() *Config {
	cfg := &Config{}
	pflag.StringVar(&cfg.Host, "host", "", "Server host address (required)")
	pflag.IntVar(&cfg.Port, "port", 0, "Server port (required)")
	pflag.Parse()

	if cfg.Host == "" || cfg.Port == 0 {
		fmt.Println("Error: --host and --port are required.")
		pflag.Usage()
		os.Exit(1)
	}
	return cfg
}
