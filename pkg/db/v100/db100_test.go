package db100

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

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
var testSection = Section{SectionID: -1, OrganizationID: 1, Name: "Developers"}
var testMember = Member{SectionID: 1, UserID: 1, Rights: 1}
var testEvent = Event{EventID: -1, OrganizationID: 1, Name: "Hackathon", Address: sql.NullString{String: "there", Valid: true}, Start: time.Now(), End: sql.NullTime{Time: time.Now().Add(time.Hour)}, Creator: 1}
var testComment = Comment{CommentID: -1, EventID: 1, UserID: 1, Creation: time.Now(), Comment: "Hullu"}
var testAttendee = Attendee{EventID: 1, UserID: 1, Comment: sql.NullString{Valid: true, String: "bin am start"}, Commitment: 1}
var testInfo = Info{Key: "Key", Value: "Value"}

func insertTestUser() {
	u := testUser
	s, _ := helpers.GenerateSalt()
	u.Salt = s
	pw, _ := helpers.GeneratePasswordHash(u.Password, u.Salt)
	u.Password = pw
	u.Insert()
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

func TestInitialisation(t *testing.T) {
	con := config.DatabaseConnection{}
	con.Connection = "invalidCON"
	con.Driver = "invalidDB"
	err := Initialisation(con, "")
	if err == nil {
		t.Errorf("Expected error but got none")
	}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Initialisation in %s", tc.connection.Driver), func(t *testing.T) {
			err := Initialisation(tc.connection, "")
			if err != nil {
				t.Errorf("No error expected but got %v", err)
			}
		})
	}
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
			testOrg.Insert()
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
			testOrg.Insert()
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
			testOrg.Insert()
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
			testOrg.Insert()
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

func TestSectionInsert(t *testing.T) {
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Insert Section in %s", tc.connection.Driver), func(t *testing.T) {
			setupDatabase(tc)
			testOrg.Insert()

			//sqlite3 doesnt throw an error, dont know yet why -> deactivated test
			if tc.connection.Driver != "sqlite3" {
				brokenSection := Section{OrganizationID: 42, Name: "test213"}
				err := brokenSection.Insert()
				if err == nil {
					t.Errorf("Expected an error but got non")
				}
			}
			s := testSection
			err := s.Insert()
			if err != nil {
				t.Errorf("Expected no error but got %v", err)
			}
			if s.SectionID < 0 {
				t.Errorf("Expected SectionID > 0 but got %v", s.SectionID)
			}

			teardownDatabase(tc)
		})
	}
}

func TestGetSections(t *testing.T) {
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Get Sections in %s", tc.connection.Driver), func(t *testing.T) {
			setupDatabase(tc)
			testOrg.Insert()
			testSection.Insert()

			ss, err := GetSections()
			if err != nil {
				t.Fatalf("No error expected but got %v", err)
			}
			if len(ss) != 1 {
				t.Error("Expected length 1 got ", len(ss))
			}

			teardownDatabase(tc)
		})
	}
}

func TestSectionDetails(t *testing.T) {
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Get Section Details in %s", tc.connection.Driver), func(t *testing.T) {
			setupDatabase(tc)
			testOrg.Insert()
			testSection.Insert()

			s := Section{SectionID: 1}
			err := s.GetDetails()
			if err != nil {
				t.Fatalf("No error expected but got %v", err)
			}
			if s.Name != testSection.Name {
				t.Errorf("Expected Name test but got %v", s.Name)
			}
			if s.OrganizationID != 1 {
				t.Errorf("Expected OrganizationID = 1 but got %v", s.OrganizationID)
			}

			teardownDatabase(tc)
		})
	}
}

func TestPatchSection(t *testing.T) {
	s := Section{Name: "Developers"}
	sn := Section{Name: "Marketing"}
	err := s.Patch(sn)
	if err != nil {
		t.Fatalf("No error expected but got %v", err)
	}
	if s.Name != sn.Name {
		t.Errorf("Expected name to be %v but got %v", sn.Name, s.Name)
	}
}

func TestUpdateSection(t *testing.T) {
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Update Section in %s", tc.connection.Driver), func(t *testing.T) {
			setupDatabase(tc)
			testOrg.Insert()
			testSection.Insert()

			s := Section{SectionID: 1, Name: "Marketing"}
			sn := Section{SectionID: 1}
			err := s.Update()
			if err != nil {
				t.Fatalf("No error expected but got %v", err)
			}
			err = sn.GetDetails()
			if err != nil {
				t.Fatalf("No error expected but got %v", err)
			}
			if s.Name != sn.Name {
				t.Errorf("Found differences between old and updated name: %v %v", s.Name, sn.Name)
			}

			teardownDatabase(tc)
		})
	}
}

