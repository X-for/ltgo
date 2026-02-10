package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	Cookie   string `json:"cookie"`
	Language string `json:"language"`
	Site     string `json:"site"`
}

func getConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".ltgo", "config.json"), nil
}

func Load() (*Config, error) {
	home, _ := os.UserHomeDir()
	configPath := filepath.Join(home, ".ltgo", "config.json") // 保持 config.json

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	// 设置默认语言 (兼容旧配置文件)
	if cfg.Language == "" {
		cfg.Language = "golang"
	}

	return &cfg, nil
}

func (c *Config) Save() error {
	path, err := getConfigPath()
	if err != nil {
		return err
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(c, "", " ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}
