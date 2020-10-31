package db100

import (
	"errors"
)

type databaseObject interface {
	getID() int
	getTablename() string
	getIDColumn() string
	getInsertColumns() []string
	getInsertFields() []interface{}
	getUpdateColumns() []string
	getUpdateFields() []interface{}
}

func queryLetter(i int) string {
	var result string
	if i == 0 {
		result = "( "
	} else {
		result = ", "
	}
	return result
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
	if returning != "" {
		query = query + ` RETURNING "` + returning + `"`
	}
	tx := db.MustBegin()
	stmt, err := tx.Prepare(query)
	if err != nil {
		tx.Rollback()
		return -1, errors.New("Error preparing Statement:" + err.Error())
	}
	if returning != "" {
		stmt.QueryRow(a...).Scan(newid)
		if err != nil {
			tx.Rollback()
			return -1, errors.New("Error executing Statement:" + err.Error())
		}
	} else {
		stmt.QueryRow(a...)
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

func updateDBO(dbo databaseObject) error {
	query := `UPDATE "` + dbo.getTablename() + `" SET`
	columns := dbo.getUpdateColumns()
	for i, c := range columns {
		query = query + ` "` + c + `" = ?`
		if i != (len(columns) - 1) {
			query = query + ","
		}
	}
	query = query + ` WHERE "` + dbo.getIDColumn() + `" = ?`
	query = db.Rebind(query)
	fields := dbo.getUpdateFields()
	fields = append(fields, dbo.getID())
	_, err := db.Exec(query, fields...)
	if err != nil {
		return errors.New("Error updating: " + err.Error())
	}
	return nil
}
