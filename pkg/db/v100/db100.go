package db100

import (
	"database/sql"
	"errors"
	"log"
	"time"

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
func Initialisation(dbc config.DatabaseConnection) error {
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

	migrations := &migrate.PackrMigrationSource{
		Box: migbox,
		Dir: "./",
	}
	migrate.SetTable("migrations")

	n, err := migrate.Exec(db.DB, dbc.Driver, migrations, migrate.Up)
	if err != nil {
		return errors.New("Error applying migrations:" + err.Error())
	}
	log.Println("Applied ", n, "Migrations")
	return nil
}

//Attendee manages Useres that attend a single Event
type Attendee struct {
	EventID    int            `json:"eventid" db:"EventID"`
	UserID     int            `json:"userid" db:"UserID"`
	Commitment int            `json:"commitment" db:"Commitment"`
	Comment    sql.NullString `json:"comment" db:"Comment"`
}

//Comment holds User message for a Event "Blackboard"
type Comment struct {
	CommentID int    `json:"commentid" db:"CommentID"`
	EventID   int    `json:"eventid" db:"EventID"`
	UserID    int    `json:"userid" db:"UserID"`
	Comment   string `json:"comment"  db:"Comment"`
}

//Event is a single Happening (Gig, Game, etc)
type Event struct {
	EventID        int            `json:"eventid" db:"EventID"`
	OrganizationID int            `json:"organizationid" db:"OrganizationID"`
	Name           string         `json:"name" db:"Name"`
	Address        sql.NullString `json:"address" db:"Address"`
	Start          time.Time      `json:"start" db:"Start"`
	End            sql.NullTime   `json:"end" db:"End"`
	Creator        int            `json:"creator" db:"Creator"`
}

//Info hold Data from the Info table which is a general key value store
type Info struct {
	Key   string `json:"key" db:"Key"`
	Value string `json:"value" db:"Value"`
}
