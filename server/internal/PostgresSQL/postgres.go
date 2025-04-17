package postgressql

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

var conn *pgx.Conn

type Config struct {
	Host     string `yaml:"POSTGRES_HOST"`
	Port     string `yaml:"POSTGRES_PORT"`
	Username string `yaml:"POSTGRES_USERNAME"`
	Password string `yaml:"POSTGRES_PASSWORD"`
	Database string `yaml:"POSTGRES_DATABASE"`
}

func New(config Config) (*pgx.Conn, error){
	var err error
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", config.Username, config.Password, config.Host, config.Port, config.Database)
	conn, err = pgx.Connect(context.Background(), connString)
	if err != nil{
		return nil, err
	}
	return conn, nil
}

