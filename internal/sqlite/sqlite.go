package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type DB struct {
	Connection *sql.DB
	ctx        context.Context // background context
	cancel     func()          // cancel background context

	// Placeholder for queries from sqlc

	// Datasource name.
	DSN string

	// Returns the current time. Defaults to time.Now().
	// Can be mocked for tests.
	Now func() time.Time
}

func NewDB(dsn string) *DB {
	db := &DB{
		DSN: dsn,
		Now: time.Now,
	}
	db.ctx, db.cancel = context.WithCancel(context.Background())
	return db

}

func (db *DB) Open() (err error) {

	if db.DSN == "" {
		return fmt.Errorf("DSN required")
	}
	if db.Connection, err = sql.Open("sqlite3", db.DSN); err != nil {
		return err
	}

	// Enable WAL. SQLite performs better with the WAL  because it allows
	// multiple readers to operate while data is being written.
	if _, err := db.Connection.Exec(`PRAGMA journal_mode = wal;`); err != nil {
		return fmt.Errorf("enable wal: %w", err)
	}

	// Enable foreign key checks. For historical reasons, SQLite does not check
	// foreign key constraints by default... which is kinda insane. There's some
	// overhead on inserts to verify foreign key integrity but it's definitely
	// worth it.
	if _, err := db.Connection.Exec(`PRAGMA foreign_keys = ON;`); err != nil {
		return fmt.Errorf("foreign keys pragma: %w", err)
	}
	return nil
}

func (db *DB) Close() error {
	// Cancel background context.
	db.cancel()

	// Close database.
	if db.Connection != nil {
		return db.Connection.Close()
	}
	return nil
}
