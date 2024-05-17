package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"time"
)

type (
	Config struct {
		Env     string  `yaml:"env" env:"ENV"`
		App     App     `yaml:"app"`
		Service Service `yaml:"service"`
		PG      PG      `yaml:"pg"`
	}

	App struct {
		Name string `yaml:"name" env:"APP_NAME"`
	}

	Service struct {
		OutputFilePath string        `yaml:"output_file_path" env:"SERVICE_OUTPUT_FILE_PATH"`
		UpdateInterval time.Duration `yaml:"update_inteval" env:"SERVICE_UPDATE_INTERVAL"`
		CancelInterval time.Duration `yaml:"cancel_inteval" env:"SERVICE_CANCEL_INTERVAL"`
		CsvSeparator   string        `yaml:"csv_separator" env:"SERVICE_CSV_SEPARATOR"`
	}

	PG struct {
		URL string `env:"PG_URL"`
	}
)

func NewConfig() (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig("./config/main.yml", cfg)
	if err != nil {
		return nil, err
	}

	err = cleanenv.ReadConfig(".env", cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
