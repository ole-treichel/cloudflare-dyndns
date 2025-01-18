package config

import (
	"encoding/json"
	"io"
	"os"
)

type DomainConfig struct {
	ZoneId  string   `json:"zoneId"`
	Domains []string `json:"domains"`
}

type Config struct {
	DomainConfigs []DomainConfig `json:"domainConfigs"`
}

func GetConfig() (Config, error) {
	var config Config
	configFile, err := os.Open("config.json")

	if err != nil {
		return config, err
	}

	defer configFile.Close()

	configBytes, err := io.ReadAll(configFile)

	if err != nil {
		return config, err
	}

	err = json.Unmarshal(configBytes, &config)

	if err != nil {
		return config, err
	}

	return config, nil
}
