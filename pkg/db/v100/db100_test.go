package db100

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/concertLabs/oaf-server/internal/helpers"
	"github.com/concertLabs/oaf-server/pkg/config"
	"github.com/jmoiron/sqlx"
	migrate "github.com/rubenv/sql-migrate"
)

type testCase struct {
	connection config.DatabaseConnection
	migration  string
}

var testCases []testCase
var testUser = User{UserID: -1, Username: "test", Password: "test", EMail: "test@test", Shownname: sql.NullString{String: "Testiman", Valid: true}}
var testOrg = Organization{OrganizationID: -1, Name: "concertLabs", Picture: []byte("not actually a picture")}

func insertTestUser() {
	u := testUser
	s, _ := helpers.GenerateSalt()
	u.Salt = s
	pw, _ := helpers.GeneratePasswordHash(u.Password, u.Salt)
	u.Password = pw
	u.Insert()
}

func insertTestOrg() {
	o := testOrg
	o.Insert()
}

func execMigrations(t testCase, dir migrate.MigrationDirection) {
	migrations := &migrate.FileMigrationSource{
		Dir: t.migration,
	}
	migrate.SetTable("migrations")

	n, err := migrate.Exec(db.DB, t.connection.Driver, migrations, dir)
	if err != nil {
		log.Fatalln("Error applying migrations:", err)
	}
	if n == 0 {
		log.Fatalln("No migrations applied in DB:", t.connection.Driver)
	}
}

func setupDatabase(t testCase) {
	connector := t.connection.Driver
	if t.connection.Driver == "postgres" {
		connector = pgDriverName
	}
	var err error
	db, err = sqlx.Connect(connector, t.connection.Connection)
	if err != nil {
		log.Fatalln("Error connecting to database:", err)
	}
	execMigrations(t, migrate.Up)

}

func teardownDatabase(t testCase) {
	execMigrations(t, migrate.Down)
}

func TestMain(m *testing.M) {
	sqliteCase := testCase{}
	sqliteCase.connection.Driver = "sqlite3"
	sqliteCase.connection.Connection = "test.db"
	sqliteCase.migration = "../../../scripts/sql/migrations/sqlite"
	testCases = append(testCases, sqliteCase)
	os.Remove(sqliteCase.connection.Connection)

	pgCase := testCase{}
	pgCase.connection.Driver = "postgres"
	pgCase.connection.Connection = "postgres://oaf-server:oaf-server@localhost:5432/oaf-server"
	pgCase.migration = "../../../scripts/sql/migrations/postgres"
	testCases = append(testCases, pgCase)

	os.Exit(m.Run())
}

func TestUserInsert(t *testing.T) {

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Insert User in %s", tc.connection.Driver), func(t *testing.T) {
			setupDatabase(tc)
			u := testUser
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
			teardownDatabase(tc)
		})
	}

}

func TestDoesUserExist(t *testing.T) {
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Does User Exist in %s", tc.connection.Driver), func(t *testing.T) {
			setupDatabase(tc)
			cont, err := DoesUserExist("foobar")

			if err != nil {
				t.Fatalf("No error expected but got %v", err)
			}

			if cont {
				t.Errorf("Expected false got %v", cont)
			}

			insertTestUser()

			cont, err = DoesUserExist(testUser.Username)

			if err != nil {
				t.Fatalf("No error expected but got %v", err)
			}

			if !cont {
				t.Errorf("Expected true got %v", cont)
			}
			teardownDatabase(tc)
		})
	}
}

func TestGetUsers(t *testing.T) {
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Get Users in %s", tc.connection.Driver), func(t *testing.T) {
			setupDatabase(tc)
			insertTestUser()
			uu, err := GetUsers()
			if err != nil {
				t.Fatalf("No error expected but got %v", err)
			}
			if len(uu) != 1 {
				t.Error("Expected length 1 got ", len(uu))
			}
			teardownDatabase(tc)
		})
	}
}

func TestGetDetailstoUsername(t *testing.T) {
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Get Details to username in %s", tc.connection.Driver), func(t *testing.T) {
			setupDatabase(tc)
			insertTestUser()
			u := User{Username: testUser.Username}
			err := u.GetDetailstoUsername()
			if err != nil {
				t.Fatalf("No error expected but got %v", err)
			}
			if u.UserID != 1 {
				t.Errorf("Expected User_ID 1 but got %v", u.UserID)
			}
			teardownDatabase(tc)
		})
	}
}

