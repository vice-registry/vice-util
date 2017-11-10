package persistence

import (
	"fmt"
	"log"

	"github.com/vice-registry/vice-util/models"
	r "gopkg.in/gorethink/gorethink.v3"
)

var tableImages = "images"

// CreateImage creates the provided image
func CreateImage(item *models.Image) (*models.Image, error) {
	id := GenerateID(defaultIDLength)
	item.ID = id
	createItem(item, tableImages)
	return item, nil
}

// UpdateImage updates the provided image
func UpdateImage(item *models.Image) (*models.Image, error) {
	err := updateItem(item, item.ID, tableImages)
	if err != nil {
		return nil, err
	}
	return item, nil
}

// DeleteImage returns a single image by id
func DeleteImage(id string) error {
	err := deleteItem(id, tableImages)
	if err != nil {
		return err
	}
	return nil
}

// GetImage returns a single image by id
func GetImage(id string) (*models.Image, error) {
	cursor, err := r.DB(connectionProperties.Database).Table(tableImages).Get(id).Run(connectionProperties.Session)
	if err != nil {
		log.Printf("Error in persistence getItem: cannot get item %s from table %s: %s", id, tableImages, err)
		return nil, err

	}
	if cursor.IsNil() {
		// no results
		log.Printf("No result for getItem on table %s and id %s", tableImages, id)
		return nil, nil
	}

	var item models.Image
	err = cursor.One(&item)
	if err != nil {
		fmt.Printf("Error scanning database result: %s", err)
		return nil, err
	}

	return &item, nil
}

// GetImages returns an array of images of the authenticated user
func GetImages(user *models.User) ([]*models.Image, error) {
	table := tableImages
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

	var items []*models.Image
	err = resp.All(&items)
	if err != nil {
		log.Printf("Error in persistence getItemsByIndex: cannot read results for table %s on index %s with value %s: %s", table, index, search, err)
		return nil, err
	}
	return items, nil
}
