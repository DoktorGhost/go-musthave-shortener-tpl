package config

import (
	"flag"
	"fmt"
	"github.com/caarlos0/env/v6"
	"log"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Host            string
	Port            string
	BaseURL         string
	FileStoragePath string
	DatabaseDSN     string
}

type EnvStruct struct {
	Hp              []string `env:"SERVER_ADDRESS" envSeparator:":"`
	BaseURL         string   `env:"BASE_URL"`
	FileStoragePath string   `env:"FILE_STORAGE_PATH"`
	DatabaseDSN     string   `env:"DATABASE_DSN"`
}

//$env:DATABASE_DSN = "postgres://admin:admin@localhost:5433/postgres?sslmode=disable"
//echo $env:DATABASE_DSN

type HostPort struct {
	Host string
	Port int
}

func (c *HostPort) String() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

func (c *HostPort) Set(value string) error {
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

func ParseConfig() *Config {
	var envStruct EnvStruct
	//считываем все переменны окружения в cfg
	if err := env.Parse(&envStruct); err != nil {
		log.Println(err)
		return nil
	}

	config := new(Config)
	hostPort := new(HostPort)

	//парсим флаги командной строки
	flag.Var(hostPort, "a", "Net address host:port")
	flag.StringVar(&config.BaseURL, "b", "", "Net address base url")
	flag.StringVar(&config.DatabaseDSN, "d", "", "Connect DB string")
	flag.StringVar(&config.FileStoragePath, "f", "/tmp/short-url-db.json", "File storage path")
	flag.Parse()

	_, exists := os.LookupEnv("SERVER_ADDRESS")
	if exists {
		config.Host = envStruct.Hp[0]
		config.Port = envStruct.Hp[1]
	} else {
		if hostPort.Host == "" {
			config.Host = "localhost"
		} else {
			config.Host = hostPort.Host
		}
		if hostPort.Port == 0 {
			config.Port = "8080"
		} else {
			config.Port = strconv.Itoa(hostPort.Port)
		}
	}

	value, exists := os.LookupEnv("BASE_URL")
	if exists {
		config.BaseURL = value
	}

	value, ok := os.LookupEnv("FILE_STORAGE_PATH")
	if ok {
		config.FileStoragePath = value
	}

	value, ok = os.LookupEnv("DATABASE_DSN")
	if ok {
		config.DatabaseDSN = value
	}

	return config
}