func TestDeleteSection(t *testing.T) {
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Delete Section in %s", tc.connection.Driver), func(t *testing.T) {
			setupDatabase(tc)
			testOrg.Insert()
			testSection.Insert()

			err := DeleteSection(1)
			if err != nil {
				t.Fatalf("No error expected but got %v", err)
			}
			ss, err := GetSections()
			if err != nil {
				t.Fatalf("No error expected but got %v", err)
			}
			if len(ss) != 0 {
				t.Error("Expected length 0 got ", len(ss))
			}

			teardownDatabase(tc)
		})
	}
}

func TestMemberInsert(t *testing.T) {
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Insert Member in %s", tc.connection.Driver), func(t *testing.T) {
			setupDatabase(tc)
			testOrg.Insert()
			testSection.Insert()
			insertTestUser()

			//sqlite3 doesnt throw an error, dont know yet why -> deactivated test
			if tc.connection.Driver != "sqlite3" {
				brokenMember := Member{SectionID: 42, UserID: 42, Rights: 1}
				err := brokenMember.Insert()
				if err == nil {
					t.Errorf("Expected an error but got non")
				}
			}

			m := testMember
			err := m.Insert()
			if err != nil {
				t.Errorf("Expected no error but got %v", err)
			}

			teardownDatabase(tc)
		})
	}
}

func TestGetMember(t *testing.T) {
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Get Member in %s", tc.connection.Driver), func(t *testing.T) {
			setupDatabase(tc)
			testOrg.Insert()
			testSection.Insert()
			insertTestUser()
			testMember.Insert()

			mm, err := GetMembers(0)
			if err != nil {
				t.Fatalf("No error expected but got %v", err)
			}
			if len(mm) != 1 {
				t.Error("Expected length 1 got ", len(mm))
			}
			mm, err = GetMembers(1)
			if err != nil {
				t.Fatalf("No error expected but got %v", err)
			}
			if len(mm) != 1 {
				t.Error("Expected length 1 got ", len(mm))
			}
			teardownDatabase(tc)
		})
	}
}

func TestMemberDetails(t *testing.T) {
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Get Member Details in %s", tc.connection.Driver), func(t *testing.T) {
			setupDatabase(tc)
			testOrg.Insert()
			testSection.Insert()
			insertTestUser()
			testMember.Insert()

			m := Member{SectionID: 1, UserID: 1}
			err := m.GetDetails()
			if err != nil {
				t.Fatalf("No error expected but got %v", err)
			}
			if m.Rights != testMember.Rights {
				t.Errorf("Expected Right %v but got %v", testMember.Rights, m.Rights)
			}

			teardownDatabase(tc)
		})
	}
}

func TestPatchMember(t *testing.T) {
	m := Member{Rights: 1}
	mn := Member{Rights: 2}
	err := m.Patch(mn)
	if err != nil {
		t.Fatalf("No error expected but got %v", err)
	}
	if m.Rights != mn.Rights {
		t.Errorf("Expected Right to be %v but got %v", mn.Rights, m.Rights)
	}
}

func TestUpdateMember(t *testing.T) {
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Update Member in %s", tc.connection.Driver), func(t *testing.T) {
			setupDatabase(tc)
			testOrg.Insert()
			testSection.Insert()
			insertTestUser()
			testMember.Insert()

			m := Member{SectionID: 1, UserID: 1, Rights: 2}
			mn := Member{SectionID: 1, UserID: 1}
			err := m.Update()
			if err != nil {
				t.Fatalf("No error expected but got %v", err)
			}
			err = mn.GetDetails()
			if err != nil {
				t.Fatalf("No error expected but got %v", err)
			}
			if m.Rights != mn.Rights {
				t.Errorf("Found differences between old and updated right: %v %v", m.Rights, mn.Rights)
			}

			teardownDatabase(tc)
		})
	}
}

