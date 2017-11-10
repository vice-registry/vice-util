package persistence

import (
	"log"
	"math/rand"

	r "gopkg.in/gorethink/gorethink.v3"
)

const idLetterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
const defaultIDLength int = 6

var connectionProperties = struct {
	Location string
	Database string
	//Username string
	//Password string
	Session *r.Session
}{}

// SetConnectionProperties set the login credentials to Rethinkdb cluster
func SetConnectionProperties(location string, database string) {
	connectionProperties.Location = location
	connectionProperties.Database = database

	session, err := r.Connect(r.ConnectOpts{
		Address: location,
	})
	if err != nil {
		log.Fatalln("cannot connect to vice-db: ", err)
	}
	connectionProperties.Session = session

}

// CloseConnection closes all open connections to the RethinkDB cluster
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

func updateItem(item interface{}, id string, table string) error {
	_, err := r.DB(connectionProperties.Database).Table(table).Get(id).Update(item).RunWrite(connectionProperties.Session)
	if err != nil {
		log.Printf("Error in persistence updateItem: cannot update itemId %s from table %s: %s", id, table, err)
		return err
	}
	return nil
}

func createItem(item interface{}, table string) error {
	_, err := r.DB(connectionProperties.Database).Table(table).Insert(item).RunWrite(connectionProperties.Session)
	if err != nil {
		log.Printf("Error in persistence createItem: cannot insert item %v from table %s: %s", item, table, err)
		return err
	}
	return nil
}

func deleteItem(id string, table string) error {
	_, err := r.DB(connectionProperties.Database).Table(table).Get(id).Delete().RunWrite(connectionProperties.Session)
	if err != nil {
		log.Printf("Error in persistence deleteItem: cannot delete itemId %s from table %s: %s", id, table, err)
		return err
	}
	return nil
}
