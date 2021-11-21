package api100

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/concertLabs/oaf-server/internal/helpers"
	"github.com/concertLabs/oaf-server/pkg/config"
	db100 "github.com/concertLabs/oaf-server/pkg/db/v100"
)

var testUser = db100.User{UserID: 1, Username: "mosch", Password: "ernst", EMail: "ernstmosch@conecertLabs.solutions", Shownname: sql.NullString{String: "Ernst Mosch", Valid: true}, SuperUser: true}

func TestMain(m *testing.M) {
	var conf config.DatabaseConnection
	conf.Connection = "test.db"
	conf.Driver = "sqlite3"
	db100.Initialisation(conf, "../../../scripts/sql/migrations/sqlite")
	s, _ := helpers.GenerateSalt()
	testUser.Salt = s
	pw, _ := helpers.GeneratePasswordHash(testUser.Password, testUser.Salt)
	testUser.Password = pw
	err := testUser.Insert()
	if err != nil {
		log.Fatalf("Error inserting testuser: %v", err)
	}
	os.Exit(m.Run())
}

func TestGetAuthorization(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal("Error creating request:", err)
	}
	req.Header.Set("Authorization", "wrong")
	_, _, err = getAuthorization(req)
	if err == nil {
		t.Errorf("Expected error but got non")
	}
	req.Header.Set("Authorization", "Bearer Token")
	s1, s2, err := getAuthorization(req)
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
	if s1 != "Bearer" {
		t.Errorf("Expected s1 to be Bearer but got %v", s1)
	}
	if s2 != "Token" {
		t.Errorf("Expected s2 to be Token but got %v", s2)
	}
}