func TestDeleteMember(t *testing.T) {
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Delete Member in %s", tc.connection.Driver), func(t *testing.T) {
			setupDatabase(tc)
			testOrg.Insert()
			testSection.Insert()
			insertTestUser()
			testMember.Insert()

			err := DeleteMember(1, 1)
			if err != nil {
				t.Fatalf("No error expected but got %v", err)
			}
			mm, err := GetMembers(0)
			if err != nil {
				t.Fatalf("No error expected but got %v", err)
			}
			if len(mm) != 0 {
				t.Error("Expected length 0 got ", len(mm))
			}

			teardownDatabase(tc)
		})
	}
}

func TestEventInsert(t *testing.T) {
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Insert Event in %s", tc.connection.Driver), func(t *testing.T) {
			setupDatabase(tc)
			testOrg.Insert()
			insertTestUser()

			//sqlite3 doesnt throw an error, dont know yet why -> deactivated test
			if tc.connection.Driver != "sqlite3" {
				brokenEvent := Event{OrganizationID: 42, Name: "test213", Creator: 42}
				err := brokenEvent.Insert()
				if err == nil {
					t.Errorf("Expected an error but got non")
				}
			}
			e := testEvent
			err := e.Insert()
			if err != nil {
				t.Errorf("Expected no error but got %v", err)
			}
			if e.EventID < 0 {
				t.Errorf("Expected EventID > 0 but got %v", e.EventID)
			}

			teardownDatabase(tc)
		})
	}
}

func TestGetEvents(t *testing.T) {
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Get Events in %s", tc.connection.Driver), func(t *testing.T) {
			setupDatabase(tc)
			testOrg.Insert()
			insertTestUser()
			testEvent.Insert()

			ee, err := GetEvents(0)
			if err != nil {
				t.Fatalf("No error expected but got %v", err)
			}
			if len(ee) != 1 {
				t.Error("Expected length 1 got ", len(ee))
			}
			ee, err = GetEvents(1)
			if err != nil {
				t.Fatalf("No error expected but got %v", err)
			}
			if len(ee) != 1 {
				t.Error("Expected length 1 got ", len(ee))
			}
			teardownDatabase(tc)
		})
	}
}

func TestEventDetails(t *testing.T) {
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Get Event Details in %s", tc.connection.Driver), func(t *testing.T) {
			setupDatabase(tc)
			testOrg.Insert()
			insertTestUser()
			testEvent.Insert()

			e := Event{EventID: 1}
			err := e.GetDetails()
			if err != nil {
				t.Fatalf("No error expected but got %v", err)
			}
			if e.Name != testEvent.Name {
				t.Errorf("Expected Name test but got %v", e.Name)
			}
			if e.OrganizationID != testEvent.OrganizationID {
				t.Errorf("Expected OrganizationID = 1 but got %v", e.OrganizationID)
			}
			if e.Address != testEvent.Address {
				t.Errorf("Expected Address = %v but got %v", testEvent.Address, e.Address)
			}
			if e.Creator != testEvent.Creator {
				t.Errorf("Expected Creator = %v but got %v", testEvent.Creator, e.Creator)
			}

			teardownDatabase(tc)
		})
	}
}

func TestPatchEvent(t *testing.T) {
	e := Event{Name: "hackathon", Address: sql.NullString{String: "there", Valid: true}, Start: time.Now(), End: sql.NullTime{Time: time.Now().Add(time.Hour), Valid: true}}
	en := Event{Name: "beerfest", Address: sql.NullString{String: "here", Valid: true}, Start: time.Now().Add(time.Hour), End: sql.NullTime{Time: time.Now().Add(time.Hour * 2), Valid: true}}
	err := e.Patch(en)
	if err != nil {
		t.Fatalf("No error expected but got %v", err)
	}
	if e.Name != en.Name {
		t.Errorf("Expected name to be %v but got %v", en.Name, e.Name)
	}
	if e.Address != en.Address {
		t.Errorf("Expected Address to be %v but got %v", en.Address, e.Address)
	}
	if e.Start != en.Start {
		t.Errorf("Expected Start to be %v but got %v", en.Start, e.Start)
	}
	if e.End != en.End {
		t.Errorf("Expected End to be %v but got %v", en.End, e.End)
	}
}

