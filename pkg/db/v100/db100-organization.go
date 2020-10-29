package db100

import (
	"errors"

	"github.com/concertLabs/oaf-server/internal/helpers"
)

//Organization holds infos about a single Org (like orchestra, football club you name it)
type Organization struct {
	OrganizationID int    `json:"organizationid" db:"OrganizationID"`
	Name           string `json:"name" db:"Name"`
	Picture        []byte `json:"picture" db:"Picture"`
}

func (o *Organization) getID() int {
	return o.OrganizationID
}

func (o *Organization) getTablename() string {
	return "Organizations"
}

func (o *Organization) getIDColumn() string {
	return "OrganizationID"
}

func (o *Organization) getInsertColumns() []string {
	return []string{"Name", "Picture"}
}

func (o *Organization) getInsertFields() []interface{} {
	var interfaceSlice []interface{}
	interfaceSlice = append(interfaceSlice, o.Name)
	interfaceSlice = append(interfaceSlice, o.Picture)
	return interfaceSlice
}

func (o *Organization) getUpdateColumns() []string {
	return o.getInsertColumns()
}

func (o *Organization) getUpdateFields() []interface{} {
	return o.getInsertFields()
}

//Insert inserts a new Organization into the database and adding the new OrganizationID into the struct
func (o *Organization) Insert() error {
	var err error
	o.OrganizationID, err = insertDBO(o)
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
	err := updateDBO(o)
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
