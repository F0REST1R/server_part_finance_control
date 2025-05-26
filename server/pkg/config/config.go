package config

import (
	postgresql "Server_part_finance_control/server/internal/PostgresSQL"
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)
type gRPC_cfg struct{
	PORT int `yaml:"PORT"`
	TIMEOUT string `yaml:"TIMEOUT"`
}

type Config struct {
	ENV string `yaml:"ENV"`
	POSTGRES postgresql.Config `yaml:"POSTGRES"`
	GRPC gRPC_cfg `yaml:"GRPC"`
}

func New() (*Config, error) {
	var cfg Config

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == ""{
		configPath = "./config.yaml"
	}

	if err := cleanenv.ReadConfig(configPath, &cfg);
	err != nil {
		return nil, err
	}
	fmt.Printf("\nУспешная загрузка конфигураций %+v\n", cfg)
	return &cfg, nil
}