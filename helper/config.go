package helper

import (
	"log"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
	Dir  string `mapstructure:"dir"`
}

func LoadConfig() Config {
	v := viper.New()
	v.SetConfigName("config")
	v.AddConfigPath(".")
	v.AddConfigPath("./configs")
	v.SetDefault("host", "0.0.0.0")
	v.SetDefault("port", 8080)
	v.SetDefault("dir", "")
	v.AutomaticEnv()
	v.SetEnvPrefix("APP")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("Config file not found. Using defaults and environment variables.")
		} else {
			log.Fatalf("Fatal error reading config file: %s", err)
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		log.Fatalf("Unable to unmarshal config into ServerConfig struct: %s", err)
	}

	if cfg.Host == "" {
		log.Fatal("Configuration error: 'host' cannot be empty.")
	}
	if cfg.Port <= 0 {
		log.Fatal("Configuration error: 'port' must be a positive integer.")
	}
	if cfg.Dir == "" {
		log.Fatal("Configuration error: 'dir' cannot be empty.")
	}

	if stat, err := os.Stat(cfg.Dir); os.IsNotExist(err) {
		log.Fatalf("Configured directory '%s' does not exist. Please create it or set a valid path.", cfg.Dir)
	} else if !stat.IsDir() {
		log.Fatalf("Configured path '%s' is not a directory.", cfg.Dir)
	}

	return cfg
}