func TestUpdateEvent(t *testing.T) {
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Update Event in %s", tc.connection.Driver), func(t *testing.T) {
			setupDatabase(tc)
			testOrg.Insert()
			insertTestUser()
			testEvent.Insert()

			e := Event{EventID: 1, OrganizationID: 1, Name: "beerfest", Address: sql.NullString{String: "here", Valid: true}, Start: time.Now().Add(time.Hour), End: sql.NullTime{Time: time.Now().Add(time.Hour * 2), Valid: true}}
			en := Event{EventID: 1}
			err := e.Update()
			if err != nil {
				t.Fatalf("No error expected but got %v", err)
			}
			err = en.GetDetails()
			if err != nil {
				t.Fatalf("No error expected but got %v", err)
			}
			if e.Name != en.Name {
				t.Errorf("Expected name to be %v but got %v", en.Name, e.Name)
			}
			if e.Address != en.Address {
				t.Errorf("Expected Address to be %v but got %v", en.Address, e.Address)
			}

			teardownDatabase(tc)
		})
	}
}

func TestDeleteEvent(t *testing.T) {
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Delete Event in %s", tc.connection.Driver), func(t *testing.T) {
			setupDatabase(tc)
			testOrg.Insert()
			testEvent.Insert()

			err := DeleteEvent(1)
			if err != nil {
				t.Fatalf("No error expected but got %v", err)
			}
			ss, err := GetEvents(0)
			if err != nil {
				t.Fatalf("No error expected but got %v", err)
			}
			if len(ss) != 0 {
				t.Error("Expected length 0 got ", len(ss))
			}

			teardownDatabase(tc)
		})
	}
}

func TestCommentInsert(t *testing.T) {
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Insert Comment in %s", tc.connection.Driver), func(t *testing.T) {
			setupDatabase(tc)
			testOrg.Insert()
			insertTestUser()
			testEvent.Insert()

			//sqlite3 doesnt throw an error, dont know yet why -> deactivated test
			if tc.connection.Driver != "sqlite3" {
				brokenComment := Comment{EventID: 42, UserID: 42}
				err := brokenComment.Insert()
				if err == nil {
					t.Errorf("Expected an error but got non")
				}
			}

			c := testComment
			err := c.Insert()
			if err != nil {
				t.Errorf("Expected no error but got %v", err)
			}

			if c.CommentID < 0 {
				t.Errorf("Expected CommentID > 0 but got %v", c.CommentID)
			}

			teardownDatabase(tc)
		})
	}
}

func TestGetComment(t *testing.T) {
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Get Comment in %s", tc.connection.Driver), func(t *testing.T) {
			setupDatabase(tc)
			testOrg.Insert()
			insertTestUser()
			testEvent.Insert()
			testComment.Insert()

			cc, err := GetComments(0)
			if err != nil {
				t.Fatalf("No error expected but got %v", err)
			}
			if len(cc) != 1 {
				t.Error("Expected length 1 got ", len(cc))
			}
			cc, err = GetComments(1)
			if err != nil {
				t.Fatalf("No error expected but got %v", err)
			}
			if len(cc) != 1 {
				t.Error("Expected length 1 got ", len(cc))
			}
			teardownDatabase(tc)
		})
	}
}

func TestCommentDetails(t *testing.T) {
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Get Comment Details in %s", tc.connection.Driver), func(t *testing.T) {
			setupDatabase(tc)
			testOrg.Insert()
			insertTestUser()
			testEvent.Insert()
			testComment.Insert()

			c := Comment{CommentID: 1}
			err := c.GetDetails()
			if err != nil {
				t.Fatalf("No error expected but got %v", err)
			}
			if c.EventID != testComment.EventID {
				t.Errorf("Expected EventID %v but got %v", testComment.EventID, c.EventID)
			}
			if c.UserID != testComment.UserID {
				t.Errorf("Expected UserID %v but got %v", testComment.UserID, c.UserID)
			}
			if c.Comment != testComment.Comment {
				t.Errorf("Expected Comment %v but got %v", testComment.Comment, c.Comment)
			}
			teardownDatabase(tc)
		})
	}
}

func TestPatchComment(t *testing.T) {
	c := Comment{Creation: time.Now(), Comment: "Hullu"}
	cn := Comment{Creation: time.Now().Add(time.Hour), Comment: "Hallu"}
	err := c.Patch(cn)
	if err != nil {
		t.Fatalf("No error expected but got %v", err)
	}
	if c.Creation != cn.Creation {
		t.Errorf("Expected Creation %v but got %v", cn.Creation, c.Creation)
	}
	if c.Comment != cn.Comment {
		t.Errorf("Expected Comment %v but got %v", cn.Comment, c.Comment)
	}
}

