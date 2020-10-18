package db100

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/concertLabs/oaf-server/internal/helpers"
	"github.com/concertLabs/oaf-server/pkg/config"
	migrate "github.com/rubenv/sql-migrate"
)

func TestMain(m *testing.M) {
	var con config.DatabaseConnection
	con.Driver = "sqlite3"
	con.Connection = "test.db"
	os.Remove(con.Connection)
	abs, _ := filepath.Abs(con.Connection)
	log.Println("Test Database Path:", abs)
	Initialisation(con)
	migrations := &migrate.FileMigrationSource{
		Dir: "../../../scripts/sql/migrations/sqlite",
	}
	migrate.SetTable("migrations")

	n, err := migrate.Exec(db.DB, con.Driver, migrations, migrate.Up)
	if err != nil {
		log.Fatalln("Error applying migrations:", err)
	}
	log.Println("Applied ", n, "Migrations")
	exit := m.Run()
	os.Exit(exit)
}

func TestUserInsert(t *testing.T) {

	u := User{Username: "test", Password: "test", EMail: "test@test", Shownname: sql.NullString{String: "Testiman", Valid: true}}
	s, err := helpers.GenerateSalt()
	if err != nil {
		t.Fatalf("No error expected but got %v", err)
	}
	u.Salt = s

	pw, err := helpers.GeneratePasswordHash(u.Password, u.Salt)
	if err != nil {
		t.Fatalf("No error expected but got %v", err)
	}
	u.Password = pw

	err = u.Insert()

	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}

	if u.UserID < 0 {
		t.Errorf("Expected Userid > 0 but got %v", u.UserID)
	}
}

func TestDoesUserExist(t *testing.T) {
	cont, err := DoesUserExist("foobar")

	if err != nil {
		t.Fatalf("No error expected but got %v", err)
	}

	if cont {
		t.Errorf("Expected false got %v", cont)
	}

	cont, err = DoesUserExist("test")

	if err != nil {
		t.Fatalf("No error expected but got %v", err)
	}

	if !cont {
		t.Errorf("Expected true got %v", cont)
	}
}

func TestGetUsers(t *testing.T) {
	uu, err := GetUsers()
	if err != nil {
		t.Fatalf("No error expected but got %v", err)
	}
	if len(uu) != 1 {
		t.Error("Expected length 1 got ", len(uu))
	}
}

func TestGetDetailstoUsername(t *testing.T) {
	u := User{Username: "test"}
	err := u.GetDetailstoUsername()
	if err != nil {
		t.Fatalf("No error expected but got %v", err)
	}
	if u.UserID != 1 {
		t.Errorf("Expected User_ID 1 but got %v", u.UserID)
	}
}

func TestGetUserDetails(t *testing.T) {
	u := User{UserID: 1}
	err := u.GetDetails()
	if err != nil {
		t.Fatalf("No error expected but got %v", err)
	}
	if u.Username != "test" {
		t.Errorf("Expected Username test but got %v", u.Username)
	}
}

func TestPatchUser(t *testing.T) {
	u := User{Username: "test", Password: "test", EMail: "test@test", Shownname: sql.NullString{String: "Testiman", Valid: true}}
	un := User{Username: "test", Password: "test", EMail: "test@otherhost", Shownname: sql.NullString{String: "Testiman2", Valid: true}}

	err := u.Patch(un)

	if err != nil {
		t.Fatalf("No error expected but got %v", err)
	}

	if u.EMail != "test@otherhost" {
		t.Errorf("Expected Email admin@otherhost but got %v", u.EMail)
	}

	if u.Shownname.String != "Testiman2" {
		t.Errorf("Expected Shownname Testiman2 but got %v", u.Shownname.String)
	}
}

func TestUpdateUser(t *testing.T) {
	u := User{UserID: 1, Username: "test1", Password: "test", EMail: "test@otherhost", Shownname: sql.NullString{String: "Testiman2", Valid: true}}
	un := User{UserID: 1}
	err := u.Update()
	if err != nil {
		t.Fatalf("No error expected but got %v", err)
	}
	err = un.GetDetails()
	if err != nil {
		t.Fatalf("No error expected but got %v", err)
	}
	if u.Username != un.Username {
		t.Errorf("Found differences between old and updated username: %v %v", u.Username, un.Username)
	}
	if u.EMail != "test@otherhost" {
		t.Errorf("Expected Email admin@otherhost but got %v", u.EMail)
	}
	if u.Shownname.String != "Testiman2" {
		t.Errorf("Expected Shownname Testiman2 but got %v", u.Shownname.String)
	}
}

func TestDeleteUser(t *testing.T) {
	err := DeleteUser(1)
	if err != nil {
		t.Fatalf("No error expected but got %v", err)
	}
	uu, err := GetUsers()
	if err != nil {
		t.Fatalf("No error expected but got %v", err)
	}
	if len(uu) != 0 {
		t.Error("Expected length 0 got ", len(uu))
	}
}
