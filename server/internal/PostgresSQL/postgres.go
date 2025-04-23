package postgressql

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

type PostgresSQL struct{
	conn *pgx.Conn
}

type Config struct {
	Host     string `yaml:"POSTGRES_HOST"`
	Port     string `yaml:"POSTGRES_PORT"`
	Username string `yaml:"POSTGRES_USERNAME"`
	Password string `yaml:"POSTGRES_PASSWORD"`
	Database string `yaml:"POSTGRES_DATABASE"`
}

func (c *PostgresSQL) New(config Config, ctx context.Context) (*PostgresSQL, error){
	var err error
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", config.Username, config.Password, config.Host, config.Port, config.Database)
	c.conn, err = pgx.Connect(context.Background(), connString)
	if err != nil{
		return nil, err
	}

	ctxPing, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	if err := c.conn.Ping(ctxPing); err != nil{
		return nil, fmt.Errorf("failed to ping databse: %w", err)
	}

	return &PostgresSQL{conn: c.conn}, nil
}

func (p *PostgresSQL) Close() error {
	if p.conn != nil{
		return p.conn.Close(context.Background())
	}
	return nil
}

func (p *PostgresSQL) HealtCheck(ctx context.Context) error{
	if p.conn == nil{
		return fmt.Errorf("connection is not initialized")
	}

	return p.conn.Ping(ctx)
}

