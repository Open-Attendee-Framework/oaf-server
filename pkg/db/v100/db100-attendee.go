package db100

import (
	"database/sql"
	"errors"

	"github.com/concertLabs/oaf-server/internal/helpers"
)

//Attendee manages Useres that attend a single Event
type Attendee struct {
	EventID    int            `json:"eventid" db:"EventID"`
	UserID     int            `json:"userid" db:"UserID"`
	Commitment int            `json:"commitment" db:"Commitment"`
	Comment    sql.NullString `json:"comment" db:"Comment"`
}

func (a *Attendee) getID() interface{} {
	return a.UserID
}

func (a *Attendee) getTablename() string {
	return "Attendees"
}

func (a *Attendee) getIDColumn() string {
	return ""
}

func (a *Attendee) getInsertColumns() []string {
	return []string{"EventID", "UserID", "Commitment", "Comment"}
}

func (a *Attendee) getInsertFields() []interface{} {
	var interfaceSlice []interface{}
	interfaceSlice = append(interfaceSlice, a.EventID)
	interfaceSlice = append(interfaceSlice, a.UserID)
	interfaceSlice = append(interfaceSlice, a.Commitment)
	interfaceSlice = append(interfaceSlice, a.Comment)
	return interfaceSlice
}

func (a *Attendee) getUpdateColumns() []string {
	return []string{}
}

func (a *Attendee) getUpdateFields() []interface{} {
	var interfaceSlice []interface{}
	return interfaceSlice
}

//Insert inserts a new Attendee into the database
func (a *Attendee) Insert() error {
	var err error
	_, err = insertDBO(a)
	if err != nil {
		return errors.New("Error inserting Attendee:" + err.Error())
	}
	return nil
}

//GetAttendees gives back all Attendees in the Database
func GetAttendees(eventid int) ([]Attendee, error) {
	var a []Attendee
	var at Attendee
	query, in := buildSelectQuery(&at, "EventID", eventid)
	err := db.Select(&a, query, in...)
	if err != nil {
		return a, errors.New("Error getting Attendees:" + err.Error())
	}
	return a, nil
}

//GetDetails takes a Attendee struct with only the UserID and EventID and tries to fetch the remaining infos
func (a *Attendee) GetDetails() error {
	query := db.Rebind(`SELECT * FROM "Attendees" WHERE "UserID" = ? AND "EventID" = ? LIMIT 1`)
	err := db.Get(a, query, a.UserID, a.EventID)
	if err != nil {
		return errors.New("Error getting Attendee details:" + err.Error())
	}
	return nil
}

//Patch patches a Member with new Info from a second struct
func (a *Attendee) Patch(aa Attendee) error {
	a.Commitment = aa.Commitment
	if aa.Comment.Valid {
		a.Comment.String = helpers.CopyIfNotEmpty(a.Comment.String, aa.Comment.String)
		a.Comment.Valid = true
	}
	return nil
}

//Update updates the Commitment and Comment Field of a Attendee in the Database
func (a *Attendee) Update() error {
	query := `UPDATE "Attendees" SET "Commitment" = ?, "Comment" = ? WHERE "UserID" = ? and "EventID" = ?`
	query = db.Rebind(query)
	_, err := db.Exec(query, a.Commitment, a.Comment, a.UserID, a.EventID)
	if err != nil {
		return errors.New("Error updating Attendee: " + err.Error())
	}
	return nil
}

//DeleteAttendee deletes a Attendee with the given UserID and EventID
func DeleteAttendee(UserID int, EventID int) error {
	query := db.Rebind(`DELETE FROM "Attendees" WHERE "UserID" = ? and "EventID" = ?`)
	_, err := db.Exec(query, UserID, EventID)
	if err != nil {
		return errors.New("Error deleting Attendee: " + err.Error())
	}
	return nil
}
