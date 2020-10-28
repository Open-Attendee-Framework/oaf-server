package db100

import (
	"errors"
	"fmt"
)

type databaseObject interface {
	getTablename() string
	getIDColumn() string
	getInsertColumns() []string
	getInsertFields() []interface{}
}

func insertDBO(dbo databaseObject) (int, error) {
	var err error
	var nid int
	query := `INSERT INTO "` + dbo.getTablename() + `" `
	columns := dbo.getInsertColumns()
	for i, c := range columns {
		query = query + queryLetter(i) + `"` + c + `"`
	}
	query = query + ") VALUES "
	for i := 0; i < len(columns); i++ {
		query = query + queryLetter(i) + "?"
	}
	query = query + ")"
	query = db.Rebind(query)
	if db.DriverName() == pgDriverName {
		nid, err = insertDBOPG(query, dbo.getIDColumn(), dbo.getInsertFields())
	} else {
		nid, err = insertDBOOther(query, dbo.getInsertFields())
	}
	return nid, err
}

func insertDBOPG(query string, returning string, a []interface{}) (int, error) {
	var newid int
	query = query + ` RETURNING "` + returning + `"`
	tx := db.MustBegin()
	fmt.Println(query)
	stmt, err := tx.Prepare(query)
	if err != nil {
		tx.Rollback()
		return -1, errors.New("Error preparing Statement:" + err.Error())
	}
	stmt.QueryRow(a...).Scan(newid)
	if err != nil {
		tx.Rollback()
		return -1, errors.New("Error executing Statement:" + err.Error())
	}
	err = tx.Commit()
	if err != nil {
		return -1, errors.New("Error executing Commit:" + err.Error())
	}
	return newid, nil
}

func insertDBOOther(query string, a []interface{}) (int, error) {
	res, err := db.Exec(query, a...)
	if err != nil {
		return -1, errors.New("Error inserting: " + err.Error())
	}
	id, err := res.LastInsertId()
	if err != nil {
		return -1, errors.New("Error fetching new ID: " + err.Error())
	}
	return int(id), nil
}
