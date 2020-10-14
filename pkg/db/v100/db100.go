package db100

import (
	"database/sql"
	"log"
	"time"

	"github.com/concertLabs/oaf-server/internal/packer"
	"github.com/concertLabs/oaf-server/pkg/config"
	"github.com/gobuffalo/packr"
	"github.com/jmoiron/sqlx"
	migrate "github.com/rubenv/sql-migrate"

	//Postgres Driver
	_ "github.com/lib/pq"

	//SQLite Driver
	_ "github.com/mattn/go-sqlite3"
)

var db *sqlx.DB

//Initialisation sets up the DB connection and applies the lates migrations
func Initialisation(dbc config.DatabaseConnection) {
	var err error
	db, err = sqlx.Connect(dbc.Driver, dbc.Connection)
	if err != nil {
		log.Fatalln(err)
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
		log.Fatalln(err)
	}
	log.Println("Applied ", n, "Migrations")
}

//Attendee manages Useres that attend a single Event
type Attendee struct {
	EventID    int            `json:"eventid"`
	UserID     int            `json:"userid"`
	Commitment int            `json:"commitment"`
	Comment    sql.NullString `json:"comment"`
}

//Comment holds User message for a Event "Blackboard"
type Comment struct {
	CommentID int    `json:"commentid"`
	EventID   int    `json:"eventid"`
	UserID    int    `json:"userid"`
	Comment   string `json:"comment"`
}

//Event is a single Happening (Gig, Game, etc)
type Event struct {
	EventID        int            `json:"eventid"`
	OrganizationID int            `json:"organizationid"`
	Name           string         `json:"name"`
	Address        sql.NullString `json:"address"`
	Start          time.Time      `json:"start"`
	End            sql.NullTime   `json:"end"`
	Creator        int            `json:"creator"`
}

//Info hold Data from the Info table which is a general key value store
type Info struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

//Member models the m:n relation between Users and Sections
type Member struct {
	SectionID int `json:"sectionid"`
	UserID    int `json:"userid"`
	Right     int `json:"right"`
}

//Organization holds infos about a single Org (like orchestra, football club you name it)
type Organization struct {
	OrganizationID int    `json:"organizationid"`
	Name           string `json:"name"`
	Picture        []byte `json:"picture"`
}

//Section is a part of the Club (like Baritons or Handball Section)
type Section struct {
	SectionID      int    `json:"sectionid"`
	OrganizationID int    `json:"organizationid"`
	Name           string `json:"name"`
}

//User is a singe User of the OAF
type User struct {
	UserID    int            `json:"userid"`
	Username  string         `json:"username"`
	Password  string         `json:"password"`
	Salt      string         `json:"-"`
	Shownname sql.NullString `json:"shownname"`
	EMail     string         `json:"email"`
	SuperUser bool           `json:"superuser"`
}
