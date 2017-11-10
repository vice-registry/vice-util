package persistence

import (
	"fmt"
	"log"

	"github.com/vice-registry/vice-util/models"
	r "gopkg.in/gorethink/gorethink.v3"
)

var tableUsers = "users"

// CreateUser creates the provided user
func CreateUser(item *models.User) (*models.User, error) {
	id := GenerateID(defaultIDLength)
	item.ID = id
	createItem(item, tableUsers)
	return item, nil
}

// UpdateUser updates the provided user
func UpdateUser(item *models.User) (*models.User, error) {
	err := updateItem(item, item.ID, tableUsers)
	if err != nil {
		return nil, err
	}
	return item, nil
}

// DeleteUser returns a single user by id
func DeleteUser(id string) error {
	err := deleteItem(id, tableUsers)
	if err != nil {
		return err
	}
	return nil
}

// GetUser returns a single user by id
func GetUser(id string) (*models.User, error) {
	cursor, err := r.DB(connectionProperties.Database).Table(tableUsers).Get(id).Run(connectionProperties.Session)
	if err != nil {
		log.Printf("Error in persistence getItem: cannot get item %s from table %s: %s", id, tableUsers, err)
		return nil, err

	}
	if cursor.IsNil() {
		// no results
		log.Printf("No result for getItem on table %s and id %s", tableUsers, id)
		return nil, nil
	}

	var item models.User
	err = cursor.One(&item)
	if err != nil {
		fmt.Printf("Error scanning database result: %s", err)
		return nil, err
	}

	return &item, nil
}

// GetUserByName returns a single user by username
func GetUserByName(name string) (*models.User, error) {
	table := tableUsers
	index := "Username"
	search := name
	resp, err := r.DB(connectionProperties.Database).Table(table).GetAllByIndex(index, search).Run(connectionProperties.Session)
	if err != nil {
		log.Printf("Error in persistence getItemsByIndex: cannot query %s on index %s with value %s: %s", table, index, search, err)
		return nil, err
	}
	defer resp.Close()

	if resp.IsNil() {
		return nil, nil
	}

	var item models.User
	err = resp.One(&item)
	if err != nil {
		log.Printf("Error in persistence getItemsByIndex: cannot read results for table %s on index %s with value %s: %s", table, index, search, err)
		return nil, err
	}
	return &item, nil
}
