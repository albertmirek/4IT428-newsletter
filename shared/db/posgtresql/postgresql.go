package postgresql

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/reflectx"
)

type Config struct {
	Host     string `envconfig:"POSTGRESQL_HOST" default:"localhost"`
	Port     int    `envconfig:"POSTGRESQL_PORT" default:"5432"`
	User     string `envconfig:"POSTGRESQL_USER" default:"user"`
	Password string `envconfig:"POSTGRESQL_PASSWORD" default:"password"`
	DBName   string `envconfig:"POSTGRESQL_DBNAME" default:"newsletter_app"`
	SSLMode  string `envconfig:"POSTGRESQL_SSLMODE" default:"disable"`
}

func NewConnectionPool(ctx context.Context, cfg Config) (*sqlx.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode)

	config, err := pgx.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	connStr := stdlib.RegisterConnConfig(config)
	db := sqlx.MustConnect("pgx", connStr)
	db.Mapper = reflectx.NewMapperFunc("db", func(s string) string { return s })

	return db, nil
}
