package db100

import (
	"database/sql"
	"errors"
	"time"

	"github.com/concertLabs/oaf-server/internal/helpers"
)

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

func (e *Event) getID() interface{} {
	return e.EventID
}

func (e *Event) getTablename() string {
	return "Events"
}

func (e *Event) getIDColumn() string {
	return "EventID"
}

func (e *Event) getInsertColumns() []string {
	result := e.getUpdateColumns()
	result = append(result, "Creator")
	return result
}

func (e *Event) getInsertFields() []interface{} {
	var interfaceSlice []interface{}
	interfaceSlice = append(interfaceSlice, e.getUpdateFields()...)
	interfaceSlice = append(interfaceSlice, e.Creator)
	return interfaceSlice
}

func (e *Event) getUpdateColumns() []string {
	return []string{"OrganizationID", "Name", "Address", "Start", "End"}
}

func (e *Event) getUpdateFields() []interface{} {
	var interfaceSlice []interface{}
	interfaceSlice = append(interfaceSlice, e.OrganizationID)
	interfaceSlice = append(interfaceSlice, e.Name)
	interfaceSlice = append(interfaceSlice, e.Address)
	interfaceSlice = append(interfaceSlice, e.Start)
	interfaceSlice = append(interfaceSlice, e.End)
	return interfaceSlice
}

//Insert inserts a new Event into the database and adding the new EventID into the struct
func (e *Event) Insert() error {
	var err error
	e.EventID, err = insertDBO(e)
	if err != nil {
		return errors.New("Error inserting Event:" + err.Error())
	}
	return nil
}

//GetEvents gives back all Events in the Database
func GetEvents(orgid int) ([]Event, error) {
	var e []Event
	query := `SELECT * FROM "Events"`
	if orgid > 0 {
		query = query + ` WHERE "OrganizationID" = ?`
	}
	query = db.Rebind(query)
	var err error
	if orgid > 0 {
		err = db.Select(&e, query, orgid)
	} else {
		err = db.Select(&e, query)
	}
	if err != nil {
		return e, errors.New("Error getting Events:" + err.Error())
	}
	return e, nil
}

//GetDetails takes a Event struct with only the EventID and tries to fetch the remaining infos
func (e *Event) GetDetails() error {
	query := db.Rebind(`SELECT * FROM "Events" WHERE "EventID" = ? LIMIT 1`)
	err := db.Get(e, query, e.EventID)
	if err != nil {
		return errors.New("Error getting Event details:" + err.Error())
	}
	return nil
}

//Patch patches a Event with new Info from a second struct
func (e *Event) Patch(ee Event) error {
	e.Name = helpers.CopyIfNotEmpty(e.Name, ee.Name)
	e.Address.String = helpers.CopyIfNotEmpty(e.Address.String, ee.Address.String)
	e.Start = ee.Start
	if ee.End.Valid {
		e.End = ee.End
	}
	return nil
}

//Update updates all Event Fields in the Database
func (e *Event) Update() error {
	err := updateDBO(e)
	if err != nil {
		return errors.New("Error updating Event:" + err.Error())
	}
	return nil
}

//DeleteEvent deletes a Event with the given EventID
func DeleteEvent(id int) error {
	query := db.Rebind(`DELETE FROM "Events" WHERE "EventID" = ?`)
	_, err := db.Exec(query, id)
	if err != nil {
		return errors.New("Error deleting Event: " + err.Error())
	}
	return nil
}