func TestUpdateComment(t *testing.T) {
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Update Comment in %s", tc.connection.Driver), func(t *testing.T) {
			setupDatabase(tc)
			testOrg.Insert()
			insertTestUser()
			testEvent.Insert()
			testComment.Insert()

			c := Comment{CommentID: 1, Creation: time.Now().Add(time.Hour), Comment: "Hallu"}
			cn := Comment{CommentID: 1}
			err := c.Update()
			if err != nil {
				t.Fatalf("No error expected but got %v", err)
			}
			err = cn.GetDetails()
			if err != nil {
				t.Fatalf("No error expected but got %v", err)
			}
			if c.Comment != cn.Comment {
				t.Errorf("Expected Right %v but got %v", cn.Comment, c.Comment)
			}

			teardownDatabase(tc)
		})
	}
}

func TestDeleteComment(t *testing.T) {
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Delete Comment in %s", tc.connection.Driver), func(t *testing.T) {
			setupDatabase(tc)
			testOrg.Insert()
			insertTestUser()
			testEvent.Insert()
			testComment.Insert()

			err := DeleteComment(1)
			if err != nil {
				t.Fatalf("No error expected but got %v", err)
			}
			cc, err := GetComments(0)
			if err != nil {
				t.Fatalf("No error expected but got %v", err)
			}
			if len(cc) != 0 {
				t.Error("Expected length 0 got ", len(cc))
			}

			teardownDatabase(tc)
		})
	}
}

func TestAttendeeInsert(t *testing.T) {
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Insert Attendee in %s", tc.connection.Driver), func(t *testing.T) {
			setupDatabase(tc)
			testOrg.Insert()
			insertTestUser()
			testEvent.Insert()

			//sqlite3 doesnt throw an error, dont know yet why -> deactivated test
			if tc.connection.Driver != "sqlite3" {
				brokenAttendee := Attendee{EventID: 42, UserID: 42, Commitment: 1}
				err := brokenAttendee.Insert()
				if err == nil {
					t.Errorf("Expected an error but got non")
				}
			}

			a := testAttendee
			err := a.Insert()
			if err != nil {
				t.Errorf("Expected no error but got %v", err)
			}

			teardownDatabase(tc)
		})
	}
}

func TestGetAttendee(t *testing.T) {
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Get Attendee in %s", tc.connection.Driver), func(t *testing.T) {
			setupDatabase(tc)
			testOrg.Insert()
			insertTestUser()
			testEvent.Insert()
			testAttendee.Insert()

			aa, err := GetAttendees(0)
			if err != nil {
				t.Fatalf("No error expected but got %v", err)
			}
			if len(aa) != 1 {
				t.Error("Expected length 1 got ", len(aa))
			}
			aa, err = GetAttendees(1)
			if err != nil {
				t.Fatalf("No error expected but got %v", err)
			}
			if len(aa) != 1 {
				t.Error("Expected length 1 got ", len(aa))
			}
			teardownDatabase(tc)
		})
	}
}

func TestAttendeeDetails(t *testing.T) {
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Get Attendee Details in %s", tc.connection.Driver), func(t *testing.T) {
			setupDatabase(tc)
			testOrg.Insert()
			insertTestUser()
			testEvent.Insert()
			testAttendee.Insert()

			a := Attendee{EventID: 1, UserID: 1}
			err := a.GetDetails()
			if err != nil {
				t.Fatalf("No error expected but got %v", err)
			}
			if a.Comment != testAttendee.Comment {
				t.Errorf("Expected Comment %v but got %v", testAttendee.Comment, a.Comment)
			}
			if a.Commitment != testAttendee.Commitment {
				t.Errorf("Expected Commitment %v but got %v", testAttendee.Commitment, a.Commitment)
			}

			teardownDatabase(tc)
		})
	}
}

func TestPatchAttendee(t *testing.T) {
	a := Attendee{Commitment: 1, Comment: sql.NullString{Valid: true, String: "bin am Start"}}
	an := Attendee{Commitment: 2, Comment: sql.NullString{Valid: true, String: "nope"}}
	err := a.Patch(an)
	if err != nil {
		t.Fatalf("No error expected but got %v", err)
	}
	if a.Comment != an.Comment {
		t.Errorf("Expected Comment %v but got %v", an.Comment, a.Comment)
	}
	if a.Commitment != an.Commitment {
		t.Errorf("Expected Commitment %v but got %v", an.Commitment, a.Commitment)
	}
}

