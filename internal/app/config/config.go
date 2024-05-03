package config

import (
	"flag"
	"fmt"
	"strconv"
	"strings"
)

// Config структура для хранения конфигурационных параметров/
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
	addr := new(Config)
	_ = flag.Value(addr)
	flag.Var(addr, "a", "Net address host:port")
	flag.StringVar(&BaseURL, "b", "http://localhost/", "Net address host:port")
	flag.Parse()

	if addr.Host == "" {
		addr.Host = "localhost"
	}
	if addr.Port == 0 {
		addr.Port = 8080
	}
	return addr
}

func ParseBaseURL() string {
	flag.StringVar(&BaseURL, "b", "http://localhost/", "Net address host:port")
	return BaseURL
}
