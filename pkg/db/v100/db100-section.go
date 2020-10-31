package db100

import (
	"errors"

	"github.com/concertLabs/oaf-server/internal/helpers"
)

//Section is a part of the Club (like Baritons or Handball Section)
type Section struct {
	SectionID      int    `json:"sectionid" db:"SectionID"`
	OrganizationID int    `json:"organizationid" db:"OrganizationID"`
	Name           string `json:"name" db:"Name"`
}

func (s *Section) getID() interface{} {
	return s.SectionID
}

func (s *Section) getTablename() string {
	return "Sections"
}

func (s *Section) getIDColumn() string {
	return "SectionID"
}

func (s *Section) getInsertColumns() []string {
	result := s.getUpdateColumns()
	result = append(result, "OrganizationID")
	return result
}

func (s *Section) getInsertFields() []interface{} {
	var interfaceSlice []interface{}
	interfaceSlice = append(interfaceSlice, s.getUpdateFields()...)
	interfaceSlice = append(interfaceSlice, s.OrganizationID)
	return interfaceSlice
}

func (s *Section) getUpdateColumns() []string {
	return []string{"Name"}
}

func (s *Section) getUpdateFields() []interface{} {
	var interfaceSlice []interface{}
	interfaceSlice = append(interfaceSlice, s.Name)
	return interfaceSlice
}

//Insert inserts a new Section into the database and adding the new SectionID into the struct
func (s *Section) Insert() error {
	var err error
	s.SectionID, err = insertDBO(s)
	if err != nil {
		return errors.New("Error inserting Section:" + err.Error())
	}
	return nil
}

//GetSections gives back all Sections in the Database
func GetSections() ([]Section, error) {
	var s []Section
	var se Section
	query, in := buildSelectQuery(&se, "", -1)
	err := db.Select(&s, query, in...)
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
	err := updateDBO(s)
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
