package config

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"os"
)

type Config struct {
	DbDSN    string `json:"db_dsn"`
	LogLevel string `json:"log_level"`
	Listen   string `json:"listen"`
}

func LoadConfig(path string) (*Config, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	v := Config{
		DbDSN:    "",
		LogLevel: "info",
	}
	err = json.Unmarshal(b, &v)
	if err != nil {
		return nil, err
	}

	lvl, err := log.ParseLevel(v.LogLevel)
	if err != nil {
		return &v, err
	}
	log.SetLevel(lvl)

	log.WithFields(log.Fields{
		"path": path,
	}).Debug("configuration loaded")

	return &v, nil
}
