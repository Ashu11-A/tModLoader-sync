package configs

import (
	"fmt"
	"os"
	"github.com/spf13/pflag"
)

type Config struct {
	Port int
}

func Load() *Config {
	cfg := &Config{}
	pflag.IntVar(&cfg.Port, "port", 0, "Port to run the server on (required)")
	pflag.Parse()

	if cfg.Port == 0 {
		fmt.Println("Error: --port is required.")
		pflag.Usage()
		os.Exit(1)
	}
	return cfg
}

func (c *Config) Address() string {
	return fmt.Sprintf(":%d", c.Port)
}
