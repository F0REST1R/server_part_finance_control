package postgressql

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

type PostgresSQL struct{
	Conn *pgx.Conn
}

type Config struct {
	Host     string `yaml:"POSTGRES_HOST"`
	Port     string `yaml:"POSTGRES_PORT"`
	Username string `yaml:"POSTGRES_USERNAME"`
	Password string `yaml:"POSTGRES_PASSWORD"`
	Database string `yaml:"POSTGRES_DATABASE"`
}

func NewPost(config Config, ctx context.Context) (*PostgresSQL, error){	
	connString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
	)

	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}

	ctxPing, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	if err := conn.Ping(ctxPing); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &PostgresSQL{Conn: conn}, nil
}

func (p *PostgresSQL) Close() error {
	if p.Conn != nil{
		return p.Conn.Close(context.Background())
	}
	return nil
}

func (p *PostgresSQL) HealtCheck(ctx context.Context) error{
	if p.Conn == nil{
		return fmt.Errorf("connection is not initialized")
	}

	return p.Conn.Ping(ctx)
}

