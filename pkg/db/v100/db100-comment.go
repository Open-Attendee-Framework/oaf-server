package db100

import (
	"errors"
	"time"

	"github.com/concertLabs/oaf-server/internal/helpers"
)

//Comment holds User message for a Event "Blackboard"
type Comment struct {
	CommentID int       `json:"commentid" db:"CommentID"`
	EventID   int       `json:"eventid" db:"EventID"`
	UserID    int       `json:"userid" db:"UserID"`
	Creation  time.Time `json:"creation" db:"Creation"`
	Comment   string    `json:"comment"  db:"Comment"`
}

func (c *Comment) getIDs() []interface{} {
	var interfaceSlice []interface{}
	interfaceSlice = append(interfaceSlice, c.CommentID)
	return interfaceSlice
}

func (c *Comment) getTablename() string {
	return "Comments"
}

func (c *Comment) getIDColumns() []string {
	return []string{"CommentID"}
}

func (c *Comment) getInsertColumns() []string {
	result := c.getUpdateColumns()
	result = append(result, "EventID")
	result = append(result, "UserID")
	return result
}

func (c *Comment) getInsertFields() []interface{} {
	var interfaceSlice []interface{}
	interfaceSlice = append(interfaceSlice, c.getUpdateFields()...)
	interfaceSlice = append(interfaceSlice, c.EventID)
	interfaceSlice = append(interfaceSlice, c.UserID)
	return interfaceSlice
}

func (c *Comment) getUpdateColumns() []string {
	return []string{"Creation", "Comment"}
}

func (c *Comment) getUpdateFields() []interface{} {
	var interfaceSlice []interface{}
	interfaceSlice = append(interfaceSlice, c.Creation)
	interfaceSlice = append(interfaceSlice, c.Comment)
	return interfaceSlice
}

//Insert inserts a new Comment into the database and adding the new CommentID into the struct
func (c *Comment) Insert() error {
	var err error
	c.CommentID, err = insertDBO(c)
	if err != nil {
		err := errors.New("Error inserting Comment:" + err.Error())
		return err
	}
	return nil
}

//GetComments gives back all Comments in the Database
func GetComments(eventid int) ([]Comment, error) {
	var c []Comment
	var co Comment
	query, in := buildSelectQuery(&co, "EventID", eventid)
	err := db.Select(&c, query, in...)
	if err != nil {
		return c, errors.New("Error getting Attendees:" + err.Error())
	}
	return c, nil
}

//GetDetails takes a Comment struct with only the CommentID and tries to fetch the remaining infos
func (c *Comment) GetDetails() error {
	err := getDetailsDBO(c)
	if err != nil {
		return errors.New("Error getting Comment details:" + err.Error())
	}
	return nil
}

//Patch patches a Comment with new Info from a second struct
func (c *Comment) Patch(cc Comment) error {
	c.Comment = helpers.CopyIfNotEmpty(c.Comment, cc.Comment)
	c.Creation = cc.Creation
	return nil
}

//Update updates all Comment Fields in the Database
func (c *Comment) Update() error {
	err := updateDBO(c)
	if err != nil {
		return errors.New("Error updating Comment:" + err.Error())
	}
	return nil
}

//DeleteComment deletes a Comment with the given CommentID
func DeleteComment(id int) error {
	query := db.Rebind(`DELETE FROM "Comments" WHERE "CommentID" = ?`)
	_, err := db.Exec(query, id)
	if err != nil {
		return errors.New("Error deleting Comment: " + err.Error())
	}
	return nil
}
