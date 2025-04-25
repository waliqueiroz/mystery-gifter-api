package postgres

//go:generate go run go.uber.org/mock/mockgen -destination mock_postgres/db.go . DB
//go:generate go run go.uber.org/mock/mockgen -destination mock_postgres/tx.go . TX

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/waliqueiroz/mystery-gifter-api/internal/infra/config"
)

const POSTGRES_UNIQUE_VIOLATION = "unique_violation"

type QueryExecutor interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	GetContext(ctx context.Context, dest any, query string, args ...any) error
	SelectContext(ctx context.Context, dest any, query string, args ...any) error
}

type TX interface {
	QueryExecutor
	Commit() error
	Rollback() error
}

type DB interface {
	QueryExecutor
	Close() error
	BeginTxx(ctx context.Context, opts *sql.TxOptions) (TX, error)
	GetDB() *sql.DB
}

// wrapper para sqlx.DB implementar a interface DB
type sqlxDB struct {
	*sqlx.DB
}

func NewDB(db *sqlx.DB) DB {
	return &sqlxDB{db}
}

func (db *sqlxDB) BeginTxx(ctx context.Context, opts *sql.TxOptions) (TX, error) {
	tx, err := db.DB.BeginTxx(ctx, opts)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (db *sqlxDB) GetDB() *sql.DB {
	return db.DB.DB
}

func Connect(databaseConfig config.DatabaseConfig) (DB, error) {
	dataSourceName := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=UTC",
		databaseConfig.Host,
		databaseConfig.Port,
		databaseConfig.Username,
		databaseConfig.Password,
		databaseConfig.Database,
	)

	db, err := sqlx.Connect("postgres", dataSourceName)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(50)
	db.SetMaxIdleConns(50)
	db.SetConnMaxLifetime(10 * time.Minute)

	return NewDB(db), nil
}

func Migrate(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance("file://./internal/infra/outgoing/postgres/migrations", "postgres", driver)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Println("Nothing to migrate.")
			return nil
		}

		return err
	}

	return nil
}
