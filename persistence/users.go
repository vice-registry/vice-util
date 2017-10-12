package persistence

import (
	"log"

	"github.com/vice-registry/vice-util/models"
	gocb "gopkg.in/couchbase/gocb.v1"
)

// CreateUser creates the provided user
func CreateUser(item *models.User) (*models.User, error) {
	id := GenerateID(defaultIDLength)
	item.ID = id
	bucket, err := couchbaseCredentials.Cluster.OpenBucket("vice-users", couchbaseCredentials.Password)
	if err != nil {
		log.Printf("Error in persistence CreateUser: cannot open bucket %s: %s", "vice-users", err)
		return nil, err
	}
	defer bucket.Close()
	_, err = bucket.Insert(id, item, 0)
	if err != nil {
		log.Printf("Error in persistence CreateUser: cannot create item %+v in bucket %s: %s", item, "vice-users", err)
		return nil, err
	}
	return item, nil
}

// UpdateUser updates the provided user
func UpdateUser(item *models.User) (*models.User, error) {
	err := updateItem(item, item.ID, "vice-users")
	if err != nil {
		return nil, err
	}
	return item, nil
}

// DeleteUser returns a single user by id
func DeleteUser(id string) error {
	err := deleteItem(id, "vice-users")
	if err != nil {
		return err
	}
	return nil
}

// GetUser returns a single user by id
func GetUser(id string) (*models.User, error) {
	var item models.User
	err := getItem(&item, id, "vice-users")
	if err != nil {
		return nil, err
	}
	return &item, nil
}

// GetUserByName returns a single user by username
func GetUserByName(name string) (*models.User, error) {
	bucket, err := couchbaseCredentials.Cluster.OpenBucket("vice-users", couchbaseCredentials.Password)
	if err != nil {
		log.Printf("Error in persistence GetUserByName: cannot open bucket %s: %s", "vice-users", err)
		return nil, err
	}
	query := gocb.NewN1qlQuery("SELECT `id`, `username`, `password`, `fullname`, `email` FROM `vice-users` AS users WHERE `username`=$1;")
	var params []interface{}
	params = append(params, name)

	rows, err := bucket.ExecuteN1qlQuery(query, params)
	if err != nil {
		log.Printf("Error in persistence GetUserByName: cannot run query on bucket %s: %s", "vice-users", err)
		return nil, err
	}
	var item models.User
	rows.One(&item)
	return &item, nil
}
