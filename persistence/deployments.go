package persistence

import (
	"fmt"
	"log"

	"github.com/vice-registry/vice-util/models"
	r "gopkg.in/gorethink/gorethink.v3"
)

var tableDeployments = "deployments"

// CreateDeployment creates the provided deployment
func CreateDeployment(item *models.Deployment) (*models.Deployment, error) {
	id := GenerateID(defaultIDLength)
	item.ID = id
	createItem(item, tableDeployments)
	return item, nil
}

// UpdateDeployment updates the provided deployment
func UpdateDeployment(item *models.Deployment) (*models.Deployment, error) {
	err := updateItem(item, item.ID, tableDeployments)
	if err != nil {
		return nil, err
	}
	return item, nil
}

// DeleteDeployment returns a single deployment by id
func DeleteDeployment(id string) error {
	err := deleteItem(id, tableDeployments)
	if err != nil {
		return err
	}
	return nil
}

// GetDeployment returns a single deployment by id
func GetDeployment(id string) (*models.Deployment, error) {
	cursor, err := r.DB(connectionProperties.Database).Table(tableDeployments).Get(id).Run(connectionProperties.Session)
	if err != nil {
		log.Printf("Error in persistence getItem: cannot get item %s from table %s: %s", id, tableDeployments, err)
		return nil, err

	}
	if cursor.IsNil() {
		// no results
		log.Printf("No result for getItem on table %s and id %s", tableDeployments, id)
		return nil, nil
	}

	var item models.Deployment
	err = cursor.One(&item)
	if err != nil {
		fmt.Printf("Error scanning database result: %s", err)
		return nil, err
	}

	return &item, nil
}

// GetDeployments returns an array of deployments of the authenticated user
func GetDeployments(user *models.User) ([]*models.Deployment, error) {
	table := tableDeployments
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

	var items []*models.Deployment
	err = resp.All(&items)
	if err != nil {
		log.Printf("Error in persistence getItemsByIndex: cannot read results for table %s on index %s with value %s: %s", table, index, search, err)
		return nil, err
	}
	return items, nil
}
