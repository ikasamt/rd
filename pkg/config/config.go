package config

import (
	"fmt"
	"os"
)

type Config struct {
	RedmineURL string
	APIKey     string
}

func Load(urlFlag, keyFlag string) (*Config, error) {
	cfg := &Config{}

	// URLの取得（フラグ優先、次に環境変数）
	if urlFlag != "" {
		cfg.RedmineURL = urlFlag
	} else {
		cfg.RedmineURL = os.Getenv("REDMINE_URL")
	}

	// APIキーの取得（フラグ優先、次に環境変数）
	if keyFlag != "" {
		cfg.APIKey = keyFlag
	} else {
		cfg.APIKey = os.Getenv("REDMINE_API_KEY")
	}

	// 検証
	if cfg.RedmineURL == "" {
		return nil, fmt.Errorf("Redmine URL is not set. Please set REDMINE_URL environment variable or use --url flag")
	}

	if cfg.APIKey == "" {
		return nil, fmt.Errorf("Redmine API key is not set. Please set REDMINE_API_KEY environment variable or use --key flag")
	}

	return cfg, nil
}