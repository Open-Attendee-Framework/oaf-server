package db100

import (
	"database/sql"
	"errors"

	"github.com/concertLabs/oaf-server/internal/helpers"
)

//User is a singe User of the OAF
type User struct {
	UserID    int            `json:"userid" db:"UserID"`
	Username  string         `json:"username" db:"Username"`
	Password  string         `json:"password" db:"Password"`
	Salt      string         `json:"-" db:"Salt"`
	Shownname sql.NullString `json:"shownname" db:"Shownname"`
	EMail     string         `json:"email" db:"EMail"`
	SuperUser bool           `json:"superuser" db:"SuperUser"`
}

//DoesUserExist gives back a boolean if the User with this Username can be found in the Database
func DoesUserExist(username string) (bool, error) {
	var u User
	query := db.Rebind(`SELECT * FROM "Users" WHERE "Username" = ? LIMIT 1`)
	err := db.Get(&u, query, username)
	b := (u.UserID > 0)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return b, errors.New("Error finding User:" + err.Error())
	}
	return b, nil
}

//GetUsers gives back all Users in the Database
func GetUsers() ([]User, error) {
	var u []User
	err := db.Select(&u, `SELECT * FROM "Users"`)
	if err != nil {
		return u, errors.New("Error getting Users:" + err.Error())
	}
	return u, nil
}

//GetDetailstoUsername takes a User struct with only the Username and tries to fetch the remaining infos
func (u *User) GetDetailstoUsername() error {
	query := db.Rebind(`SELECT * FROM "Users" WHERE "Username" = ? LIMIT 1`)
	err := db.Get(u, query, u.Username)
	if err != nil {
		return errors.New("Error getting user details:" + err.Error())
	}
	return nil
}

//GetDetails takes a User struct with only the UserID and tries to fetch the remaining infos
func (u *User) GetDetails() error {
	query := db.Rebind(`SELECT * FROM "Users" WHERE "UserID" = ? LIMIT 1`)
	err := db.Get(u, query, u.UserID)
	if err != nil {
		return errors.New("Error getting user details:" + err.Error())
	}
	return nil
}

//Patch patches a User with new Info from a second struct
func (u *User) Patch(ou User) error {
	u.Username = helpers.CopyIfNotEmpty(u.Username, ou.Username)
	if ou.Password != "" {
		p, err := helpers.GeneratePasswordHash(ou.Password, u.Salt)
		if err != nil {
			return errors.New("Error hashing password: " + err.Error())
		}
		u.Password = p
	}
	u.EMail = helpers.CopyIfNotEmpty(u.EMail, ou.EMail)
	u.Shownname.String = helpers.CopyIfNotEmpty(u.Shownname.String, ou.Shownname.String)
	return nil
}

//Update updates all Userfields in the Database
func (u *User) Update() error {
	err := updateDBO(u)
	if err != nil {
		return errors.New("Error updating user:" + err.Error())
	}
	return nil
}

func (u *User) getID() interface{} {
	return u.UserID
}

func (u *User) getTablename() string {
	return "Users"
}

func (u *User) getIDColumn() string {
	return "UserID"
}

func (u *User) getInsertColumns() []string {
	return []string{"Username", "Password", "Salt", "EMail", "SuperUser"}
}

func (u *User) getInsertFields() []interface{} {
	var interfaceSlice []interface{}
	interfaceSlice = append(interfaceSlice, u.Username)
	interfaceSlice = append(interfaceSlice, u.Password)
	interfaceSlice = append(interfaceSlice, u.Salt)
	interfaceSlice = append(interfaceSlice, u.EMail)
	interfaceSlice = append(interfaceSlice, u.SuperUser)
	return interfaceSlice
}

func (u *User) getUpdateColumns() []string {
	return u.getInsertColumns()
}

func (u *User) getUpdateFields() []interface{} {
	return u.getInsertFields()
}

//Insert inserts a new User into the database and adding the new UserID into the struct
func (u *User) Insert() error {
	var err error
	u.UserID, err = insertDBO(u)
	if err != nil {
		return errors.New("Error inserting User:" + err.Error())
	}
	return nil
}

//DeleteUser delets user with the given UserID
func DeleteUser(id int) error {
	query := db.Rebind(`DELETE FROM "Users" WHERE "UserID" = ?`)
	_, err := db.Exec(query, id)
	if err != nil {
		return errors.New("Error deleting User: " + err.Error())
	}
	return nil
}
