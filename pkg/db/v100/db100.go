package db100

import (
	"errors"
	"log"

	"github.com/concertLabs/oaf-server/internal/packer"
	"github.com/concertLabs/oaf-server/pkg/config"
	"github.com/gobuffalo/packr"
	"github.com/jmoiron/sqlx"
	migrate "github.com/rubenv/sql-migrate"

	//Postgres Driver
	_ "github.com/jackc/pgx/stdlib"

	//SQLite Driver
	_ "github.com/mattn/go-sqlite3"
)

const (
	pgDriverName = "pgx"
)

var db *sqlx.DB

//Initialisation sets up the DB connection and applies the lates migrations
func Initialisation(dbc config.DatabaseConnection, migdir string) error {
	var err error
	connector := dbc.Driver
	if dbc.Driver == "postgres" {
		connector = pgDriverName
	}
	db, err = sqlx.Connect(connector, dbc.Connection)
	if err != nil {
		return errors.New("Error connecting to database:" + err.Error())
	}

	var migbox packr.Box

	switch dbc.Driver {
	case "sqlite3":
		migbox = packer.SQLiteBox
		log.Println("SQLite should only been used for local testing environments. You have been warned")
	case "postgres":
		migbox = packer.PostgresBox
	}

	var migrations migrate.MigrationSource
	if migdir == "" {
		migrations = &migrate.PackrMigrationSource{
			Box: migbox,
			Dir: "./",
		}
	} else {
		migrations = &migrate.FileMigrationSource{
			Dir: migdir,
		}
	}

	migrate.SetTable("migrations")

	n, err := migrate.Exec(db.DB, dbc.Driver, migrations, migrate.Up)
	if err != nil {
		return errors.New("Error applying migrations:" + err.Error())
	}
	log.Println("Applied ", n, "Migrations")
	return nil
}
