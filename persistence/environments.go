package persistence

import (
	"log"

	"github.com/vice-registry/vice-util/models"
	gocb "gopkg.in/couchbase/gocb.v1"
)

// CreateEnvironment creates the provided environment
func CreateEnvironment(item *models.Environment) (*models.Environment, error) {
	id := GenerateID(defaultIDLength)
	item.ID = id
	bucket, err := couchbaseCredentials.Cluster.OpenBucket("vice-environments", couchbaseCredentials.Password)
	if err != nil {
		log.Printf("Error in persistence CreateEnvironment: cannot open bucket %s: %s", "vice-environments", err)
		return nil, err
	}
	defer bucket.Close()
	_, err = bucket.Insert(id, item, 0)
	if err != nil {
		log.Printf("Error in persistence CreateEnvironment: cannot create item %+v in bucket %s: %s", item, "vice-environments", err)
		return nil, err
	}
	return item, nil
}

// UpdateEnvironment updates the provided environment
func UpdateEnvironment(item *models.Environment) (*models.Environment, error) {
	err := updateItem(item, item.ID, "vice-environments")
	if err != nil {
		return nil, err
	}
	return item, nil
}

// DeleteEnvironment returns a single environment by id
func DeleteEnvironment(id string) error {
	err := deleteItem(id, "vice-environments")
	if err != nil {
		return err
	}
	return nil
}

// GetEnvironment returns a single environment by id
func GetEnvironment(id string) (*models.Environment, error) {
	var item models.Environment
	err := getItem(&item, id, "vice-environments")
	if err != nil {
		return nil, err
	}
	return &item, nil
}

// GetEnvironments returns an array of environments of the authenticated user
func GetEnvironments(user *models.User) ([]*models.Environment, error) {
	query := gocb.NewN1qlQuery("SELECT environments.* FROM `vice-environments` AS environments WHERE `userid` LIKE  $1;")
	params := []interface{}{"%"}
	if user != nil {
		params = []interface{}{user.ID}
	}
	bucket, err := couchbaseCredentials.Cluster.OpenBucket("vice-environments", couchbaseCredentials.Password)
	if err != nil {
		log.Printf("Error in persistence GetEnvironments: cannot open bucket %s: %s", "vice-environments", err)
		return nil, err
	}
	rows, err := bucket.ExecuteN1qlQuery(query, params)
	if err != nil {
		log.Printf("Error in persistence GetEnvironments: cannot run query on bucket %s: %s", "vice-environments", err)
		return nil, err
	}
	var items []*models.Environment
	var item models.Environment
	for rows.Next(&item) {
		copy := new(models.Environment)
		*copy = item
		if item.ID != "" {
			items = append(items, copy)
		}
		item = models.Environment{}
	}
	return items, nil
}
