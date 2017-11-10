package persistence

import (
	"fmt"
	"log"

	"github.com/vice-registry/vice-util/models"
	r "gopkg.in/gorethink/gorethink.v3"
)

var tableEnvironments = "environments"

// CreateEnvironment creates the provided environment
func CreateEnvironment(item *models.Environment) (*models.Environment, error) {
	id := GenerateID(defaultIDLength)
	item.ID = id
	createItem(item, tableEnvironments)
	return item, nil
}

// UpdateEnvironment updates the provided environment
func UpdateEnvironment(item *models.Environment) (*models.Environment, error) {
	err := updateItem(item, item.ID, tableEnvironments)
	if err != nil {
		return nil, err
	}
	return item, nil
}

// DeleteEnvironment returns a single environment by id
func DeleteEnvironment(id string) error {
	err := deleteItem(id, tableEnvironments)
	if err != nil {
		return err
	}
	return nil
}

// GetEnvironment returns a single environment by id
func GetEnvironment(id string) (*models.Environment, error) {
	cursor, err := r.DB(connectionProperties.Database).Table(tableEnvironments).Get(id).Run(connectionProperties.Session)
	if err != nil {
		log.Printf("Error in persistence getItem: cannot get item %s from table %s: %s", id, tableEnvironments, err)
		return nil, err

	}
	if cursor.IsNil() {
		// no results
		log.Printf("No result for getItem on table %s and id %s", tableEnvironments, id)
		return nil, nil
	}

	var item models.Environment
	err = cursor.One(&item)
	if err != nil {
		fmt.Printf("Error scanning database result: %s", err)
		return nil, err
	}

	return &item, nil
}

// GetEnvironments returns an array of environments of the authenticated user
func GetEnvironments(user *models.User) ([]*models.Environment, error) {
	table := tableEnvironments
	index := "Userid"
	search := user.ID
	resp, err := r.DB(connectionProperties.Database).Table(table).GetAllByIndex(index, search).Run(connectionProperties.Session)
	if err != nil {
		log.Printf("Error in persistence getItemsByIndex: cannot query %s on index %s with value %s: %s", table, index, search, err)
		return nil, err
	}
	defer resp.Close()

	if resp.IsNil() {
		return nil, nil
	}

	var items []*models.Environment
	err = resp.All(&items)
	if err != nil {
		log.Printf("Error in persistence getItemsByIndex: cannot read results for table %s on index %s with value %s: %s", table, index, search, err)
		return nil, err
	}
	return items, nil
}
