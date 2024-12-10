package postgres

import (
	"embed"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"
	"log/slog"

	_ "github.com/lib/pq"
)

type DB struct {
	db  *sqlx.DB
	log *slog.Logger
}

//go:embed migrations/*.sql
var migrations embed.FS

func New(log *slog.Logger, dbPath string) (*DB, error) {
	const fn = "db.postgres.New"

	db, err := sqlx.Connect("postgres", dbPath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	log.Info("successfully connected to db",
		slog.String("fn", fn),
		slog.String("path", dbPath),
	)

	goose.SetBaseFS(migrations)

	if err = goose.SetDialect("postgres"); err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	if err = goose.Up(db.DB, "migrations"); err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	return &DB{
		db:  db,
		log: log,
	}, nil
}
