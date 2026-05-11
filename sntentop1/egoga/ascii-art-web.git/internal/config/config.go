package config

import (
	"encoding/json"
	"io"
	"os"
)

type Config struct {
	Banner_base_path string   `json:"banner_base_path"`
	ValidBanners     []string `json:"valid_banners"`
}

func LoadConfig(path string) (*Config, error) {
	configFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer configFile.Close()

	configContent, err := io.ReadAll(configFile)
	if err != nil {
		return nil, err
	}

	config := Config{}

	err = json.Unmarshal(configContent, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
