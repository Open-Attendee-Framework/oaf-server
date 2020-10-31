package db100

import (
	"errors"
)

//Info hold Data from the Info table which is a general key value store
type Info struct {
	Key   string `json:"key" db:"Key"`
	Value string `json:"value" db:"Value"`
}

func (i *Info) getID() interface{} {
	return i.Key
}

func (i *Info) getTablename() string {
	return "Info"
}

func (i *Info) getIDColumn() string {
	return "Key"
}

func (i *Info) getInsertColumns() []string {
	result := i.getUpdateColumns()
	result = append(result, "Key")
	return result
}

func (i *Info) getInsertFields() []interface{} {
	var interfaceSlice = i.getUpdateFields()
	interfaceSlice = append(interfaceSlice, i.Key)
	return interfaceSlice
}

func (i *Info) getUpdateColumns() []string {
	return []string{"Value"}
}

func (i *Info) getUpdateFields() []interface{} {
	var interfaceSlice []interface{}
	interfaceSlice = append(interfaceSlice, i.Value)
	return interfaceSlice
}

//Insert inserts a new Organization into the database and adding the new OrganizationID into the struct
func (i *Info) Insert() error {
	var err error
	_, err = insertDBO(i)
	if err != nil {
		return errors.New("Error inserting Info:" + err.Error())
	}
	return nil
}

//GetInfos gives back all Infos in the Database
func GetInfos() ([]Info, error) {
	var i []Info
	var nfo Info
	query, in := buildSelectQuery(&nfo, "", -1)
	err := db.Select(&i, query, in...)
	if err != nil {
		return i, errors.New("Error getting Info:" + err.Error())
	}
	return i, nil
}

//GetDetails takes a Info struct with only the Key and tries to fetch the remaining infos
func (i *Info) GetDetails() error {
	query := db.Rebind(`SELECT * FROM "Info" WHERE "Key" = ? LIMIT 1`)
	err := db.Get(i, query, i.Key)
	if err != nil {
		return errors.New("Error getting Info details:" + err.Error())
	}
	return nil
}

//Patch patches a Organization with new Info from a second struct
func (i *Info) Patch(ii Info) error {
	i.Value = ii.Value
	return nil
}

//Update updates all Organization Fields in the Database
func (i *Info) Update() error {
	err := updateDBO(i)
	if err != nil {
		return errors.New("Error updating Info:" + err.Error())
	}
	return nil
}

//DeleteInfo deletes a Info with the given Key
func DeleteInfo(key string) error {
	query := db.Rebind(`DELETE FROM "Info" WHERE "Key" = ?`)
	_, err := db.Exec(query, key)
	if err != nil {
		return errors.New("Error deleting Key: " + err.Error())
	}
	return nil
}
