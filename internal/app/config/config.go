package config

import (
	"flag"
	"fmt"
	"github.com/caarlos0/env/v6"
	"log"
	"strconv"
	"strings"
)

type HostPort struct {
	Hp              []string `env:"SERVER_ADDRESS" envSeparator:":"`
	BaseURL         string   `env:"BASE_URL"`
	FileStoragePath string   `env:"FILE_STORAGE"`
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
var FileStoragePath string

func ParseConfig() *Config {
	var cfg HostPort
	if err := env.Parse(&cfg); err != nil {
		log.Println(err)
		return nil
	}

	//парсим флаги командной строки
	addr := &Config{}
	flag.Var(addr, "a", "Net address host:port")
	baseURL := flag.String("b", "", "Net address base url")
	storagePath := flag.String("f", "", "File storage path")
	flag.Parse()

	if len(cfg.Hp) == 0 {
		if addr.Host == "" {
			addr.Host = "localhost"
		}
		if addr.Port == 0 {
			addr.Port = 8080
		}
	} else {
		addr.Host = cfg.Hp[0]
		port, err := strconv.Atoi(cfg.Hp[1])
		if err != nil {
			log.Println(err)
			return nil
		}
		addr.Port = port
	}

	BaseURL = ""
	if cfg.BaseURL != "" {
		BaseURL = strings.TrimSuffix(cfg.BaseURL, "/")
	} else if *baseURL != "" {
		BaseURL = *baseURL
	}

	if cfg.FileStoragePath != "" {
		FileStoragePath = cfg.FileStoragePath
	} else if *storagePath != "" {
		FileStoragePath = *storagePath
	} else {
		FileStoragePath = "/tmp/short-url-db.json"
		//FileStoragePath = "C:\\Users\\Олег\\go\\src\\yandex-praktikum\\project3\\go-musthave-shortener-tpl\\tmp\\short-url-db.json"
	}

	return addr
}
