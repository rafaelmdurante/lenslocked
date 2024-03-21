package models

import (
	"database/sql"
	"fmt"
	"io/fs"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/pressly/goose/v3"
)

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	SSLMode  string
}

func (cfg PostgresConfig) String() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database, cfg.SSLMode)
}

func DefaultPostgresConfig() PostgresConfig {
	return PostgresConfig{
		Host:     "localhost",
		Port:     "5432",
		User:     "baloo",
		Password: "junglebook",
		Database: "lenslocked",
		SSLMode:  "disable",
	}
}

// Open will open a SQL connection with the provided Postgres database.
// Callers of Open need to ensure the connection is eventually closed via the
// db.Close() method
func Open(config PostgresConfig) (*sql.DB, error) {
	db, err := sql.Open("pgx", config.String())

	if err != nil {
		return nil, fmt.Errorf("open: %w", err)
	}

	return db, nil
}

// Migrate runs the migrations, but it requires the .sql files to exist in the
// computer. That means binaries without embedded sql files will not work.
func Migrate(db *sql.DB, dir string) error {
    err := goose.SetDialect("postgres")
    if err != nil {
        return fmt.Errorf("migrate: failed to set dialect: %w", err)
    }

    err = goose.Up(db, dir)
    if err != nil {
        return fmt.Errorf("migrate: failed to run migration: %w", err)
    }

    return nil
}

// MigrateFS embeds migration files into the binary so it works without having
// the sql files available in the computer.
func MigrateFS(db *sql.DB, migrationsFS fs.FS, dir string) error {
    // in case the dir is an empty string, they probably meant the current dir
    // and goose wants a period for that
    if dir == "" {
        dir = "."
    }

    goose.SetBaseFS(migrationsFS)
    defer func() {
        // ensure that we remove the FS on the off chance some other of our app
        // uses goose for migrations and doesn't want to use our FS.
        goose.SetBaseFS(nil)
    }()

    return Migrate(db, dir)
}

