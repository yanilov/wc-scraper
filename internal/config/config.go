package config

import (
	"os"

	"github.com/yanilov/wc-scraper/internal/scraper"
	"gopkg.in/yaml.v3"
)

type AppConfig struct {
	ScraperConfig   scraper.ScraperConfig         `yaml:"scraper"`
	JobLoaderConfig scraper.ScrapeJobLoaderConfig `yaml:"job_loader"`
	WordBankUrl     string                        `yaml:"word_bank_url"`
	JobsUrl         string                        `yaml:"jobs_url"`
}

func FromFile(path string) (*AppConfig, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var result AppConfig
	err = yaml.Unmarshal(file, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
