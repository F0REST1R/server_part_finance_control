package config

import (
	postgressql "Server_part_finance_control/server/internal/PostgresSQL"
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	POSTGRES postgressql.Config `yaml:"POSTGRES"`
}

func New() (*Config, error) {
	var cfg Config
	if err := cleanenv.ReadConfig("C:/Users/alexa/Desktop/IT STUDY/GO/Golang_project/Server_part_finance_control/config.yaml", &cfg);
	err != nil {
		return nil, err
	}
	fmt.Printf("\nУспешная загрузка конфигураций %+v\n", cfg)
	return &cfg, nil
}