func TestGetUserDetails(t *testing.T) {
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Get User Details in %s", tc.connection.Driver), func(t *testing.T) {
			setupDatabase(tc)
			insertTestUser()
			u := User{UserID: 1}
			err := u.GetDetails()
			if err != nil {
				t.Fatalf("No error expected but got %v", err)
			}
			if u.Username != testUser.Username {
				t.Errorf("Expected Username test but got %v", u.Username)
			}
			teardownDatabase(tc)
		})
	}
}

func TestPatchUser(t *testing.T) {
	u := testUser
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
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Update User in %s", tc.connection.Driver), func(t *testing.T) {
			setupDatabase(tc)
			insertTestUser()
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
			teardownDatabase(tc)
		})
	}
}

func TestDeleteUser(t *testing.T) {
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Delete User in %s", tc.connection.Driver), func(t *testing.T) {
			setupDatabase(tc)
			insertTestUser()
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
			teardownDatabase(tc)
		})
	}
}

func TestOrganizationInsert(t *testing.T) {
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Insert Organization in %s", tc.connection.Driver), func(t *testing.T) {
			setupDatabase(tc)
			o := testOrg

			err := o.Insert()

			if err != nil {
				t.Errorf("Expected no error but got %v", err)
			}

			if o.OrganizationID < 0 {
				t.Errorf("Expected OrganizationID > 0 but got %v", o.OrganizationID)
			}
			teardownDatabase(tc)
		})
	}
}

func TestGetOrganizations(t *testing.T) {
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Get Organization in %s", tc.connection.Driver), func(t *testing.T) {
			setupDatabase(tc)
			insertTestOrg()
			oo, err := GetOrganizations()
			if err != nil {
				t.Fatalf("No error expected but got %v", err)
			}
			if len(oo) != 1 {
				t.Error("Expected length 1 got ", len(oo))
			}
			teardownDatabase(tc)
		})
	}
}

func TestOrganizationDetails(t *testing.T) {
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Get Organiztaion Details in %s", tc.connection.Driver), func(t *testing.T) {
			setupDatabase(tc)
			insertTestOrg()
			o := Organization{OrganizationID: 1}
			err := o.GetDetails()
			if err != nil {
				t.Fatalf("No error expected but got %v", err)
			}
			if o.Name != testOrg.Name {
				t.Errorf("Expected Name test but got %v", o.Name)
			}
			if len(o.Picture) != len(testOrg.Picture) {
				t.Errorf("Length of byte array different %v vs %v", len(o.Picture), len(testOrg.Picture))
			}
			teardownDatabase(tc)
		})
	}
}

func TestPatchOrganization(t *testing.T) {
	o := Organization{Name: "concertLabs", Picture: []byte("not actually a picture")}
	on := Organization{Name: "concertLabs2", Picture: []byte("still not a picture")}

	err := o.Patch(on)

	if err != nil {
		t.Fatalf("No error expected but got %v", err)
	}

	if o.Name != on.Name {
		t.Errorf("Expected name to be %v but got %v", on.Name, o.Name)
	}
	if len(o.Picture) != len([]byte("still not a picture")) {
		t.Errorf("Length of byte array different %v vs %v", len(o.Picture), len([]byte("still not a picture")))
	}
}

func TestUpdateOrganization(t *testing.T) {
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Update Organization in %s", tc.connection.Driver), func(t *testing.T) {
			setupDatabase(tc)
			insertTestOrg()
			o := Organization{OrganizationID: 1, Name: "concertLabs2", Picture: []byte("still not a picture")}
			on := Organization{OrganizationID: 1}
			err := o.Update()
			if err != nil {
				t.Fatalf("No error expected but got %v", err)
			}
			err = on.GetDetails()
			if err != nil {
				t.Fatalf("No error expected but got %v", err)
			}
			if o.Name != on.Name {
				t.Errorf("Found differences between old and updated name: %v %v", o.Name, on.Name)
			}
			if len(o.Picture) != len(on.Picture) {
				t.Errorf("Length of byte array different %v vs %v", len(o.Picture), len(on.Picture))
			}
			teardownDatabase(tc)
		})
	}
}

func TestDeleteOrganization(t *testing.T) {
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Delete Organization in %s", tc.connection.Driver), func(t *testing.T) {
			setupDatabase(tc)
			insertTestOrg()
			err := DeleteOrganization(1)
			if err != nil {
				t.Fatalf("No error expected but got %v", err)
			}
			oo, err := GetOrganizations()
			if err != nil {
				t.Fatalf("No error expected but got %v", err)
			}
			if len(oo) != 0 {
				t.Error("Expected length 0 got ", len(oo))
			}
			teardownDatabase(tc)
		})
	}
}
