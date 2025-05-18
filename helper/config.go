package helper

import "fmt"

type Config struct {
	Host     string   `json:"host"`
	Port     uint     `json:"port"`
	WikiRoot string   `json:"wiki-root"`
	Custom   []string `json:"custom"`
}

func GetConfig() (Config, error) {
	var cfg Config
	val, err := ParseJson("config.json", &cfg)
	if err != nil {
		return cfg, fmt.Errorf("unable to get config: %w", err)
	}

	cfgPtr, ok := val.(*Config)
	if !ok {
		return cfg, fmt.Errorf("unable to parse config")
	}

	cfg = *cfgPtr

	// might move this to main later
	if cfg.Port > 65535 {
		return cfg, fmt.Errorf("Invalid port number: %d", cfg.Port)
	}

	return cfg, nil
}