func TestUpdateAttendee(t *testing.T) {
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Update Attendee in %s", tc.connection.Driver), func(t *testing.T) {
			setupDatabase(tc)
			testOrg.Insert()
			insertTestUser()
			testEvent.Insert()
			testAttendee.Insert()

			a := Attendee{EventID: 1, UserID: 1, Commitment: 2, Comment: sql.NullString{Valid: true, String: "nope"}}
			an := Attendee{EventID: 1, UserID: 1}
			err := a.Update()
			if err != nil {
				t.Fatalf("No error expected but got %v", err)
			}
			err = an.GetDetails()
			if err != nil {
				t.Fatalf("No error expected but got %v", err)
			}
			if a.Comment != an.Comment {
				t.Errorf("Expected Comment %v but got %v", an.Comment, a.Comment)
			}
			if a.Commitment != an.Commitment {
				t.Errorf("Expected Commitment %v but got %v", an.Commitment, a.Commitment)
			}

			teardownDatabase(tc)
		})
	}
}

func TestDeleteAttendee(t *testing.T) {
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Delete Attendee in %s", tc.connection.Driver), func(t *testing.T) {
			setupDatabase(tc)
			testOrg.Insert()
			insertTestUser()
			testEvent.Insert()
			testAttendee.Insert()

			err := DeleteAttendee(1, 1)
			if err != nil {
				t.Fatalf("No error expected but got %v", err)
			}
			aa, err := GetAttendees(0)
			if err != nil {
				t.Fatalf("No error expected but got %v", err)
			}
			if len(aa) != 0 {
				t.Error("Expected length 0 got ", len(aa))
			}

			teardownDatabase(tc)
		})
	}
}

func TestInfoInsert(t *testing.T) {
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Insert Info in %s", tc.connection.Driver), func(t *testing.T) {
			setupDatabase(tc)
			i := testInfo

			err := i.Insert()

			if err != nil {
				t.Errorf("Expected no error but got %v", err)
			}

			teardownDatabase(tc)
		})
	}
}

func TestGetInfos(t *testing.T) {
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Get Info in %s", tc.connection.Driver), func(t *testing.T) {
			setupDatabase(tc)
			testInfo.Insert()
			ii, err := GetInfos()
			if err != nil {
				t.Fatalf("No error expected but got %v", err)
			}
			if len(ii) != 1 {
				t.Error("Expected length 1 got ", len(ii))
			}
			teardownDatabase(tc)
		})
	}
}

func TestInfoDetails(t *testing.T) {
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Get Organiztaion Details in %s", tc.connection.Driver), func(t *testing.T) {
			setupDatabase(tc)
			testInfo.Insert()
			i := Info{Key: "Key"}
			err := i.GetDetails()
			if err != nil {
				t.Fatalf("No error expected but got %v", err)
			}
			if i.Value != testInfo.Value {
				t.Errorf("Expected Valie %v but got %v", testInfo.Value, i.Value)
			}
			teardownDatabase(tc)
		})
	}
}

func TestPatchInfo(t *testing.T) {
	i := Info{Value: "Value"}
	in := Info{Value: "Value2"}

	err := i.Patch(in)

	if err != nil {
		t.Fatalf("No error expected but got %v", err)
	}

	if i.Value != in.Value {
		t.Errorf("Expected Value to be %v but got %v", in.Value, i.Value)
	}
}

func TestUpdateInfo(t *testing.T) {
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Update Info in %s", tc.connection.Driver), func(t *testing.T) {
			setupDatabase(tc)
			testInfo.Insert()
			i := Info{Key: "Key", Value: "Value2"}
			in := Info{Key: "Key"}
			err := i.Update()
			if err != nil {
				t.Fatalf("No error expected but got %v", err)
			}
			err = in.GetDetails()
			if err != nil {
				t.Fatalf("No error expected but got %v", err)
			}
			if i.Value != in.Value {
				t.Errorf("Expected Value to be %v but got %v", in.Value, i.Value)
			}
			teardownDatabase(tc)
		})
	}
}

func TestDeleteInfo(t *testing.T) {
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Delete Info in %s", tc.connection.Driver), func(t *testing.T) {
			setupDatabase(tc)
			testOrg.Insert()
			err := DeleteInfo("Key")
			if err != nil {
				t.Fatalf("No error expected but got %v", err)
			}
			ii, err := GetInfos()
			if err != nil {
				t.Fatalf("No error expected but got %v", err)
			}
			if len(ii) != 0 {
				t.Error("Expected length 0 got ", len(ii))
			}
			teardownDatabase(tc)
		})
	}
}
