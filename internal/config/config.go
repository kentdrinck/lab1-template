package config

import (
	"fmt"
	"os"

	"github.com/pelletier/go-toml/v2"
)

type DBConfig struct {
	Host     string `toml:"host"`
	Port     string `toml:"port"`
	User     string `toml:"user"`
	Password string `toml:"password"`
	DBName   string `toml:"dbname"`
}

type AppConfig struct {
	DB   DBConfig `toml:"db"`
	Addr string   `toml:"addr"`
}

func LoadConfig(filename string) (*AppConfig, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения файла: %w", err)
	}

	var config AppConfig
	err = toml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("ошибка парсинга TOML: %w", err)
	}

	return &config, nil
}
