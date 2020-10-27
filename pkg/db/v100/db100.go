package db100

import (
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/concertLabs/oaf-server/internal/helpers"
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

type inserter interface {
	insertPG(query string) error
	insertOther(query string) error
}

func insertData(i inserter, query string) error {
	var err error
	if db.DriverName() == pgDriverName {
		err = i.insertPG(query)
	} else {
		err = i.insertOther(query)
	}
	return err
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

//Member models the m:n relation between Users and Sections
type Member struct {
	SectionID int `json:"sectionid" db:"SectionID"`
	UserID    int `json:"userid" db:"UserID"`
	Right     int `json:"right" db:"Right"`
}

//Organization holds infos about a single Org (like orchestra, football club you name it)
type Organization struct {
	OrganizationID int    `json:"organizationid" db:"OrganizationID"`
	Name           string `json:"name" db:"Name"`
	Picture        []byte `json:"picture" db:"Picture"`
}

func (o *Organization) insertPG(query string) error {
	query = query + ` RETURNING "OrganizationID"`
	tx := db.MustBegin()
	stmt, err := tx.Prepare(query)
	if err != nil {
		tx.Rollback()
		return errors.New("Error preparing Statement:" + err.Error())
	}
	stmt.QueryRow(o.Name, o.Picture).Scan(&o.OrganizationID)
	if err != nil {
		tx.Rollback()
		return errors.New("Error executing Statement:" + err.Error())
	}
	err = tx.Commit()
	if err != nil {
		return errors.New("Error executing Commit:" + err.Error())
	}
	return nil
}

func (o *Organization) insertOther(query string) error {
	res, err := db.Exec(query, o.Name, o.Picture)
	if err != nil {
		return errors.New("Error inserting Organiztaion: " + err.Error())
	}
	newid, err := res.LastInsertId()
	if err != nil {
		return errors.New("Error fetching new ID: " + err.Error())
	}
	o.OrganizationID = int(newid)
	return nil
}

//Insert inserts a new Organization into the database and adding the new OrganizationID into the struct
func (o *Organization) Insert() error {
	query := db.Rebind(`INSERT INTO "Organizations" ("Name", "Picture") VALUES (?, ?)`)
	err := insertData(o, query)
	if err != nil {
		return errors.New("Error inserting Organization:" + err.Error())
	}
	return nil
}

//GetOrganizations gives back all Organizations in the Database
func GetOrganizations() ([]Organization, error) {
	var o []Organization
	err := db.Select(&o, `SELECT * FROM "Organizations"`)
	if err != nil {
		return o, errors.New("Error getting Organizations:" + err.Error())
	}
	return o, nil
}

//GetDetails takes a Organization struct with only the OrganizationID and tries to fetch the remaining infos
func (o *Organization) GetDetails() error {
	query := db.Rebind(`SELECT * FROM "Organizations" WHERE "OrganizationID" = ? LIMIT 1`)
	err := db.Get(o, query, o.OrganizationID)
	if err != nil {
		return errors.New("Error getting Organization details:" + err.Error())
	}
	return nil
}

//Patch patches a Organization with new Info from a second struct
func (o *Organization) Patch(oo Organization) error {
	o.Name = helpers.CopyIfNotEmpty(o.Name, oo.Name)
	if len(oo.Picture) > 0 {
		o.Picture = oo.Picture
	}
	return nil
}

//Update updates all Organization Fields in the Database
func (o *Organization) Update() error {
	query := db.Rebind(`UPDATE "Organizations" SET "Name" = ?, "Picture" = ? WHERE "OrganizationID" = ?`)
	_, err := db.Exec(query, o.Name, o.Picture, o.OrganizationID)
	if err != nil {
		return errors.New("Error updating Organization:" + err.Error())
	}
	return nil
}

//DeleteOrganization deletes a organization with the given OrganizationID
func DeleteOrganization(id int) error {
	query := db.Rebind(`DELETE FROM "Organizations" WHERE "OrganizationID" = ?`)
	_, err := db.Exec(query, id)
	if err != nil {
		return errors.New("Error deleting Organization: " + err.Error())
	}
	return nil
}

//Section is a part of the Club (like Baritons or Handball Section)
type Section struct {
	SectionID      int    `json:"sectionid" db:"SectionID"`
	OrganizationID int    `json:"organizationid" db:"OrganizationID"`
	Name           string `json:"name" db:"Name"`
}

func (s *Section) insertPG(query string) error {
	query = query + ` RETURNING "SectionID"`
	tx := db.MustBegin()
	stmt, err := tx.Prepare(query)
	if err != nil {
		tx.Rollback()
		return errors.New("Error preparing Statement:" + err.Error())
	}
	stmt.QueryRow(s.OrganizationID, s.Name).Scan(&s.SectionID)
	if err != nil {
		tx.Rollback()
		return errors.New("Error executing Statement:" + err.Error())
	}
	err = tx.Commit()
	if err != nil {
		return errors.New("Error executing Commit:" + err.Error())
	}
	return nil
}

func (s *Section) insertOther(query string) error {
	res, err := db.Exec(query, s.OrganizationID, s.Name)
	if err != nil {
		return errors.New("Error inserting Section: " + err.Error())
	}
	newid, err := res.LastInsertId()
	if err != nil {
		return errors.New("Error fetching new ID: " + err.Error())
	}
	s.SectionID = int(newid)
	return nil
}

//Insert inserts a new Section into the database and adding the new SectionID into the struct
func (s *Section) Insert() error {
	query := db.Rebind(`INSERT INTO "Sections" ("OrganizationID", "Name") VALUES (?, ?)`)
	err := insertData(s, query)
	if err != nil {
		return errors.New("Error inserting Section:" + err.Error())
	}
	return nil
}

//GetSections gives back all Sections in the Database
func GetSections() ([]Section, error) {
	var s []Section
	err := db.Select(&s, `SELECT * FROM "Sections"`)
	if err != nil {
		return s, errors.New("Error getting Sections:" + err.Error())
	}
	return s, nil
}

//GetDetails takes a Section struct with only the SectionID and tries to fetch the remaining infos
func (s *Section) GetDetails() error {
	query := db.Rebind(`SELECT * FROM "Sections" WHERE "SectionID" = ? LIMIT 1`)
	err := db.Get(s, query, s.SectionID)
	if err != nil {
		return errors.New("Error getting Section details:" + err.Error())
	}
	return nil
}

//Patch patches a Section with new Info from a second struct
func (s *Section) Patch(ss Section) error {
	s.Name = helpers.CopyIfNotEmpty(s.Name, ss.Name)
	return nil
}

//Update updates all Section Fields in the Database
func (s *Section) Update() error {
	query := db.Rebind(`UPDATE "Sections" SET "Name" = ? WHERE "SectionID" = ?`)
	_, err := db.Exec(query, s.Name, s.SectionID)
	if err != nil {
		return errors.New("Error updating Sections:" + err.Error())
	}
	return nil
}

//DeleteSection deletes a Section with the given SectionID
func DeleteSection(id int) error {
	query := db.Rebind(`DELETE FROM "Sections" WHERE "SectionID" = ?`)
	_, err := db.Exec(query, id)
	if err != nil {
		return errors.New("Error deleting Section: " + err.Error())
	}
	return nil
}

//User is a singe User of the OAF
type User struct {
	UserID    int            `json:"userid" db:"UserID"`
	Username  string         `json:"username" db:"Username"`
	Password  string         `json:"password" db:"Password"`
	Salt      string         `json:"-" db:"Salt"`
	Shownname sql.NullString `json:"shownname" db:"Shownname"`
	EMail     string         `json:"email" db:"EMail"`
	SuperUser bool           `json:"superuser" db:"SuperUser"`
}

//DoesUserExist gives back a boolean if the User with this Username can be found in the Database
func DoesUserExist(username string) (bool, error) {
	var u User
	query := db.Rebind(`SELECT * FROM "Users" WHERE "Username" = ? LIMIT 1`)
	err := db.Get(&u, query, username)
	b := (u.UserID > 0)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return b, errors.New("Error finding User:" + err.Error())
	}
	return b, nil
}

//GetUsers gives back all Users in the Database
func GetUsers() ([]User, error) {
	var u []User
	err := db.Select(&u, `SELECT * FROM "Users"`)
	if err != nil {
		return u, errors.New("Error getting Users:" + err.Error())
	}
	return u, nil
}

//GetDetailstoUsername takes a User struct with only the Username and tries to fetch the remaining infos
func (u *User) GetDetailstoUsername() error {
	query := db.Rebind(`SELECT * FROM "Users" WHERE "Username" = ? LIMIT 1`)
	err := db.Get(u, query, u.Username)
	if err != nil {
		return errors.New("Error getting user details:" + err.Error())
	}
	return nil
}

//GetDetails takes a User struct with only the UserID and tries to fetch the remaining infos
func (u *User) GetDetails() error {
	query := db.Rebind(`SELECT * FROM "Users" WHERE "UserID" = ? LIMIT 1`)
	err := db.Get(u, query, u.UserID)
	if err != nil {
		return errors.New("Error getting user details:" + err.Error())
	}
	return nil
}

//Patch patches a User with new Info from a second struct
func (u *User) Patch(ou User) error {
	u.Username = helpers.CopyIfNotEmpty(u.Username, ou.Username)
	if ou.Password != "" {
		p, err := helpers.GeneratePasswordHash(ou.Password, u.Salt)
		if err != nil {
			return errors.New("Error hashing password: " + err.Error())
		}
		u.Password = p
	}
	u.EMail = helpers.CopyIfNotEmpty(u.EMail, ou.EMail)
	u.Shownname.String = helpers.CopyIfNotEmpty(u.Shownname.String, ou.Shownname.String)
	return nil
}

//Update updates all Userfields in the Database
func (u *User) Update() error {
	query := db.Rebind(`UPDATE "Users" SET "Username" = ?, "Password" = ?, "Salt" = ?, "Shownname" = ?, "EMail" = ? WHERE "UserID" = ?`)
	_, err := db.Exec(query, u.Username, u.Password, u.Salt, u.Shownname, u.EMail, u.UserID)
	if err != nil {
		return errors.New("Error updating user:" + err.Error())
	}
	return nil
}

func (u *User) insertPG(query string) error {
	query = query + ` RETURNING "UserID"`
	tx := db.MustBegin()
	stmt, err := tx.Prepare(query)
	if err != nil {
		tx.Rollback()
		return errors.New("Error preparing Statement:" + err.Error())
	}
	stmt.QueryRow(u.Username, u.Password, u.Salt, u.EMail, false).Scan(&u.UserID)
	if err != nil {
		tx.Rollback()
		return errors.New("Error executing Statement:" + err.Error())
	}
	err = tx.Commit()
	if err != nil {
		return errors.New("Error executing Commit:" + err.Error())
	}
	return nil
}

func (u *User) insertOther(query string) error {
	res, err := db.Exec(query, u.Username, u.Password, u.Salt, u.EMail, false)
	if err != nil {
		return errors.New("Error inserting User: " + err.Error())
	}
	newid, err := res.LastInsertId()
	if err != nil {
		return errors.New("Error fetching new ID: " + err.Error())
	}
	u.UserID = int(newid)
	return nil
}

//Insert inserts a new User into the database and adding the new UserID into the struct
func (u *User) Insert() error {
	query := db.Rebind(`INSERT INTO "Users" ("Username", "Password", "Salt", "EMail", "SuperUser") VALUES (?, ?, ?, ?, ?)`)
	err := insertData(u, query)
	if err != nil {
		return errors.New("Error inserting User:" + err.Error())
	}
	return nil
}

//DeleteUser delets user with the given UserID
func DeleteUser(id int) error {
	query := db.Rebind(`DELETE FROM "Users" WHERE "UserID" = ?`)
	_, err := db.Exec(query, id)
	if err != nil {
		return errors.New("Error deleting User: " + err.Error())
	}
	return nil
}
