package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type App struct {
	ServerAddress string `json:"server_address"`
}

type Http struct {
	Port         string `json:"port"`
	IdleTimeout  int    `json:"idle_timeout"`
	WriteTimeout int    `json:"write_timeout"`
	ReadTimeout  int    `json:"read_timeout"`
}

type Database struct {
	DBname string `json:"db_name"`
	DBhost string `json:"db_host"`
	DBuser string `json:"db_user"`
	DBpass string `json:"db_pass"`
	DBport int    `json:"db_port"`
}

type Config struct {
	App      `json:"app"`
	Http     `json:"http"`
	Database `json:"database"`
}

func InitConfig(path string) (Config, error) {
	var cfg Config
	file, err := os.Open(path)
	if err != nil {
		return Config{}, fmt.Errorf("open %s error: %w", path, err)
	}
	if err = json.NewDecoder(file).Decode(&cfg); err != nil {
		return Config{}, fmt.Errorf("decode was error: %w", err)
	}
	return cfg, nil
}
