package smtp

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
)

type SMTPConfig struct {
	SMTPHost     string `yaml:"smtpHost"`
	SMTPPort     string `yaml:"smtpPort"`
	SMTPUser     string `yaml:"username"`
	ClientID     string `yaml:"clientID"`
	ClientSecret string `yaml:"secretKey"`
}

type Config struct {
	SMTP SMTPConfig `yaml:"SMTP"`
}

func ParseFromYaml() (*Config, error) {
	cfgPath := "smtp/config.yaml"
	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file not found: %s", cfgPath)
	}
	var cfg Config
	if err := cleanenv.ReadConfig(cfgPath, &cfg); err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}
	return &cfg, nil
}
