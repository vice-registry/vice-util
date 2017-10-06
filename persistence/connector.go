package persistence

import (
	"math/rand"

	"log"

	"gopkg.in/couchbase/gocb.v1"
)

const idLetterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
const defaultIDLength int = 6

var couchbaseCredentials = struct {
	Location string
	Username string
	Password string
	Cluster  *gocb.Cluster
}{}

// SetCouchbaseCredentials set the login credentials to Couchbase cluster
func SetCouchbaseCredentials(location string, username string, password string) {
	couchbaseCredentials.Location = location
	couchbaseCredentials.Username = username
	couchbaseCredentials.Password = password

	cluster, err := gocb.Connect("couchbase://" + couchbaseCredentials.Location)
	if err != nil {
		log.Fatalln("cannot connect to couchbase: ", err)
	}
	couchbaseCredentials.Cluster = cluster
}

// CloseConnection closes all open connections to the Couchbase cluster
func CloseConnection() {
	// TODO
}

// GenerateID provides a random string with n characters
func GenerateID(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = idLetterBytes[rand.Intn(len(idLetterBytes))]
	}
	return string(b)
}

func getItem(item interface{}, id string, bucketName string) error {
	bucket, err := couchbaseCredentials.Cluster.OpenBucket(bucketName, couchbaseCredentials.Password)
	if err != nil {
		log.Printf("Error in persistence getItem: cannot open bucket %s: %s", bucketName, err)
		return err
	}
	defer bucket.Close()
	_, err = bucket.Get(id, item)
	if err != nil {
		log.Printf("Error in persistence getItem: cannot get itemId %s from bucket %s: %s", id, bucketName, err)
		return err
	}
	return nil
}

func updateItem(item interface{}, id string, bucketName string) error {
	bucket, err := couchbaseCredentials.Cluster.OpenBucket(bucketName, couchbaseCredentials.Password)
	if err != nil {
		log.Printf("Error in persistence updateItem: cannot open bucket %s: %s", bucketName, err)
		return err
	}
	defer bucket.Close()
	_, err = bucket.Replace(id, item, 0, 0)
	if err != nil {
		log.Printf("Error in persistence updateItem: cannot update itemId %s from bucket %s: %s", id, bucketName, err)
		return err
	}
	return nil
}

func deleteItem(id string, bucketName string) error {
	bucket, err := couchbaseCredentials.Cluster.OpenBucket(bucketName, couchbaseCredentials.Password)
	if err != nil {
		log.Printf("Error in persistence deleteItem: cannot open bucket %s: %s", bucketName, err)
		return err
	}
	defer bucket.Close()
	_, err = bucket.Remove(id, 0)
	if err != nil {
		log.Printf("Error in persistence deleteItem: cannot delete itemId %s from bucket %s: %s", id, bucketName, err)
		return err
	}
	return nil
}
