package config

import (
	"flag"
	"fmt"
	"github.com/caarlos0/env/v6"
	"strconv"
	"strings"
)

// Config структура для хранения конфигурационных параметров/
type HostPort struct {
	Hp      []string `env:"SERVER_ADDRESS" envSeparator:":"`
	BaseURL string   `env:"BASE_URL"`
}

type Config struct {
	Host string
	Port int
}

func (c *Config) String() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

func (c *Config) Set(value string) error {
	hp := strings.Split(value, ":")
	if len(hp) != 2 {
		return fmt.Errorf("invalid host:port format: %s", value)
	}

	port, err := strconv.Atoi(hp[1])
	if err != nil {
		return err
	}
	c.Host = hp[0]
	c.Port = port
	return nil
}

var BaseURL string

func ParseConfig() *Config {
	var cfg HostPort
	env.Parse(&cfg)
	fmt.Println(cfg)

	addr := new(Config)
	_ = flag.Value(addr)
	flag.Var(addr, "a", "Net address host:port")
	flag.StringVar(&BaseURL, "b", "", "Net address base url")
	flag.Parse()

	if len(cfg.Hp) == 0 {
		if addr.Host == "" {
			addr.Host = "localhost"
			fmt.Println("defaulting to localhost:", addr.Host)
		}
		if addr.Port == 0 {
			addr.Port = 8080
			fmt.Println("defaulting to 8080:", addr.Port)
		}
	} else {
		addr.Host = cfg.Hp[0]
		port, err := strconv.Atoi(cfg.Hp[1])
		if err != nil {
			return nil
		}
		addr.Port = port
	}
	if cfg.BaseURL != "" {
		BaseURL = strings.TrimSuffix(cfg.BaseURL, "/")
	}

	return addr
}
